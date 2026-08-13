package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/your-org/hrms-backend/docs"
	"github.com/your-org/hrms-backend/internal/config"
	"github.com/your-org/hrms-backend/internal/handler"
	"github.com/your-org/hrms-backend/internal/middleware"
	"github.com/your-org/hrms-backend/internal/repository"
	"github.com/your-org/hrms-backend/internal/service"
	"github.com/your-org/hrms-backend/pkg/aiservice"
)

func Setup(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// ── Interactive Swagger UI Route ──────────────────────────────────────
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// AI Client
	aiClient := aiservice.NewClient(cfg.AIServiceURL, cfg.AIServiceSecret)

	// ── Initialize Repositories ──────────────────────────────────────────
	authRepo := repository.NewAuthRepository(db)
	empRepo := repository.NewEmployeeRepository(db)
	attRepo := repository.NewAttendanceRepository(db)
	leaveRepo := repository.NewLeaveRepository(db)
	payrollRepo := repository.NewPayrollRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	// ── Initialize Services ──────────────────────────────────────────────
	authSvc := service.NewAuthService(authRepo, cfg)
	empSvc := service.NewEmployeeService(empRepo)
	attSvc := service.NewAttendanceService(attRepo)
	leaveSvc := service.NewLeaveService(leaveRepo)
	payrollSvc := service.NewPayrollService(payrollRepo)
	workflowSvc := service.NewWorkflowService(workflowRepo)
	notifSvc := service.NewNotificationService(notifRepo)

	// ── Initialize Handlers ──────────────────────────────────────────────
	authH := handler.NewAuthHandler(authSvc)
	empH := handler.NewEmployeeHandler(empSvc)
	attH := handler.NewAttendanceHandler(attSvc, aiClient)
	leaveH := handler.NewLeaveHandler(leaveSvc)
	payrollH := handler.NewPayrollHandler(payrollSvc, empRepo)
	workflowH := handler.NewWorkflowHandler(workflowSvc)
	notifH := handler.NewNotificationHandler(notifSvc)

	// ═════════════════════════════════════════════════════════════════════
	// ROUTES
	// ═════════════════════════════════════════════════════════════════════

	api := r.Group("/api/v1")
	{
		// ── Public Auth ──────────────────────────────────────────────
		auth := api.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
			auth.POST("/refresh", authH.RefreshToken)
			auth.POST("/forgot-password", authH.ForgotPassword)
			auth.POST("/reset-password", authH.ResetPassword)
		}

		// ── Protected Routes (JWT required) ──────────────────────────
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth(cfg))
		{
			// ── Profile ──────────────────────────────────────────────
			protected.GET("/auth/me", authH.GetMe)
			protected.PUT("/auth/change-password", authH.ChangePassword)
			protected.POST("/auth/logout", authH.Logout)
			protected.POST("/auth/logout-all", authH.LogoutAll)
			protected.GET("/employees/me", empH.GetMyProfile)

			// ── User Management ──────────────────────────────────────
			users := protected.Group("/users")
			{
				users.GET("", middleware.RequirePermission("users:read"), authH.ListUsers)
				users.POST("", middleware.RequirePermission("users:write"), authH.CreateUser)
				users.PUT("/:id", middleware.RequirePermission("users:write"), authH.UpdateUser)
				users.DELETE("/:id", middleware.RequirePermission("users:delete"), authH.DeactivateUser)
			}

			// ── Role Management ──────────────────────────────────────
			roles := protected.Group("/roles")
			{
				roles.GET("", middleware.RequirePermission("settings:read"), authH.ListRoles)
				roles.POST("", middleware.RequirePermission("settings:write"), authH.CreateRole)
				roles.POST("/assign", middleware.RequirePermission("users:write"), authH.AssignRole)
				roles.GET("/:id", middleware.RequirePermission("settings:read"), authH.GetRole)
				roles.PUT("/:id", middleware.RequirePermission("settings:write"), authH.UpdateRole)
				roles.DELETE("/:id", middleware.RequirePermission("settings:write"), authH.DeleteRole)
				roles.GET("/:id/permissions", middleware.RequirePermission("settings:read"), authH.GetRolePermissions)
				roles.PUT("/:id/permissions", middleware.RequirePermission("settings:write"), authH.SetRolePermissions)
			}

			protected.GET("/permissions", middleware.RequirePermission("settings:read"), authH.ListPermissions)

			// ── Department Management ─────────────────────────────────
			departments := protected.Group("/departments")
			{
				departments.GET("", middleware.RequirePermission("employee:read"), empH.ListDepartments)
				departments.POST("", middleware.RequirePermission("employee:write"), empH.CreateDepartment)
				departments.GET("/:id", middleware.RequirePermission("employee:read"), empH.GetDepartment)
				departments.PUT("/:id", middleware.RequirePermission("employee:write"), empH.UpdateDepartment)
				departments.DELETE("/:id", middleware.RequirePermission("employee:write"), empH.DeleteDepartment)
			}

			// ── Designation Management ────────────────────────────────
			designations := protected.Group("/designations")
			{
				designations.GET("", middleware.RequirePermission("employee:read"), empH.ListDesignations)
				designations.POST("", middleware.RequirePermission("employee:write"), empH.CreateDesignation)
				designations.GET("/:id", middleware.RequirePermission("employee:read"), empH.GetDesignation)
				designations.PUT("/:id", middleware.RequirePermission("employee:write"), empH.UpdateDesignation)
				designations.DELETE("/:id", middleware.RequirePermission("employee:write"), empH.DeleteDesignation)
			}

			// ── Employee Management ───────────────────────────────────
			employees := protected.Group("/employees")
			{
				employees.GET("", middleware.RequirePermission("employee:read"), empH.List)
				employees.POST("", middleware.RequirePermission("employee:write"), empH.Create)
				employees.GET("/:id", middleware.RequirePermission("employee:read"), empH.GetByID)
				employees.PUT("/:id", middleware.RequirePermission("employee:write"), empH.Update)
				employees.DELETE("/:id", middleware.RequirePermission("employee:delete"), empH.Delete)

				// Bank Details
				employees.POST("/:id/bank-details", middleware.RequirePermission("employee:write"), empH.SaveBankDetails)
				employees.GET("/:id/bank-details", middleware.RequirePermission("employee:read"), empH.GetBankDetails)

				// Emergency Contacts
				employees.POST("/:id/emergency-contacts", middleware.RequirePermission("employee:write"), empH.AddEmergencyContact)
				employees.GET("/:id/emergency-contacts", middleware.RequirePermission("employee:read"), empH.ListEmergencyContacts)
				employees.DELETE("/emergency-contacts/:id", middleware.RequirePermission("employee:write"), empH.DeleteEmergencyContact)

				// Documents
				employees.POST("/:id/documents", middleware.RequirePermission("employee:write"), empH.AddDocument)
				employees.GET("/:id/documents", middleware.RequirePermission("employee:read"), empH.ListDocuments)
				employees.PATCH("/documents/:id/verify", middleware.RequirePermission("employee:write"), empH.VerifyDocument)
				employees.DELETE("/documents/:id", middleware.RequirePermission("employee:write"), empH.DeleteDocument)

				// Onboarding Tasks
				employees.POST("/:id/onboarding", middleware.RequirePermission("employee:write"), empH.CreateOnboardingTask)
				employees.GET("/:id/onboarding", middleware.RequirePermission("employee:read"), empH.ListOnboardingTasks)
				employees.PATCH("/onboarding/:id", middleware.RequirePermission("employee:write"), empH.UpdateOnboardingTaskStatus)
			}

			// ── Leave Management Module ───────────────────────────────
			leaveTypes := protected.Group("/leave-types")
			{
				leaveTypes.GET("", middleware.RequirePermission("leave:read"), leaveH.ListLeaveTypes)
				leaveTypes.POST("", middleware.RequirePermission("settings:write"), leaveH.CreateLeaveType)
				leaveTypes.PUT("/:id", middleware.RequirePermission("settings:write"), leaveH.UpdateLeaveType)
				leaveTypes.DELETE("/:id", middleware.RequirePermission("settings:write"), leaveH.DeleteLeaveType)
			}

			leavePolicies := protected.Group("/leave-policies")
			{
				leavePolicies.GET("", middleware.RequirePermission("leave:read"), leaveH.ListLeavePolicies)
				leavePolicies.POST("", middleware.RequirePermission("settings:write"), leaveH.CreateLeavePolicy)
				leavePolicies.DELETE("/:id", middleware.RequirePermission("settings:write"), leaveH.DeleteLeavePolicy)
			}

			leaveBalances := protected.Group("/leave-balances")
			{
				leaveBalances.POST("", middleware.RequirePermission("leave:write"), leaveH.SetLeaveBalance)
				leaveBalances.GET("/employee/:emp_id", middleware.RequirePermission("leave:read"), leaveH.GetLeaveBalances)
			}

			leaves := protected.Group("/leaves")
			{
				leaves.GET("/org", middleware.RequirePermission("leave:read"), leaveH.ListByOrg)
				leaves.POST("/apply/:emp_id", middleware.RequirePermission("leave:write"), leaveH.Apply)
				leaves.GET("/employee/:emp_id", middleware.RequirePermission("leave:read"), leaveH.ListByEmp)
				leaves.GET("/:id", middleware.RequirePermission("leave:read"), leaveH.GetLeaveByID)
				leaves.PATCH("/:id/approve", middleware.RequirePermission("leave:approve"), leaveH.Approve)
				leaves.PATCH("/:id/reject", middleware.RequirePermission("leave:approve"), leaveH.Reject)
				leaves.POST("/:id/cancel", middleware.RequirePermission("leave:write"), leaveH.Cancel)
			}

			// ── Payroll & Compensation Module ─────────────────────────
			salaryStructs := protected.Group("/salary-structures")
			{
				salaryStructs.GET("", middleware.RequirePermission("payroll:read"), payrollH.ListSalaryStructures)
				salaryStructs.POST("", middleware.RequirePermission("payroll:write"), payrollH.CreateSalaryStructure)
				salaryStructs.GET("/:id", middleware.RequirePermission("payroll:read"), payrollH.GetSalaryStructure)
				salaryStructs.DELETE("/:id", middleware.RequirePermission("payroll:write"), payrollH.DeleteSalaryStructure)
			}

			employeeSalaries := protected.Group("/employee-salaries")
			{
				employeeSalaries.POST("/assign", middleware.RequirePermission("payroll:write"), payrollH.AssignEmployeeSalary)
				employeeSalaries.GET("/employee/:emp_id", middleware.RequirePermission("payroll:read"), payrollH.GetEmployeeSalary)
			}

			payrollRuns := protected.Group("/payroll-runs")
			{
				payrollRuns.GET("", middleware.RequirePermission("payroll:read"), payrollH.ListPayrollRuns)
				payrollRuns.POST("", middleware.RequirePermission("payroll:write"), payrollH.CreatePayrollRun)
				payrollRuns.GET("/:id", middleware.RequirePermission("payroll:read"), payrollH.GetPayrollRun)
				payrollRuns.POST("/:id/process", middleware.RequirePermission("payroll:write"), payrollH.ProcessPayrollRun)
			}

			payslips := protected.Group("/payslips")
			{
				payslips.GET("/run/:run_id", middleware.RequirePermission("payroll:read"), payrollH.ListPayslipsByRun)
				payslips.GET("/employee/:emp_id", middleware.RequirePermission("payroll:read"), payrollH.ListPayslipsByEmp)
				payslips.GET("/:id", middleware.RequirePermission("payroll:read"), payrollH.GetPayslip)
			}

			reimbursements := protected.Group("/reimbursements")
			{
				reimbursements.GET("/org", middleware.RequirePermission("payroll:read"), payrollH.ListOrgReimbursements)
				reimbursements.POST("/employee/:emp_id", middleware.RequirePermission("payroll:write"), payrollH.SubmitReimbursement)
				reimbursements.GET("/employee/:emp_id", middleware.RequirePermission("payroll:read"), payrollH.ListMyReimbursements)
				reimbursements.PATCH("/:id/approve", middleware.RequirePermission("payroll:approve"), payrollH.ApproveReimbursement)
				reimbursements.PATCH("/:id/reject", middleware.RequirePermission("payroll:approve"), payrollH.RejectReimbursement)
			}

			taxDeclarations := protected.Group("/tax-declarations")
			{
				taxDeclarations.POST("/employee/:emp_id", middleware.RequirePermission("payroll:write"), payrollH.SaveTaxDeclaration)
				taxDeclarations.GET("/employee/:emp_id", middleware.RequirePermission("payroll:read"), payrollH.GetTaxDeclaration)
			}

			// ── Workflow & Custom Approvals Engine ────────────────────
			workflows := protected.Group("/workflows")
			{
				workflows.GET("/templates", middleware.RequirePermission("settings:read"), workflowH.ListTemplates)
				workflows.POST("/templates", middleware.RequirePermission("settings:write"), workflowH.CreateTemplate)
				workflows.GET("/templates/:id", middleware.RequirePermission("settings:read"), workflowH.GetTemplate)
				workflows.PUT("/templates/:id", middleware.RequirePermission("settings:write"), workflowH.UpdateTemplate)
				workflows.DELETE("/templates/:id", middleware.RequirePermission("settings:write"), workflowH.DeleteTemplate)
			}

			approvals := protected.Group("/approvals")
			{
				approvals.POST("/submit", middleware.RequirePermission("leave:write"), workflowH.SubmitRequest)
				approvals.GET("/org", middleware.RequirePermission("leave:read"), workflowH.ListOrgRequests)
				approvals.GET("/pending-me", middleware.RequirePermission("leave:read"), workflowH.ListMyPendingApprovals)
				approvals.GET("/my-requests", middleware.RequirePermission("leave:read"), workflowH.ListMySubmittedRequests)
				approvals.GET("/:id", middleware.RequirePermission("leave:read"), workflowH.GetRequest)
				approvals.PATCH("/:id/approve", middleware.RequirePermission("leave:approve"), workflowH.ApproveStep)
				approvals.PATCH("/:id/reject", middleware.RequirePermission("leave:approve"), workflowH.RejectStep)
			}

			// ── Notifications & Announcements Module ─────────────────
			notifications := protected.Group("/notifications")
			{
				notifications.GET("/me", notifH.ListMyNotifications)
				notifications.GET("/unread-count", notifH.GetUnreadCount)
				notifications.PATCH("/read-all", notifH.MarkAllAsRead)
				notifications.PATCH("/:id/read", notifH.MarkAsRead)
				notifications.DELETE("/:id", notifH.DeleteNotification)
			}

			announcements := protected.Group("/announcements")
			{
				announcements.GET("", notifH.ListAnnouncements)
				announcements.POST("", middleware.RequirePermission("settings:write"), notifH.CreateAnnouncement)
				announcements.GET("/:id", notifH.GetAnnouncement)
				announcements.DELETE("/:id", middleware.RequirePermission("settings:write"), notifH.DeleteAnnouncement)
			}

			notifPrefs := protected.Group("/notification-preferences")
			{
				notifPrefs.GET("", notifH.GetPreferences)
				notifPrefs.POST("", notifH.SavePreference)
			}

			pushTokens := protected.Group("/push-tokens")
			{
				pushTokens.POST("", notifH.RegisterPushToken)
				pushTokens.DELETE("", notifH.RemovePushToken)
			}

			// ── Attendance (requires attendance:read/write) ──────
			attendance := protected.Group("/attendance")
			{
				attendance.GET("/today", middleware.RequirePermission("attendance:read"), attH.GetToday)
				attendance.POST("/check-in", middleware.RequirePermission("attendance:write"), attH.CheckIn)
				attendance.POST("/face-checkin", middleware.RequirePermission("attendance:write"), attH.FaceCheckIn)
				attendance.POST("/check-out/:emp_id", middleware.RequirePermission("attendance:write"), attH.CheckOut)
			}
		}
	}

	return r
}
