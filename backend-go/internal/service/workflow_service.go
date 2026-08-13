package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/repository"
)

type WorkflowService interface {
	// Templates
	CreateTemplate(ctx context.Context, wt *models.WorkflowTemplate) (*models.WorkflowTemplate, error)
	GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.WorkflowTemplate, error)
	ListTemplates(ctx context.Context, orgID uuid.UUID) ([]*models.WorkflowTemplate, error)
	UpdateTemplate(ctx context.Context, wt *models.WorkflowTemplate) error
	DeleteTemplate(ctx context.Context, id uuid.UUID) error

	// Approval Engine
	SubmitApprovalRequest(ctx context.Context, ar *models.ApprovalRequest) (*models.ApprovalRequest, error)
	GetApprovalRequest(ctx context.Context, id uuid.UUID) (*models.ApprovalRequest, error)
	ListOrgRequests(ctx context.Context, orgID uuid.UUID, status string) ([]*models.ApprovalRequest, error)
	ListMyPendingApprovals(ctx context.Context, userID uuid.UUID) ([]*models.ApprovalRequest, error)
	ListMySubmittedRequests(ctx context.Context, userID uuid.UUID) ([]*models.ApprovalRequest, error)

	// Step Action (Approve / Reject)
	ProcessStepAction(ctx context.Context, reqID, approverID uuid.UUID, action string, comment *string) error
}

type workflowService struct {
	repo repository.WorkflowRepository
}

func NewWorkflowService(repo repository.WorkflowRepository) WorkflowService {
	return &workflowService{repo: repo}
}

// ── Templates ────────────────────────────────────────────────────────────────

func (s *workflowService) CreateTemplate(ctx context.Context, wt *models.WorkflowTemplate) (*models.WorkflowTemplate, error) {
	if wt.Name == "" || wt.Module == "" {
		return nil, errors.New("template name and module are required")
	}
	wt.IsActive = true
	if err := s.repo.CreateTemplate(ctx, wt); err != nil {
		return nil, err
	}
	return wt, nil
}

func (s *workflowService) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.WorkflowTemplate, error) {
	wt, err := s.repo.GetTemplateByID(ctx, id)
	if err != nil || wt == nil {
		return nil, errors.New("workflow template not found")
	}
	return wt, nil
}

func (s *workflowService) ListTemplates(ctx context.Context, orgID uuid.UUID) ([]*models.WorkflowTemplate, error) {
	return s.repo.ListTemplatesByOrg(ctx, orgID)
}

func (s *workflowService) UpdateTemplate(ctx context.Context, wt *models.WorkflowTemplate) error {
	return s.repo.UpdateTemplate(ctx, wt)
}

func (s *workflowService) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTemplate(ctx, id)
}

// ── Approval Engine ──────────────────────────────────────────────────────────

func (s *workflowService) SubmitApprovalRequest(ctx context.Context, ar *models.ApprovalRequest) (*models.ApprovalRequest, error) {
	if ar.Module == "" || ar.ResourceID == uuid.Nil {
		return nil, errors.New("module and resource_id are required")
	}
	ar.CurrentStep = 1
	ar.Status = "pending"

	if err := s.repo.CreateApprovalRequest(ctx, ar); err != nil {
		return nil, err
	}
	return ar, nil
}

func (s *workflowService) GetApprovalRequest(ctx context.Context, id uuid.UUID) (*models.ApprovalRequest, error) {
	ar, err := s.repo.GetApprovalRequestByID(ctx, id)
	if err != nil || ar == nil {
		return nil, errors.New("approval request not found")
	}
	return ar, nil
}

func (s *workflowService) ListOrgRequests(ctx context.Context, orgID uuid.UUID, status string) ([]*models.ApprovalRequest, error) {
	return s.repo.ListApprovalRequestsByOrg(ctx, orgID, status)
}

func (s *workflowService) ListMyPendingApprovals(ctx context.Context, userID uuid.UUID) ([]*models.ApprovalRequest, error) {
	return s.repo.ListPendingApprovalsForUser(ctx, userID)
}

func (s *workflowService) ListMySubmittedRequests(ctx context.Context, userID uuid.UUID) ([]*models.ApprovalRequest, error) {
	return s.repo.ListSubmittedRequestsForUser(ctx, userID)
}

func (s *workflowService) ProcessStepAction(ctx context.Context, reqID, approverID uuid.UUID, action string, comment *string) error {
	if action != "approved" && action != "rejected" {
		return errors.New("action must be 'approved' or 'rejected'")
	}

	ar, err := s.repo.GetApprovalRequestByID(ctx, reqID)
	if err != nil || ar == nil {
		return errors.New("approval request not found")
	}
	if ar.Status != "pending" {
		return fmt.Errorf("approval request is already '%s'", ar.Status)
	}

	// Record step action
	step := &models.ApprovalStep{
		ApprovalRequestID: reqID,
		StepNo:            ar.CurrentStep,
		ApproverID:        &approverID,
		Action:            &action,
		Comment:           comment,
	}
	if err := s.repo.AddApprovalStep(ctx, step); err != nil {
		return err
	}

	if action == "rejected" {
		// Mark whole request as rejected
		return s.repo.UpdateApprovalRequestStepAndStatus(ctx, reqID, ar.CurrentStep, "rejected")
	}

	// If approved, check if there are further steps (assuming 3-step max by default)
	maxSteps := 3
	if ar.CurrentStep >= maxSteps {
		// All steps complete -> Approved
		return s.repo.UpdateApprovalRequestStepAndStatus(ctx, reqID, ar.CurrentStep, "approved")
	}

	// Advance to next step
	nextStep := ar.CurrentStep + 1
	return s.repo.UpdateApprovalRequestStepAndStatus(ctx, reqID, nextStep, "pending")
}
