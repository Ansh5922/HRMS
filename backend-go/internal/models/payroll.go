package models

import (
	"time"

	"github.com/google/uuid"
)

type SalaryStructure struct {
	ID          uuid.UUID         `json:"id"`
	OrgID       uuid.UUID         `json:"org_id"`
	Name        string            `json:"name"`
	Description *string           `json:"description"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	Components  []SalaryComponent `json:"components,omitempty"`
}

type SalaryComponent struct {
	ID           uuid.UUID `json:"id"`
	StructureID  uuid.UUID `json:"structure_id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Type         string    `json:"type"`      // earning, deduction, stat_deduction
	CalcType     string    `json:"calc_type"` // fixed, pct_of_basic, pct_of_gross, formula
	Value        *float64  `json:"value"`
	Formula      *string   `json:"formula"`
	IsTaxable    bool      `json:"is_taxable"`
	DisplayOrder int       `json:"display_order"`
}

type EmployeeSalary struct {
	ID            uuid.UUID  `json:"id"`
	EmpID         uuid.UUID  `json:"emp_id"`
	StructureID   uuid.UUID  `json:"structure_id"`
	CTC           float64    `json:"ctc"`
	Basic         float64    `json:"basic"`
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`
	ApprovedBy    *uuid.UUID `json:"approved_by"`
	ApprovedAt    *time.Time `json:"approved_at"`
	CreatedAt     time.Time  `json:"created_at"`

	// Joined
	StructureName *string `json:"structure_name,omitempty"`
	EmployeeName  *string `json:"employee_name,omitempty"`
	EmpCode       *string `json:"emp_code,omitempty"`
}

type PayrollRun struct {
	ID              uuid.UUID  `json:"id"`
	OrgID           uuid.UUID  `json:"org_id"`
	Month           int        `json:"month"`
	Year            int        `json:"year"`
	Status          string     `json:"status"` // draft, processing, pending_approval, approved, disbursed
	TotalGross      float64    `json:"total_gross"`
	TotalDeductions float64    `json:"total_deductions"`
	TotalNet        float64    `json:"total_net"`
	RunBy           *uuid.UUID `json:"run_by"`
	ApprovedBy      *uuid.UUID `json:"approved_by"`
	ApprovedAt      *time.Time `json:"approved_at"`
	DisbursedAt     *time.Time `json:"disbursed_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type Payslip struct {
	ID              uuid.UUID `json:"id"`
	EmpID           uuid.UUID `json:"emp_id"`
	PayrollRunID    uuid.UUID `json:"payroll_run_id"`
	WorkingDays     int       `json:"working_days"`
	LopDays         float64   `json:"lop_days"`
	PaidDays        float64   `json:"paid_days"`
	Gross           float64   `json:"gross"`
	TotalDeductions float64   `json:"total_deductions"`
	NetPay          float64   `json:"net_pay"`
	Breakdown       *string   `json:"breakdown"` // JSON string of earnings and deductions
	PayslipURL      *string   `json:"payslip_url"`
	GeneratedAt     time.Time `json:"generated_at"`

	// Joined
	EmployeeName *string `json:"employee_name,omitempty"`
	EmpCode      *string `json:"emp_code,omitempty"`
}

type Reimbursement struct {
	ID           uuid.UUID  `json:"id"`
	EmpID        uuid.UUID  `json:"emp_id"`
	Category     string     `json:"category"`
	Amount       float64    `json:"amount"`
	Description  *string    `json:"description"`
	ReceiptURL   *string    `json:"receipt_url"`
	Status       string     `json:"status"` // pending, approved, rejected, paid
	ApprovedBy   *uuid.UUID `json:"approved_by"`
	PayrollRunID *uuid.UUID `json:"payroll_run_id"`
	CreatedAt    time.Time  `json:"created_at"`

	// Joined
	EmployeeName *string `json:"employee_name,omitempty"`
}

type TaxDeclaration struct {
	ID          uuid.UUID `json:"id"`
	EmpID       uuid.UUID `json:"emp_id"`
	Year        int       `json:"year"`
	Regime      string    `json:"regime"` // old, new
	Sections    *string   `json:"sections"`
	SubmittedAt time.Time `json:"submitted_at"`
}