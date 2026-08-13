package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkflowTemplate struct {
	ID        uuid.UUID  `json:"id"`
	OrgID     uuid.UUID  `json:"org_id"`
	Name      string     `json:"name"`
	Module    string     `json:"module"` // leave, expense, document, shift
	Steps     *string    `json:"steps"`  // JSON string array of step definitions
	IsActive  bool       `json:"is_active"`
	CreatedBy *uuid.UUID `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
}

type ApprovalRequest struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       uuid.UUID  `json:"org_id"`
	WorkflowID  *uuid.UUID `json:"workflow_id"`
	Module      string     `json:"module"`
	ResourceID  uuid.UUID  `json:"resource_id"`
	RequestedBy uuid.UUID  `json:"requested_by"`
	CurrentStep int        `json:"current_step"`
	Status      string     `json:"status"` // pending, approved, rejected, cancelled
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Joined
	WorkflowName   *string         `json:"workflow_name,omitempty"`
	RequestedByName *string        `json:"requested_by_name,omitempty"`
	Steps          []*ApprovalStep `json:"steps,omitempty"`
}

type ApprovalStep struct {
	ID                uuid.UUID  `json:"id"`
	ApprovalRequestID uuid.UUID  `json:"approval_request_id"`
	StepNo            int        `json:"step_no"`
	ApproverID        *uuid.UUID `json:"approver_id"`
	Action            *string    `json:"action"` // approved, rejected
	Comment           *string    `json:"comment"`
	ActionedAt        *time.Time `json:"actioned_at"`

	// Joined
	ApproverName *string `json:"approver_name,omitempty"`
}
