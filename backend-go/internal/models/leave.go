package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaveType struct {
	ID               uuid.UUID `json:"id"`
	OrgID            uuid.UUID `json:"org_id"`
	Name             string    `json:"name"`
	Code             string    `json:"code"`
	MaxDaysPerYear   *int      `json:"max_days_per_year"`
	CarryForward     bool      `json:"carry_forward"`
	CarryForwardMax  *int      `json:"carry_forward_max"`
	IsPaid           bool      `json:"is_paid"`
	RequiresDoc      bool      `json:"requires_doc"`
	MinNoticeDays    int       `json:"min_notice_days"`
	GenderApplicable string    `json:"gender_applicable"`
	RequiresApproval bool      `json:"requires_approval"`
	CreatedAt        time.Time `json:"created_at"`
}

type LeavePolicy struct {
	ID            uuid.UUID  `json:"id"`
	OrgID         uuid.UUID  `json:"org_id"`
	LeaveTypeID   uuid.UUID  `json:"leave_type_id"`
	DesignationID *uuid.UUID `json:"designation_id"`
	DepartmentID  *uuid.UUID `json:"department_id"`
	AnnualQuota   float64    `json:"annual_quota"`
	AccrualType   string     `json:"accrual_type"`
	CreatedAt     time.Time  `json:"created_at"`

	// Joined
	LeaveTypeName *string `json:"leave_type_name,omitempty"`
}

type LeaveBalance struct {
	ID          uuid.UUID `json:"id"`
	EmpID       uuid.UUID `json:"emp_id"`
	LeaveTypeID uuid.UUID `json:"leave_type_id"`
	Year        int       `json:"year"`
	Total       float64   `json:"total"`
	Used        float64   `json:"used"`
	Pending     float64   `json:"pending"`
	Available   float64   `json:"available"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Joined
	LeaveTypeName *string `json:"leave_type_name,omitempty"`
	LeaveTypeCode *string `json:"leave_type_code,omitempty"`
}

type LeaveApplication struct {
	ID          uuid.UUID `json:"id"`
	EmpID       uuid.UUID `json:"emp_id"`
	LeaveTypeID uuid.UUID `json:"leave_type_id"`
	FromDate    time.Time `json:"from_date"`
	ToDate      time.Time `json:"to_date"`
	Days        float64   `json:"days"`
	Session     string    `json:"session"`
	Reason      *string   `json:"reason"`
	DocURL      *string   `json:"doc_url"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Joined
	EmployeeName  *string `json:"employee_name,omitempty"`
	EmpCode       *string `json:"emp_code,omitempty"`
	LeaveTypeName *string `json:"leave_type_name,omitempty"`
	LeaveTypeCode *string `json:"leave_type_code,omitempty"`
}

type LeaveApproval struct {
	ID            uuid.UUID  `json:"id"`
	ApplicationID uuid.UUID  `json:"application_id"`
	ApproverID    uuid.UUID  `json:"approver_id"`
	Level         int        `json:"level"`
	Action        string     `json:"action"`
	Comment       *string    `json:"comment"`
	ActionedAt    time.Time  `json:"actioned_at"`

	// Joined
	ApproverName *string `json:"approver_name,omitempty"`
}