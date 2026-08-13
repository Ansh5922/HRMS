package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/repository"
	"github.com/your-org/hrms-backend/internal/service"
	"github.com/your-org/hrms-backend/pkg/response"
)

type PayrollHandler struct {
	service service.PayrollService
	empRepo repository.EmployeeRepository
}

func NewPayrollHandler(svc service.PayrollService, empRepo repository.EmployeeRepository) *PayrollHandler {
	return &PayrollHandler{service: svc, empRepo: empRepo}
}

// ── Request DTOs ─────────────────────────────────────────────────────────────

type AssignSalaryReq struct {
	EmpID          string  `json:"emp_id" binding:"required"`
	StructureID    string  `json:"structure_id" binding:"required"`
	CTC            float64 `json:"ctc" binding:"required"`
	Basic          float64 `json:"basic"`
	EffectiveFrom  string  `json:"effective_from" binding:"required"` // YYYY-MM-DD
}

type CreatePayrollRunReq struct {
	Month int `json:"month" binding:"required"`
	Year  int `json:"year" binding:"required"`
}

type SubmitReimbursementReq struct {
	Category    string  `json:"category" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Description *string `json:"description"`
	ReceiptURL  *string `json:"receipt_url"`
}

type SaveTaxDeclarationReq struct {
	Year     int     `json:"year" binding:"required"`
	Regime   string  `json:"regime" binding:"required"` // old, new
	Sections *string `json:"sections"`
}

// ── Salary Structure Endpoints ───────────────────────────────────────────────

func (h *PayrollHandler) CreateSalaryStructure(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	var st models.SalaryStructure
	if err := c.ShouldBindJSON(&st); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	st.OrgID = orgID

	res, err := h.service.CreateSalaryStructure(c.Request.Context(), &st)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *PayrollHandler) GetSalaryStructure(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid structure ID")
		return
	}

	st, err := h.service.GetSalaryStructure(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "salary structure")
		return
	}

	response.OK(c, st)
}

func (h *PayrollHandler) ListSalaryStructures(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	list, err := h.service.ListSalaryStructures(c.Request.Context(), orgID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, list)
}

func (h *PayrollHandler) DeleteSalaryStructure(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid structure ID")
		return
	}

	if err := h.service.DeleteSalaryStructure(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "salary structure deleted")
}

// ── Employee Salary Assignment Endpoints ─────────────────────────────────────

func (h *PayrollHandler) AssignEmployeeSalary(c *gin.Context) {
	var req AssignSalaryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	empID, err := uuid.Parse(req.EmpID)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}
	stID, err := uuid.Parse(req.StructureID)
	if err != nil {
		response.BadRequest(c, "invalid structure_id")
		return
	}

	sal := &models.EmployeeSalary{
		EmpID:       empID,
		StructureID: stID,
		CTC:         req.CTC,
		Basic:       req.Basic,
	}

	res, err := h.service.AssignEmployeeSalary(c.Request.Context(), sal)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, res, "salary structure assigned to employee")
}

func (h *PayrollHandler) GetEmployeeSalary(c *gin.Context) {
	empIDStr := c.Param("emp_id")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}

	sal, err := h.service.GetEmployeeSalary(c.Request.Context(), empID)
	if err != nil {
		response.NotFound(c, "employee salary record")
		return
	}

	response.OK(c, sal)
}

// ── Payroll Execution & Processing Endpoints ─────────────────────────────────

func (h *PayrollHandler) CreatePayrollRun(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	var req CreatePayrollRunReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	run := &models.PayrollRun{
		OrgID: orgID,
		Month: req.Month,
		Year:  req.Year,
		RunBy: &userID,
	}

	res, err := h.service.CreatePayrollRun(c.Request.Context(), run)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *PayrollHandler) ProcessPayrollRun(c *gin.Context) {
	idStr := c.Param("id")
	runID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid run ID")
		return
	}

	res, err := h.service.ProcessPayrollRun(c.Request.Context(), runID, h.empRepo)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, res, "payroll run processed successfully")
}

func (h *PayrollHandler) GetPayrollRun(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid run ID")
		return
	}

	run, err := h.service.GetPayrollRun(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "payroll run")
		return
	}

	response.OK(c, run)
}

func (h *PayrollHandler) ListPayrollRuns(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	runs, err := h.service.ListPayrollRuns(c.Request.Context(), orgID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, runs)
}

// ── Payslips Endpoints ───────────────────────────────────────────────────────

func (h *PayrollHandler) GetPayslip(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid payslip ID")
		return
	}

	p, err := h.service.GetPayslipByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "payslip")
		return
	}

	response.OK(c, p)
}

func (h *PayrollHandler) ListPayslipsByRun(c *gin.Context) {
	runIDStr := c.Param("run_id")
	runID, err := uuid.Parse(runIDStr)
	if err != nil {
		response.BadRequest(c, "invalid run_id")
		return
	}

	payslips, err := h.service.ListPayslipsByRun(c.Request.Context(), runID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, payslips)
}

func (h *PayrollHandler) ListPayslipsByEmp(c *gin.Context) {
	empIDStr := c.Param("emp_id")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}

	payslips, err := h.service.ListPayslipsByEmp(c.Request.Context(), empID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, payslips)
}

// ── Reimbursements Endpoints ─────────────────────────────────────────────────

func (h *PayrollHandler) SubmitReimbursement(c *gin.Context) {
	empIDStr := c.Param("emp_id")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}

	var req SubmitReimbursementReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rm := &models.Reimbursement{
		EmpID:       empID,
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
		ReceiptURL:  req.ReceiptURL,
	}

	res, err := h.service.SubmitReimbursement(c.Request.Context(), rm)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *PayrollHandler) ListMyReimbursements(c *gin.Context) {
	empIDStr := c.Param("emp_id")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}

	list, err := h.service.ListReimbursementsByEmp(c.Request.Context(), empID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, list)
}

func (h *PayrollHandler) ListOrgReimbursements(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	status := c.Query("status")
	list, err := h.service.ListReimbursementsByOrg(c.Request.Context(), orgID, status)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, list)
}

func (h *PayrollHandler) ApproveReimbursement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid reimbursement ID")
		return
	}

	userIDStr := c.GetString("user_id")
	approverID, _ := uuid.Parse(userIDStr)

	if err := h.service.ApproveReimbursement(c.Request.Context(), id, approverID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, nil, "reimbursement approved")
}

func (h *PayrollHandler) RejectReimbursement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid reimbursement ID")
		return
	}

	userIDStr := c.GetString("user_id")
	approverID, _ := uuid.Parse(userIDStr)

	if err := h.service.RejectReimbursement(c.Request.Context(), id, approverID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, nil, "reimbursement rejected")
}

// ── Tax Declaration Endpoints ────────────────────────────────────────────────

func (h *PayrollHandler) SaveTaxDeclaration(c *gin.Context) {
	empIDStr := c.Param("emp_id")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}

	var req SaveTaxDeclarationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	t := &models.TaxDeclaration{
		EmpID:    empID,
		Year:     req.Year,
		Regime:   req.Regime,
		Sections: req.Sections,
	}

	res, err := h.service.SaveTaxDeclaration(c.Request.Context(), t)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, res, "tax declaration saved")
}

func (h *PayrollHandler) GetTaxDeclaration(c *gin.Context) {
	empIDStr := c.Param("emp_id")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}

	yearStr := c.Query("year")
	year, _ := strconv.Atoi(yearStr)
	if year == 0 {
		response.BadRequest(c, "year query param is required")
		return
	}

	t, err := h.service.GetTaxDeclaration(c.Request.Context(), empID, year)
	if err != nil {
		response.NotFound(c, "tax declaration")
		return
	}

	response.OK(c, t)
}