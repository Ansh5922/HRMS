package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/service"
	"github.com/your-org/hrms-backend/pkg/response"
)

type EmployeeHandler struct {
	service service.EmployeeService
}

func NewEmployeeHandler(svc service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: svc}
}

func getIDParam(c *gin.Context) string {
	id := c.Param("id")
	if id == "" {
		id = c.Param("emp_id")
	}
	return id
}

// ── Request DTOs ─────────────────────────────────────────────────────────────

type CreateDeptReq struct {
	Name         string  `json:"name" binding:"required"`
	Code         *string `json:"code"`
	ParentDeptID *string `json:"parent_dept_id"`
	HeadEmpID    *string `json:"head_emp_id"`
	Description  *string `json:"description"`
}

type UpdateDeptReq struct {
	Name         *string `json:"name"`
	Code         *string `json:"code"`
	ParentDeptID *string `json:"parent_dept_id"`
	HeadEmpID    *string `json:"head_emp_id"`
	Description  *string `json:"description"`
	IsActive     *bool   `json:"is_active"`
}

type CreateDesigReq struct {
	Title        string  `json:"title" binding:"required"`
	DepartmentID *string `json:"department_id"`
	Grade        *string `json:"grade"`
	Level        *int    `json:"level"`
}

type UpdateDesigReq struct {
	Title        *string `json:"title"`
	DepartmentID *string `json:"department_id"`
	Grade        *string `json:"grade"`
	Level        *int    `json:"level"`
}

type SaveBankDetailsReq struct {
	BankName      string  `json:"bank_name" binding:"required"`
	AccountNumber string  `json:"account_number" binding:"required"`
	IFSCCode      string  `json:"ifsc_code" binding:"required"`
	AccountType   *string `json:"account_type"`
}

type AddEmergencyContactReq struct {
	Name     string  `json:"name" binding:"required"`
	Relation *string `json:"relation"`
	Phone    string  `json:"phone" binding:"required"`
	Email    *string `json:"email"`
	Address  *string `json:"address"`
}

type AddDocumentReq struct {
	DocType  string  `json:"doc_type" binding:"required"`
	FileURL  string  `json:"file_url" binding:"required"`
	FileName *string `json:"file_name"`
}

type CreateOnboardingTaskReq struct {
	TaskName    string  `json:"task_name" binding:"required"`
	Description *string `json:"description"`
	AssignedTo  *string `json:"assigned_to"`
	DueDate     *string `json:"due_date"`
}

type UpdateTaskStatusReq struct {
	Status string `json:"status" binding:"required"`
}

// ── Employee Core Endpoints ──────────────────────────────────────────────────

func (h *EmployeeHandler) Create(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	var emp models.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	emp.OrgID = orgID

	res, err := h.service.CreateEmployee(c.Request.Context(), &emp)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *EmployeeHandler) GetByID(c *gin.Context) {
	idStr := getIDParam(c)
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid employee ID")
		return
	}

	emp, err := h.service.GetEmployeeByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "employee")
		return
	}

	response.OK(c, emp)
}

func (h *EmployeeHandler) GetMyProfile(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	emp, err := h.service.GetEmployeeByUserID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "employee profile")
		return
	}

	response.OK(c, emp)
}

func (h *EmployeeHandler) List(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	search := c.Query("search")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var deptID, desigID *uuid.UUID
	if d := c.Query("dept_id"); d != "" {
		if parsed, err := uuid.Parse(d); err == nil {
			deptID = &parsed
		}
	}
	if des := c.Query("desig_id"); des != "" {
		if parsed, err := uuid.Parse(des); err == nil {
			desigID = &parsed
		}
	}

	emps, total, err := h.service.ListEmployees(c.Request.Context(), orgID, search, deptID, desigID, status, page, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"employees": emps,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	idStr := getIDParam(c)
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid employee ID")
		return
	}

	var emp models.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	emp.ID = id
	emp.OrgID = orgID

	if err := h.service.UpdateEmployee(c.Request.Context(), &emp); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "employee updated successfully")
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	idStr := getIDParam(c)
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid employee ID")
		return
	}

	if err := h.service.DeleteEmployee(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "employee deleted successfully")
}

// ── Department Endpoints ─────────────────────────────────────────────────────

func (h *EmployeeHandler) CreateDepartment(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	var req CreateDeptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	dept := &models.Department{
		OrgID:       orgID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    true,
	}
	if req.ParentDeptID != nil {
		if p, err := uuid.Parse(*req.ParentDeptID); err == nil {
			dept.ParentDeptID = &p
		}
	}
	if req.HeadEmpID != nil {
		if h, err := uuid.Parse(*req.HeadEmpID); err == nil {
			dept.HeadEmpID = &h
		}
	}

	res, err := h.service.CreateDepartment(c.Request.Context(), dept)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *EmployeeHandler) GetDepartment(c *gin.Context) {
	idStr := getIDParam(c)
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid department ID")
		return
	}

	dept, err := h.service.GetDepartment(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "department")
		return
	}

	response.OK(c, dept)
}

func (h *EmployeeHandler) ListDepartments(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	depts, err := h.service.ListDepartments(c.Request.Context(), orgID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, depts)
}

func (h *EmployeeHandler) UpdateDepartment(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	idStr := getIDParam(c)
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid department ID")
		return
	}

	var req UpdateDeptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	dept := &models.Department{ID: id, OrgID: orgID}
	if req.Name != nil {
		dept.Name = *req.Name
	}
	if req.Code != nil {
		dept.Code = req.Code
	}
	if req.Description != nil {
		dept.Description = req.Description
	}
	if req.IsActive != nil {
		dept.IsActive = *req.IsActive
	}
	if req.ParentDeptID != nil {
		p, _ := uuid.Parse(*req.ParentDeptID)
		dept.ParentDeptID = &p
	}
	if req.HeadEmpID != nil {
		hEmp, _ := uuid.Parse(*req.HeadEmpID)
		dept.HeadEmpID = &hEmp
	}

	if err := h.service.UpdateDepartment(c.Request.Context(), dept); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "department updated")
}

func (h *EmployeeHandler) DeleteDepartment(c *gin.Context) {
	idStr := getIDParam(c)
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid department ID")
		return
	}

	if err := h.service.DeleteDepartment(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "department deleted")
}

// ── Designation Endpoints ────────────────────────────────────────────────────

func (h *EmployeeHandler) CreateDesignation(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	var req CreateDesigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	des := &models.Designation{
		OrgID: orgID,
		Title: req.Title,
		Grade: req.Grade,
		Level: req.Level,
	}
	if req.DepartmentID != nil {
		if d, err := uuid.Parse(*req.DepartmentID); err == nil {
			des.DepartmentID = &d
		}
	}

	res, err := h.service.CreateDesignation(c.Request.Context(), des)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *EmployeeHandler) GetDesignation(c *gin.Context) {
	idStr := getIDParam(c)
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid designation ID")
		return
	}

	des, err := h.service.GetDesignation(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "designation")
		return
	}

	response.OK(c, des)
}

func (h *EmployeeHandler) ListDesignations(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	desigs, err := h.service.ListDesignations(c.Request.Context(), orgID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, desigs)
}

func (h *EmployeeHandler) UpdateDesignation(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	idStr := getIDParam(c)
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid designation ID")
		return
	}

	var req UpdateDesigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	des := &models.Designation{ID: id, OrgID: orgID}
	if req.Title != nil {
		des.Title = *req.Title
	}
	if req.Grade != nil {
		des.Grade = req.Grade
	}
	if req.Level != nil {
		des.Level = req.Level
	}
	if req.DepartmentID != nil {
		d, _ := uuid.Parse(*req.DepartmentID)
		des.DepartmentID = &d
	}

	if err := h.service.UpdateDesignation(c.Request.Context(), des); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "designation updated")
}

func (h *EmployeeHandler) DeleteDesignation(c *gin.Context) {
	idStr := getIDParam(c)
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid designation ID")
		return
	}

	if err := h.service.DeleteDesignation(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "designation deleted")
}

// ── Bank Details Endpoints ───────────────────────────────────────────────────

func (h *EmployeeHandler) SaveBankDetails(c *gin.Context) {
	empIDStr := getIDParam(c)
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid employee ID")
		return
	}

	var req SaveBankDetailsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	bank := &models.EmployeeBankDetails{
		EmpID:        empID,
		BankName:     req.BankName,
		AccountNoEnc: req.AccountNumber,
		IFSCCode:     req.IFSCCode,
		AccountType:  req.AccountType,
		IsPrimary:    true,
	}

	res, err := h.service.SaveBankDetails(c.Request.Context(), bank)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, res, "bank details saved")
}

func (h *EmployeeHandler) GetBankDetails(c *gin.Context) {
	empIDStr := getIDParam(c)
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid employee ID")
		return
	}

	res, err := h.service.GetBankDetails(c.Request.Context(), empID)
	if err != nil {
		response.NotFound(c, "bank details")
		return
	}

	response.OK(c, res)
}

// ── Emergency Contacts Endpoints ─────────────────────────────────────────────

func (h *EmployeeHandler) AddEmergencyContact(c *gin.Context) {
	empIDStr := getIDParam(c)
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid employee ID")
		return
	}

	var req AddEmergencyContactReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	contact := &models.EmergencyContact{
		EmpID:    empID,
		Name:     req.Name,
		Relation: req.Relation,
		Phone:    req.Phone,
		Email:    req.Email,
		Address:  req.Address,
	}

	res, err := h.service.AddEmergencyContact(c.Request.Context(), contact)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *EmployeeHandler) ListEmergencyContacts(c *gin.Context) {
	empIDStr := getIDParam(c)
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid employee ID")
		return
	}

	contacts, err := h.service.ListEmergencyContacts(c.Request.Context(), empID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, contacts)
}

func (h *EmployeeHandler) DeleteEmergencyContact(c *gin.Context) {
	idStr := getIDParam(c)
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid contact ID")
		return
	}

	if err := h.service.DeleteEmergencyContact(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "emergency contact deleted")
}

// ── Employee Document Endpoints ──────────────────────────────────────────────

func (h *EmployeeHandler) AddDocument(c *gin.Context) {
	empIDStr := getIDParam(c)
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid employee ID")
		return
	}

	var req AddDocumentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	doc := &models.EmployeeDocument{
		EmpID:    empID,
		DocType:  req.DocType,
		FileURL:  req.FileURL,
		FileName: req.FileName,
	}

	res, err := h.service.AddDocument(c.Request.Context(), doc)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *EmployeeHandler) ListDocuments(c *gin.Context) {
	empIDStr := getIDParam(c)
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid employee ID")
		return
	}

	docs, err := h.service.ListDocuments(c.Request.Context(), empID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, docs)
}

func (h *EmployeeHandler) VerifyDocument(c *gin.Context) {
	docIDStr := getIDParam(c)
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		response.BadRequest(c, "invalid document ID")
		return
	}

	userIDStr := c.GetString("user_id")
	verifierID, _ := uuid.Parse(userIDStr)

	if err := h.service.VerifyDocument(c.Request.Context(), docID, verifierID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "document verified")
}

func (h *EmployeeHandler) DeleteDocument(c *gin.Context) {
	docIDStr := getIDParam(c)
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		response.BadRequest(c, "invalid document ID")
		return
	}

	if err := h.service.DeleteDocument(c.Request.Context(), docID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "document deleted")
}

// ── Onboarding Task Endpoints ────────────────────────────────────────────────

func (h *EmployeeHandler) CreateOnboardingTask(c *gin.Context) {
	empIDStr := getIDParam(c)
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid employee ID")
		return
	}

	var req CreateOnboardingTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	task := &models.OnboardingTask{
		EmpID:       empID,
		TaskName:    req.TaskName,
		Description: req.Description,
	}
	if req.AssignedTo != nil {
		if u, err := uuid.Parse(*req.AssignedTo); err == nil {
			task.AssignedTo = &u
		}
	}

	res, err := h.service.CreateOnboardingTask(c.Request.Context(), task)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *EmployeeHandler) ListOnboardingTasks(c *gin.Context) {
	empIDStr := getIDParam(c)
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid employee ID")
		return
	}

	tasks, err := h.service.ListOnboardingTasks(c.Request.Context(), empID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, tasks)
}

func (h *EmployeeHandler) UpdateOnboardingTaskStatus(c *gin.Context) {
	taskIDStr := getIDParam(c)
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		response.BadRequest(c, "invalid task ID")
		return
	}

	var req UpdateTaskStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateOnboardingTask(c.Request.Context(), taskID, req.Status); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "onboarding task updated")
}