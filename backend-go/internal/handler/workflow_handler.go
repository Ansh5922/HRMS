package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/service"
	"github.com/your-org/hrms-backend/pkg/response"
)

type WorkflowHandler struct {
	service service.WorkflowService
}

func NewWorkflowHandler(svc service.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{service: svc}
}

// ── Request DTOs ─────────────────────────────────────────────────────────────

type CreateTemplateReq struct {
	Name   string  `json:"name" binding:"required"`
	Module string  `json:"module" binding:"required"`
	Steps  *string `json:"steps"` // JSON array string
}

type SubmitApprovalRequestReq struct {
	WorkflowID *string `json:"workflow_id"`
	Module     string  `json:"module" binding:"required"`
	ResourceID string  `json:"resource_id" binding:"required"`
}

type ProcessStepActionReq struct {
	Comment *string `json:"comment"`
}

// ── Template Endpoints ───────────────────────────────────────────────────────

func (h *WorkflowHandler) CreateTemplate(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	var req CreateTemplateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	wt := &models.WorkflowTemplate{
		OrgID:     orgID,
		Name:      req.Name,
		Module:    req.Module,
		Steps:     req.Steps,
		CreatedBy: &userID,
	}

	res, err := h.service.CreateTemplate(c.Request.Context(), wt)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *WorkflowHandler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid template ID")
		return
	}

	wt, err := h.service.GetTemplateByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "workflow template")
		return
	}

	response.OK(c, wt)
}

func (h *WorkflowHandler) ListTemplates(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	templates, err := h.service.ListTemplates(c.Request.Context(), orgID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, templates)
}

func (h *WorkflowHandler) UpdateTemplate(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid template ID")
		return
	}

	var wt models.WorkflowTemplate
	if err := c.ShouldBindJSON(&wt); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	wt.ID = id
	wt.OrgID = orgID

	if err := h.service.UpdateTemplate(c.Request.Context(), &wt); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "workflow template updated")
}

func (h *WorkflowHandler) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid template ID")
		return
	}

	if err := h.service.DeleteTemplate(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "workflow template deleted")
}

// ── Approval Request Endpoints ───────────────────────────────────────────────

func (h *WorkflowHandler) SubmitRequest(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	var req SubmitApprovalRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resourceID, err := uuid.Parse(req.ResourceID)
	if err != nil {
		response.BadRequest(c, "invalid resource_id")
		return
	}

	ar := &models.ApprovalRequest{
		OrgID:       orgID,
		Module:      req.Module,
		ResourceID:  resourceID,
		RequestedBy: userID,
	}

	if req.WorkflowID != nil {
		if wID, err := uuid.Parse(*req.WorkflowID); err == nil {
			ar.WorkflowID = &wID
		}
	}

	res, err := h.service.SubmitApprovalRequest(c.Request.Context(), ar)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *WorkflowHandler) GetRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid request ID")
		return
	}

	ar, err := h.service.GetApprovalRequest(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "approval request")
		return
	}

	response.OK(c, ar)
}

func (h *WorkflowHandler) ListOrgRequests(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	status := c.Query("status")
	list, err := h.service.ListOrgRequests(c.Request.Context(), orgID, status)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, list)
}

func (h *WorkflowHandler) ListMyPendingApprovals(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	list, err := h.service.ListMyPendingApprovals(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, list)
}

func (h *WorkflowHandler) ListMySubmittedRequests(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	list, err := h.service.ListMySubmittedRequests(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, list)
}

func (h *WorkflowHandler) ApproveStep(c *gin.Context) {
	idStr := c.Param("id")
	reqID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid request ID")
		return
	}

	userIDStr := c.GetString("user_id")
	approverID, _ := uuid.Parse(userIDStr)

	var req ProcessStepActionReq
	_ = c.ShouldBindJSON(&req)

	if err := h.service.ProcessStepAction(c.Request.Context(), reqID, approverID, "approved", req.Comment); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, nil, "approval step approved")
}

func (h *WorkflowHandler) RejectStep(c *gin.Context) {
	idStr := c.Param("id")
	reqID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid request ID")
		return
	}

	userIDStr := c.GetString("user_id")
	approverID, _ := uuid.Parse(userIDStr)

	var req ProcessStepActionReq
	_ = c.ShouldBindJSON(&req)

	if err := h.service.ProcessStepAction(c.Request.Context(), reqID, approverID, "rejected", req.Comment); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, nil, "approval step rejected")
}
