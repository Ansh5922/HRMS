package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/repository"
)

type PayrollService interface {
	// Salary Structures
	CreateSalaryStructure(ctx context.Context, s *models.SalaryStructure) (*models.SalaryStructure, error)
	GetSalaryStructure(ctx context.Context, id uuid.UUID) (*models.SalaryStructure, error)
	ListSalaryStructures(ctx context.Context, orgID uuid.UUID) ([]*models.SalaryStructure, error)
	DeleteSalaryStructure(ctx context.Context, id uuid.UUID) error

	// Employee Salaries
	AssignEmployeeSalary(ctx context.Context, sal *models.EmployeeSalary) (*models.EmployeeSalary, error)
	GetEmployeeSalary(ctx context.Context, empID uuid.UUID) (*models.EmployeeSalary, error)

	// Payroll Execution & Processing
	CreatePayrollRun(ctx context.Context, run *models.PayrollRun) (*models.PayrollRun, error)
	ProcessPayrollRun(ctx context.Context, runID uuid.UUID, empRepo repository.EmployeeRepository) (*models.PayrollRun, error)
	GetPayrollRun(ctx context.Context, id uuid.UUID) (*models.PayrollRun, error)
	ListPayrollRuns(ctx context.Context, orgID uuid.UUID) ([]*models.PayrollRun, error)

	// Payslips
	GetPayslipByID(ctx context.Context, id uuid.UUID) (*models.Payslip, error)
	ListPayslipsByRun(ctx context.Context, runID uuid.UUID) ([]*models.Payslip, error)
	ListPayslipsByEmp(ctx context.Context, empID uuid.UUID) ([]*models.Payslip, error)

	// Reimbursements
	SubmitReimbursement(ctx context.Context, r *models.Reimbursement) (*models.Reimbursement, error)
	ListReimbursementsByEmp(ctx context.Context, empID uuid.UUID) ([]*models.Reimbursement, error)
	ListReimbursementsByOrg(ctx context.Context, orgID uuid.UUID, status string) ([]*models.Reimbursement, error)
	ApproveReimbursement(ctx context.Context, id, approverID uuid.UUID) error
	RejectReimbursement(ctx context.Context, id, approverID uuid.UUID) error

	// Tax Declarations
	SaveTaxDeclaration(ctx context.Context, t *models.TaxDeclaration) (*models.TaxDeclaration, error)
	GetTaxDeclaration(ctx context.Context, empID uuid.UUID, year int) (*models.TaxDeclaration, error)
}

type payrollService struct {
	repo repository.PayrollRepository
}

func NewPayrollService(repo repository.PayrollRepository) PayrollService {
	return &payrollService{repo: repo}
}

// ── Salary Structures ────────────────────────────────────────────────────────

func (s *payrollService) CreateSalaryStructure(ctx context.Context, structObj *models.SalaryStructure) (*models.SalaryStructure, error) {
	if structObj.Name == "" {
		return nil, errors.New("structure name is required")
	}

	if err := s.repo.CreateStructure(ctx, structObj); err != nil {
		return nil, err
	}

	// Save components if provided
	for i := range structObj.Components {
		comp := &structObj.Components[i]
		comp.StructureID = structObj.ID
		_ = s.repo.AddComponent(ctx, comp)
	}

	return s.repo.GetStructureByID(ctx, structObj.ID)
}

func (s *payrollService) GetSalaryStructure(ctx context.Context, id uuid.UUID) (*models.SalaryStructure, error) {
	st, err := s.repo.GetStructureByID(ctx, id)
	if err != nil || st == nil {
		return nil, errors.New("salary structure not found")
	}
	return st, nil
}

func (s *payrollService) ListSalaryStructures(ctx context.Context, orgID uuid.UUID) ([]*models.SalaryStructure, error) {
	return s.repo.ListStructuresByOrg(ctx, orgID)
}

func (s *payrollService) DeleteSalaryStructure(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteStructure(ctx, id)
}

// ── Employee Salaries ────────────────────────────────────────────────────────

func (s *payrollService) AssignEmployeeSalary(ctx context.Context, sal *models.EmployeeSalary) (*models.EmployeeSalary, error) {
	if sal.CTC <= 0 {
		return nil, errors.New("CTC must be greater than 0")
	}
	if sal.Basic <= 0 {
		sal.Basic = sal.CTC * 0.50 // Default basic to 50% of CTC
	}

	if err := s.repo.AssignSalary(ctx, sal); err != nil {
		return nil, err
	}
	return sal, nil
}

func (s *payrollService) GetEmployeeSalary(ctx context.Context, empID uuid.UUID) (*models.EmployeeSalary, error) {
	sal, err := s.repo.GetSalaryByEmpID(ctx, empID)
	if err != nil || sal == nil {
		return nil, errors.New("no salary structure assigned to employee")
	}
	return sal, nil
}

// ── Payroll Execution & Processing ──────────────────────────────────────────

func (s *payrollService) CreatePayrollRun(ctx context.Context, run *models.PayrollRun) (*models.PayrollRun, error) {
	if run.Month < 1 || run.Month > 12 || run.Year < 2000 {
		return nil, errors.New("invalid month or year")
	}
	run.Status = "draft"
	if err := s.repo.CreatePayrollRun(ctx, run); err != nil {
		return nil, err
	}
	return run, nil
}

func (s *payrollService) ProcessPayrollRun(ctx context.Context, runID uuid.UUID, empRepo repository.EmployeeRepository) (*models.PayrollRun, error) {
	run, err := s.repo.GetPayrollRunByID(ctx, runID)
	if err != nil || run == nil {
		return nil, errors.New("payroll run not found")
	}

	// Fetch all employees in organization
	employees, _, err := empRepo.ListByOrg(ctx, run.OrgID, "", nil, nil, "active", 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch employees for payroll: %w", err)
	}

	var totalOrgGross, totalOrgDeductions, totalOrgNet float64

	for _, emp := range employees {
		sal, err := s.repo.GetSalaryByEmpID(ctx, emp.ID)
		if err != nil || sal == nil {
			continue // Skip employee if salary structure is not defined
		}

		st, _ := s.repo.GetStructureByID(ctx, sal.StructureID)

		monthlyBasic := sal.Basic / 12.0
		monthlyCTC := sal.CTC / 12.0

		var earningsTotal, deductionsTotal float64
		breakdownMap := make(map[string]float64)

		breakdownMap["Basic Salary"] = monthlyBasic
		earningsTotal += monthlyBasic

		if st != nil {
			for _, comp := range st.Components {
				var compAmount float64
				if comp.Value != nil {
					compAmount = *comp.Value / 12.0
				} else if comp.CalcType == "pct_of_basic" {
					compAmount = monthlyBasic * 0.40 // Example 40% HRA
				} else {
					compAmount = 1500.0
				}

				if comp.Type == "earning" {
					breakdownMap[comp.Name] = compAmount
					earningsTotal += compAmount
				} else {
					breakdownMap[comp.Name] = compAmount
					deductionsTotal += compAmount
				}
			}
		}

		// Ensure total gross matches monthly CTC
		if earningsTotal < monthlyCTC {
			allowance := monthlyCTC - earningsTotal
			breakdownMap["Special Allowance"] = allowance
			earningsTotal += allowance
		}

		netPay := earningsTotal - deductionsTotal
		breakdownJSON, _ := json.Marshal(breakdownMap)
		bStr := string(breakdownJSON)

		payslip := &models.Payslip{
			EmpID:           emp.ID,
			PayrollRunID:    runID,
			WorkingDays:     30,
			LopDays:         0,
			PaidDays:        30,
			Gross:           earningsTotal,
			TotalDeductions: deductionsTotal,
			NetPay:          netPay,
			Breakdown:       &bStr,
		}

		_ = s.repo.CreatePayslip(ctx, payslip)

		totalOrgGross += earningsTotal
		totalOrgDeductions += deductionsTotal
		totalOrgNet += netPay
	}

	// Update run status & totals
	if err := s.repo.UpdateRunTotalsAndStatus(ctx, runID, totalOrgGross, totalOrgDeductions, totalOrgNet, "approved"); err != nil {
		return nil, err
	}

	return s.repo.GetPayrollRunByID(ctx, runID)
}

func (s *payrollService) GetPayrollRun(ctx context.Context, id uuid.UUID) (*models.PayrollRun, error) {
	run, err := s.repo.GetPayrollRunByID(ctx, id)
	if err != nil || run == nil {
		return nil, errors.New("payroll run not found")
	}
	return run, nil
}

func (s *payrollService) ListPayrollRuns(ctx context.Context, orgID uuid.UUID) ([]*models.PayrollRun, error) {
	return s.repo.ListPayrollRunsByOrg(ctx, orgID)
}

// ── Payslips ─────────────────────────────────────────────────────────────────

func (s *payrollService) GetPayslipByID(ctx context.Context, id uuid.UUID) (*models.Payslip, error) {
	p, err := s.repo.GetPayslipByID(ctx, id)
	if err != nil || p == nil {
		return nil, errors.New("payslip not found")
	}
	return p, nil
}

func (s *payrollService) ListPayslipsByRun(ctx context.Context, runID uuid.UUID) ([]*models.Payslip, error) {
	return s.repo.ListPayslipsByRun(ctx, runID)
}

func (s *payrollService) ListPayslipsByEmp(ctx context.Context, empID uuid.UUID) ([]*models.Payslip, error) {
	return s.repo.ListPayslipsByEmp(ctx, empID)
}

// ── Reimbursements ───────────────────────────────────────────────────────────

func (s *payrollService) SubmitReimbursement(ctx context.Context, rMB *models.Reimbursement) (*models.Reimbursement, error) {
	if rMB.Amount <= 0 {
		return nil, errors.New("reimbursement amount must be greater than 0")
	}
	if rMB.Category == "" {
		return nil, errors.New("category is required")
	}
	rMB.Status = "pending"

	if err := s.repo.CreateReimbursement(ctx, rMB); err != nil {
		return nil, err
	}
	return rMB, nil
}

func (s *payrollService) ListReimbursementsByEmp(ctx context.Context, empID uuid.UUID) ([]*models.Reimbursement, error) {
	return s.repo.ListReimbursementsByEmp(ctx, empID)
}

func (s *payrollService) ListReimbursementsByOrg(ctx context.Context, orgID uuid.UUID, status string) ([]*models.Reimbursement, error) {
	return s.repo.ListReimbursementsByOrg(ctx, orgID, status)
}

func (s *payrollService) ApproveReimbursement(ctx context.Context, id, approverID uuid.UUID) error {
	return s.repo.UpdateReimbursementStatus(ctx, id, "approved", &approverID)
}

func (s *payrollService) RejectReimbursement(ctx context.Context, id, approverID uuid.UUID) error {
	return s.repo.UpdateReimbursementStatus(ctx, id, "rejected", &approverID)
}

// ── Tax Declarations ─────────────────────────────────────────────────────────

func (s *payrollService) SaveTaxDeclaration(ctx context.Context, t *models.TaxDeclaration) (*models.TaxDeclaration, error) {
	if t.Regime != "old" && t.Regime != "new" {
		t.Regime = "new"
	}
	if err := s.repo.SaveTaxDeclaration(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *payrollService) GetTaxDeclaration(ctx context.Context, empID uuid.UUID, year int) (*models.TaxDeclaration, error) {
	t, err := s.repo.GetTaxDeclaration(ctx, empID, year)
	if err != nil || t == nil {
		return nil, errors.New("tax declaration not found for specified year")
	}
	return t, nil
}