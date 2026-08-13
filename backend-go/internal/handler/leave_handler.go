package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/service"
	"github.com/your-org/hrms-backend/pkg/response"
)

type LeaveHandler struct {
	service service.LeaveService
}

func NewLeaveHandler(svc service.LeaveService) *LeaveHandler {
	return &LeaveHandler{service: svc}
}

// ── Request DTOs ─────────────────────────────────────────────────────────────

type CreateLeaveTypeReq struct {
	Name             string `json:"name" binding:"required"`
	Code             string `json:"code" binding:"required"`
	MaxDaysPerYear   *int   `json:"max_days_per_year"`
	CarryForward     bool   `json:"carry_forward"`
	CarryForwardMax  *int   `json:"carry_forward_max"`
	IsPaid           bool   `json:"is_paid"`
	RequiresDoc      bool   `json:"requires_doc"`
	MinNoticeDays    int    `json:"min_notice_days"`
	GenderApplicable string `json:"gender_applicable"`
	RequiresApproval bool   `json:"requires_approval"`
}

type CreateLeavePolicyReq struct {
	LeaveTypeID   string  `json:"leave_type_id" binding:"required"`
	DesignationID *string `json:"designation_id"`
	DepartmentID  *string `json:"department_id"`
	AnnualQuota   float64 `json:"annual_quota" binding:"required"`
	AccrualType   string  `json:"accrual_type"`
}

type SetLeaveBalanceReq struct {
	EmpID       string  `json:"emp_id" binding:"required"`
	LeaveTypeID string  `json:"leave_type_id" binding:"required"`
	Year        int     `json:"year"`
	Total       float64 `json:"total" binding:"required"`
}

type ApplyLeaveReq struct {
	LeaveTypeID string  `json:"leave_type_id" binding:"required"`
	FromDate    string  `json:"from_date" binding:"required"` // YYYY-MM-DD
	ToDate      string  `json:"to_date" binding:"required"`   // YYYY-MM-DD
	Session     string  `json:"session"`                     // full, first_half, second_half
	Reason      *string `json:"reason"`
	DocURL      *string `json:"doc_url"`
}

type ActionLeaveReq struct {
	Comment *string `json:"comment"`
}

// ── Leave Type Endpoints ─────────────────────────────────────────────────────

func (h *LeaveHandler) CreateLeaveType(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	var req CreateLeaveTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	lt := &models.LeaveType{
		OrgID:            orgID,
		Name:             req.Name,
		Code:             req.Code,
		MaxDaysPerYear:   req.MaxDaysPerYear,
		CarryForward:     req.CarryForward,
		CarryForwardMax:  req.CarryForwardMax,
		IsPaid:           req.IsPaid,
		RequiresDoc:      req.RequiresDoc,
		MinNoticeDays:    req.MinNoticeDays,
		GenderApplicable: req.GenderApplicable,
		RequiresApproval: req.RequiresApproval,
	}

	res, err := h.service.CreateLeaveType(c.Request.Context(), lt)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *LeaveHandler) ListLeaveTypes(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	types, err := h.service.ListLeaveTypes(c.Request.Context(), orgID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, types)
}

func (h *LeaveHandler) UpdateLeaveType(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid leave type ID")
		return
	}

	var lt models.LeaveType
	if err := c.ShouldBindJSON(&lt); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	lt.ID = id
	lt.OrgID = orgID

	if err := h.service.UpdateLeaveType(c.Request.Context(), &lt); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "leave type updated")
}

func (h *LeaveHandler) DeleteLeaveType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid leave type ID")
		return
	}

	if err := h.service.DeleteLeaveType(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "leave type deleted")
}

// ── Leave Policy Endpoints ───────────────────────────────────────────────────

func (h *LeaveHandler) CreateLeavePolicy(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	var req CreateLeavePolicyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	typeID, err := uuid.Parse(req.LeaveTypeID)
	if err != nil {
		response.BadRequest(c, "invalid leave_type_id")
		return
	}

	lp := &models.LeavePolicy{
		OrgID:       orgID,
		LeaveTypeID: typeID,
		AnnualQuota: req.AnnualQuota,
		AccrualType: req.AccrualType,
	}
	if req.DesignationID != nil {
		if des, err := uuid.Parse(*req.DesignationID); err == nil {
			lp.DesignationID = &des
		}
	}
	if req.DepartmentID != nil {
		if d, err := uuid.Parse(*req.DepartmentID); err == nil {
			lp.DepartmentID = &d
		}
	}

	res, err := h.service.CreateLeavePolicy(c.Request.Context(), lp)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *LeaveHandler) ListLeavePolicies(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	policies, err := h.service.ListLeavePolicies(c.Request.Context(), orgID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, policies)
}

func (h *LeaveHandler) DeleteLeavePolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid policy ID")
		return
	}

	if err := h.service.DeleteLeavePolicy(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "leave policy deleted")
}

// ── Leave Balance Endpoints ──────────────────────────────────────────────────

func (h *LeaveHandler) SetLeaveBalance(c *gin.Context) {
	var req SetLeaveBalanceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	empID, err := uuid.Parse(req.EmpID)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}
	typeID, err := uuid.Parse(req.LeaveTypeID)
	if err != nil {
		response.BadRequest(c, "invalid leave_type_id")
		return
	}

	lb := &models.LeaveBalance{
		EmpID:       empID,
		LeaveTypeID: typeID,
		Year:        req.Year,
		Total:       req.Total,
	}

	res, err := h.service.SetLeaveBalance(c.Request.Context(), lb)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, res, "leave balance allocated")
}

func (h *LeaveHandler) GetLeaveBalances(c *gin.Context) {
	empIDStr := c.Param("emp_id")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}

	yearStr := c.Query("year")
	year, _ := strconv.Atoi(yearStr)

	balances, err := h.service.GetLeaveBalances(c.Request.Context(), empID, year)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, balances)
}

// ── Leave Application Endpoints ──────────────────────────────────────────────

func (h *LeaveHandler) Apply(c *gin.Context) {
	empIDStr := c.Param("emp_id")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}

	var req ApplyLeaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	typeID, err := uuid.Parse(req.LeaveTypeID)
	if err != nil {
		response.BadRequest(c, "invalid leave_type_id")
		return
	}

	fromDate, err := time.Parse("2006-01-02", req.FromDate)
	if err != nil {
		response.BadRequest(c, "from_date must be YYYY-MM-DD")
		return
	}

	toDate, err := time.Parse("2006-01-02", req.ToDate)
	if err != nil {
		response.BadRequest(c, "to_date must be YYYY-MM-DD")
		return
	}

	app := &models.LeaveApplication{
		EmpID:       empID,
		LeaveTypeID: typeID,
		FromDate:    fromDate,
		ToDate:      toDate,
		Session:     req.Session,
		Reason:      req.Reason,
		DocURL:      req.DocURL,
	}

	created, err := h.service.ApplyLeave(c.Request.Context(), app)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, created)
}

func (h *LeaveHandler) GetLeaveByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid leave ID")
		return
	}

	app, err := h.service.GetLeaveByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "leave application")
		return
	}

	response.OK(c, app)
}

func (h *LeaveHandler) ListByEmp(c *gin.Context) {
	empIDStr := c.Param("emp_id")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}

	apps, err := h.service.ListLeavesByEmp(c.Request.Context(), empID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, apps)
}

func (h *LeaveHandler) ListByOrg(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	status := c.Query("status")
	apps, err := h.service.ListLeavesByOrg(c.Request.Context(), orgID, status)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, apps)
}

func (h *LeaveHandler) Approve(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid leave ID")
		return
	}

	userIDStr := c.GetString("user_id")
	approverID, _ := uuid.Parse(userIDStr)

	var req ActionLeaveReq
	_ = c.ShouldBindJSON(&req)

	if err := h.service.ApproveLeave(c.Request.Context(), id, approverID, req.Comment); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, nil, "leave application approved successfully")
}

func (h *LeaveHandler) Reject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid leave ID")
		return
	}

	userIDStr := c.GetString("user_id")
	approverID, _ := uuid.Parse(userIDStr)

	var req ActionLeaveReq
	_ = c.ShouldBindJSON(&req)

	if err := h.service.RejectLeave(c.Request.Context(), id, approverID, req.Comment); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, nil, "leave application rejected")
}

func (h *LeaveHandler) Cancel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid leave ID")
		return
	}

	empIDStr := c.Query("emp_id")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid emp_id query param")
		return
	}

	if err := h.service.CancelLeave(c.Request.Context(), id, empID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, nil, "leave application cancelled")
}
