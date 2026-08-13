package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/hrms-backend/internal/models"
)

type WorkflowRepository interface {
	// Workflow Templates
	CreateTemplate(ctx context.Context, wt *models.WorkflowTemplate) error
	GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.WorkflowTemplate, error)
	ListTemplatesByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.WorkflowTemplate, error)
	UpdateTemplate(ctx context.Context, wt *models.WorkflowTemplate) error
	DeleteTemplate(ctx context.Context, id uuid.UUID) error

	// Approval Requests
	CreateApprovalRequest(ctx context.Context, ar *models.ApprovalRequest) error
	GetApprovalRequestByID(ctx context.Context, id uuid.UUID) (*models.ApprovalRequest, error)
	ListApprovalRequestsByOrg(ctx context.Context, orgID uuid.UUID, status string) ([]*models.ApprovalRequest, error)
	ListPendingApprovalsForUser(ctx context.Context, userID uuid.UUID) ([]*models.ApprovalRequest, error)
	ListSubmittedRequestsForUser(ctx context.Context, userID uuid.UUID) ([]*models.ApprovalRequest, error)
	UpdateApprovalRequestStepAndStatus(ctx context.Context, id uuid.UUID, step int, status string) error

	// Approval Steps
	AddApprovalStep(ctx context.Context, step *models.ApprovalStep) error
	ListStepsByRequestID(ctx context.Context, reqID uuid.UUID) ([]*models.ApprovalStep, error)
}

type workflowRepo struct {
	db *pgxpool.Pool
}

func NewWorkflowRepository(db *pgxpool.Pool) WorkflowRepository {
	return &workflowRepo{db: db}
}

// ── Workflow Templates ───────────────────────────────────────────────────────

func (r *workflowRepo) CreateTemplate(ctx context.Context, wt *models.WorkflowTemplate) error {
	query := `INSERT INTO workflow_templates (org_id, name, module, steps, is_active, created_by)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, wt.OrgID, wt.Name, wt.Module, wt.Steps, wt.IsActive, wt.CreatedBy).
		Scan(&wt.ID, &wt.CreatedAt)
}

func (r *workflowRepo) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.WorkflowTemplate, error) {
	query := `SELECT id, org_id, name, module, steps, is_active, created_by, created_at
		FROM workflow_templates WHERE id = $1`
	wt := &models.WorkflowTemplate{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&wt.ID, &wt.OrgID, &wt.Name, &wt.Module, &wt.Steps, &wt.IsActive, &wt.CreatedBy, &wt.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return wt, err
}

func (r *workflowRepo) ListTemplatesByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.WorkflowTemplate, error) {
	query := `SELECT id, org_id, name, module, steps, is_active, created_by, created_at
		FROM workflow_templates WHERE org_id = $1 ORDER BY name`
	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*models.WorkflowTemplate
	for rows.Next() {
		wt := &models.WorkflowTemplate{}
		if err := rows.Scan(&wt.ID, &wt.OrgID, &wt.Name, &wt.Module, &wt.Steps, &wt.IsActive, &wt.CreatedBy, &wt.CreatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, wt)
	}
	return templates, nil
}

func (r *workflowRepo) UpdateTemplate(ctx context.Context, wt *models.WorkflowTemplate) error {
	query := `UPDATE workflow_templates SET name=$1, module=$2, steps=$3, is_active=$4 WHERE id=$5 AND org_id=$6`
	_, err := r.db.Exec(ctx, query, wt.Name, wt.Module, wt.Steps, wt.IsActive, wt.ID, wt.OrgID)
	return err
}

func (r *workflowRepo) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM workflow_templates WHERE id=$1`, id)
	return err
}

// ── Approval Requests ────────────────────────────────────────────────────────

func (r *workflowRepo) CreateApprovalRequest(ctx context.Context, ar *models.ApprovalRequest) error {
	query := `INSERT INTO approval_requests (org_id, workflow_id, module, resource_id, requested_by, current_step, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, ar.OrgID, ar.WorkflowID, ar.Module, ar.ResourceID, ar.RequestedBy, ar.CurrentStep, ar.Status).
		Scan(&ar.ID, &ar.CreatedAt, &ar.UpdatedAt)
}

func (r *workflowRepo) GetApprovalRequestByID(ctx context.Context, id uuid.UUID) (*models.ApprovalRequest, error) {
	query := `SELECT ar.id, ar.org_id, ar.workflow_id, ar.module, ar.resource_id, ar.requested_by, ar.current_step, ar.status, ar.created_at, ar.updated_at,
		wt.name as workflow_name, u.email as requested_by_name
	FROM approval_requests ar
	LEFT JOIN workflow_templates wt ON wt.id = ar.workflow_id
	LEFT JOIN users u ON u.id = ar.requested_by
	WHERE ar.id = $1`

	ar := &models.ApprovalRequest{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&ar.ID, &ar.OrgID, &ar.WorkflowID, &ar.Module, &ar.ResourceID, &ar.RequestedBy, &ar.CurrentStep, &ar.Status, &ar.CreatedAt, &ar.UpdatedAt,
		&ar.WorkflowName, &ar.RequestedByName,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Fetch steps
	steps, _ := r.ListStepsByRequestID(ctx, id)
	ar.Steps = steps
	return ar, nil
}

func (r *workflowRepo) ListApprovalRequestsByOrg(ctx context.Context, orgID uuid.UUID, status string) ([]*models.ApprovalRequest, error) {
	query := `SELECT ar.id, ar.org_id, ar.workflow_id, ar.module, ar.resource_id, ar.requested_by, ar.current_step, ar.status, ar.created_at, ar.updated_at,
		wt.name as workflow_name, u.email as requested_by_name
	FROM approval_requests ar
	LEFT JOIN workflow_templates wt ON wt.id = ar.workflow_id
	LEFT JOIN users u ON u.id = ar.requested_by
	WHERE ar.org_id = $1`

	args := []interface{}{orgID}
	if status != "" {
		query += ` AND ar.status = $2`
		args = append(args, status)
	}
	query += ` ORDER BY ar.created_at DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.ApprovalRequest
	for rows.Next() {
		ar := &models.ApprovalRequest{}
		if err := rows.Scan(
			&ar.ID, &ar.OrgID, &ar.WorkflowID, &ar.Module, &ar.ResourceID, &ar.RequestedBy, &ar.CurrentStep, &ar.Status, &ar.CreatedAt, &ar.UpdatedAt,
			&ar.WorkflowName, &ar.RequestedByName,
		); err != nil {
			return nil, err
		}
		list = append(list, ar)
	}
	return list, nil
}

func (r *workflowRepo) ListPendingApprovalsForUser(ctx context.Context, userID uuid.UUID) ([]*models.ApprovalRequest, error) {
	query := `SELECT ar.id, ar.org_id, ar.workflow_id, ar.module, ar.resource_id, ar.requested_by, ar.current_step, ar.status, ar.created_at, ar.updated_at,
		wt.name as workflow_name, u.email as requested_by_name
	FROM approval_requests ar
	LEFT JOIN workflow_templates wt ON wt.id = ar.workflow_id
	LEFT JOIN users u ON u.id = ar.requested_by
	WHERE ar.status = 'pending' ORDER BY ar.created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.ApprovalRequest
	for rows.Next() {
		ar := &models.ApprovalRequest{}
		if err := rows.Scan(
			&ar.ID, &ar.OrgID, &ar.WorkflowID, &ar.Module, &ar.ResourceID, &ar.RequestedBy, &ar.CurrentStep, &ar.Status, &ar.CreatedAt, &ar.UpdatedAt,
			&ar.WorkflowName, &ar.RequestedByName,
		); err != nil {
			return nil, err
		}
		list = append(list, ar)
	}
	return list, nil
}

func (r *workflowRepo) ListSubmittedRequestsForUser(ctx context.Context, userID uuid.UUID) ([]*models.ApprovalRequest, error) {
	query := `SELECT ar.id, ar.org_id, ar.workflow_id, ar.module, ar.resource_id, ar.requested_by, ar.current_step, ar.status, ar.created_at, ar.updated_at,
		wt.name as workflow_name
	FROM approval_requests ar
	LEFT JOIN workflow_templates wt ON wt.id = ar.workflow_id
	WHERE ar.requested_by = $1 ORDER BY ar.created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.ApprovalRequest
	for rows.Next() {
		ar := &models.ApprovalRequest{}
		if err := rows.Scan(
			&ar.ID, &ar.OrgID, &ar.WorkflowID, &ar.Module, &ar.ResourceID, &ar.RequestedBy, &ar.CurrentStep, &ar.Status, &ar.CreatedAt, &ar.UpdatedAt,
			&ar.WorkflowName,
		); err != nil {
			return nil, err
		}
		list = append(list, ar)
	}
	return list, nil
}

func (r *workflowRepo) UpdateApprovalRequestStepAndStatus(ctx context.Context, id uuid.UUID, step int, status string) error {
	query := `UPDATE approval_requests SET current_step=$1, status=$2, updated_at=NOW() WHERE id=$3`
	_, err := r.db.Exec(ctx, query, step, status, id)
	return err
}

// ── Approval Steps ───────────────────────────────────────────────────────────

func (r *workflowRepo) AddApprovalStep(ctx context.Context, step *models.ApprovalStep) error {
	query := `INSERT INTO approval_steps (approval_request_id, step_no, approver_id, action, comment, actioned_at)
		VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`
	return r.db.QueryRow(ctx, query, step.ApprovalRequestID, step.StepNo, step.ApproverID, step.Action, step.Comment).
		Scan(&step.ID)
}

func (r *workflowRepo) ListStepsByRequestID(ctx context.Context, reqID uuid.UUID) ([]*models.ApprovalStep, error) {
	query := `SELECT st.id, st.approval_request_id, st.step_no, st.approver_id, st.action, st.comment, st.actioned_at,
		u.email as approver_name
	FROM approval_steps st
	LEFT JOIN users u ON u.id = st.approver_id
	WHERE st.approval_request_id = $1 ORDER BY st.step_no ASC`

	rows, err := r.db.Query(ctx, query, reqID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []*models.ApprovalStep
	for rows.Next() {
		st := &models.ApprovalStep{}
		if err := rows.Scan(&st.ID, &st.ApprovalRequestID, &st.StepNo, &st.ApproverID, &st.Action, &st.Comment, &st.ActionedAt, &st.ApproverName); err != nil {
			return nil, err
		}
		steps = append(steps, st)
	}
	return steps, nil
}
