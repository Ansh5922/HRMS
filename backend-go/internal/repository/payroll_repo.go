package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/hrms-backend/internal/models"
)

type PayrollRepository interface {
	// Structures & Components
	CreateStructure(ctx context.Context, s *models.SalaryStructure) error
	AddComponent(ctx context.Context, c *models.SalaryComponent) error
	GetStructureByID(ctx context.Context, id uuid.UUID) (*models.SalaryStructure, error)
	ListStructuresByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.SalaryStructure, error)
	DeleteStructure(ctx context.Context, id uuid.UUID) error

	// Employee Salaries
	AssignSalary(ctx context.Context, sal *models.EmployeeSalary) error
	GetSalaryByEmpID(ctx context.Context, empID uuid.UUID) (*models.EmployeeSalary, error)

	// Payroll Runs & Processing
	CreatePayrollRun(ctx context.Context, run *models.PayrollRun) error
	GetPayrollRunByID(ctx context.Context, id uuid.UUID) (*models.PayrollRun, error)
	ListPayrollRunsByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.PayrollRun, error)
	UpdateRunTotalsAndStatus(ctx context.Context, id uuid.UUID, gross, deductions, net float64, status string) error

	// Payslips
	CreatePayslip(ctx context.Context, p *models.Payslip) error
	GetPayslipByID(ctx context.Context, id uuid.UUID) (*models.Payslip, error)
	ListPayslipsByRun(ctx context.Context, runID uuid.UUID) ([]*models.Payslip, error)
	ListPayslipsByEmp(ctx context.Context, empID uuid.UUID) ([]*models.Payslip, error)

	// Reimbursements
	CreateReimbursement(ctx context.Context, r *models.Reimbursement) error
	ListReimbursementsByEmp(ctx context.Context, empID uuid.UUID) ([]*models.Reimbursement, error)
	ListReimbursementsByOrg(ctx context.Context, orgID uuid.UUID, status string) ([]*models.Reimbursement, error)
	UpdateReimbursementStatus(ctx context.Context, id uuid.UUID, status string, approverID *uuid.UUID) error

	// Tax Declarations
	SaveTaxDeclaration(ctx context.Context, t *models.TaxDeclaration) error
	GetTaxDeclaration(ctx context.Context, empID uuid.UUID, year int) (*models.TaxDeclaration, error)
}

type payrollRepo struct {
	db *pgxpool.Pool
}

func NewPayrollRepository(db *pgxpool.Pool) PayrollRepository {
	return &payrollRepo{db: db}
}

// ── Structures & Components ──────────────────────────────────────────────────

func (r *payrollRepo) CreateStructure(ctx context.Context, s *models.SalaryStructure) error {
	query := `INSERT INTO salary_structures (org_id, name, description, is_active)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, s.OrgID, s.Name, s.Description, s.IsActive).Scan(&s.ID, &s.CreatedAt)
}

func (r *payrollRepo) AddComponent(ctx context.Context, c *models.SalaryComponent) error {
	query := `INSERT INTO salary_components (structure_id, name, code, type, calc_type, value, formula, is_taxable, display_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	return r.db.QueryRow(ctx, query, c.StructureID, c.Name, c.Code, c.Type, c.CalcType, c.Value, c.Formula, c.IsTaxable, c.DisplayOrder).Scan(&c.ID)
}

func (r *payrollRepo) GetStructureByID(ctx context.Context, id uuid.UUID) (*models.SalaryStructure, error) {
	query := `SELECT id, org_id, name, description, is_active, created_at FROM salary_structures WHERE id = $1`
	s := &models.SalaryStructure{}
	err := r.db.QueryRow(ctx, query, id).Scan(&s.ID, &s.OrgID, &s.Name, &s.Description, &s.IsActive, &s.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Fetch components
	compQuery := `SELECT id, structure_id, name, code, type, calc_type, value, formula, is_taxable, display_order
		FROM salary_components WHERE structure_id = $1 ORDER BY display_order ASC`
	rows, err := r.db.Query(ctx, compQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var c models.SalaryComponent
			if err := rows.Scan(&c.ID, &c.StructureID, &c.Name, &c.Code, &c.Type, &c.CalcType, &c.Value, &c.Formula, &c.IsTaxable, &c.DisplayOrder); err == nil {
				s.Components = append(s.Components, c)
			}
		}
	}
	return s, nil
}

func (r *payrollRepo) ListStructuresByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.SalaryStructure, error) {
	query := `SELECT id, org_id, name, description, is_active, created_at FROM salary_structures WHERE org_id = $1 ORDER BY name`
	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var structs []*models.SalaryStructure
	for rows.Next() {
		s := &models.SalaryStructure{}
		if err := rows.Scan(&s.ID, &s.OrgID, &s.Name, &s.Description, &s.IsActive, &s.CreatedAt); err != nil {
			return nil, err
		}
		structs = append(structs, s)
	}
	return structs, nil
}

func (r *payrollRepo) DeleteStructure(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM salary_structures WHERE id=$1`, id)
	return err
}

// ── Employee Salaries ────────────────────────────────────────────────────────

func (r *payrollRepo) AssignSalary(ctx context.Context, sal *models.EmployeeSalary) error {
	query := `INSERT INTO employee_salaries (emp_id, structure_id, ctc, basic, effective_from)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, sal.EmpID, sal.StructureID, sal.CTC, sal.Basic, sal.EffectiveFrom).
		Scan(&sal.ID, &sal.CreatedAt)
}

func (r *payrollRepo) GetSalaryByEmpID(ctx context.Context, empID uuid.UUID) (*models.EmployeeSalary, error) {
	query := `SELECT es.id, es.emp_id, es.structure_id, es.ctc, es.basic, es.effective_from, es.effective_to, es.approved_by, es.approved_at, es.created_at,
		st.name as structure_name, e.first_name || ' ' || e.last_name as employee_name, e.emp_code
	FROM employee_salaries es
	JOIN salary_structures st ON st.id = es.structure_id
	JOIN employees e ON e.id = es.emp_id
	WHERE es.emp_id = $1 ORDER BY es.effective_from DESC LIMIT 1`

	sal := &models.EmployeeSalary{}
	err := r.db.QueryRow(ctx, query, empID).Scan(
		&sal.ID, &sal.EmpID, &sal.StructureID, &sal.CTC, &sal.Basic, &sal.EffectiveFrom, &sal.EffectiveTo, &sal.ApprovedBy, &sal.ApprovedAt, &sal.CreatedAt,
		&sal.StructureName, &sal.EmployeeName, &sal.EmpCode,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return sal, err
}

// ── Payroll Runs ─────────────────────────────────────────────────────────────

func (r *payrollRepo) CreatePayrollRun(ctx context.Context, run *models.PayrollRun) error {
	query := `INSERT INTO payroll_runs (org_id, month, year, status, run_by)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, run.OrgID, run.Month, run.Year, run.Status, run.RunBy).Scan(&run.ID, &run.CreatedAt)
}

func (r *payrollRepo) GetPayrollRunByID(ctx context.Context, id uuid.UUID) (*models.PayrollRun, error) {
	query := `SELECT id, org_id, month, year, status, COALESCE(total_gross, 0), COALESCE(total_deductions, 0), COALESCE(total_net, 0), run_by, approved_by, approved_at, disbursed_at, created_at
		FROM payroll_runs WHERE id = $1`

	run := &models.PayrollRun{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&run.ID, &run.OrgID, &run.Month, &run.Year, &run.Status, &run.TotalGross, &run.TotalDeductions, &run.TotalNet,
		&run.RunBy, &run.ApprovedBy, &run.ApprovedAt, &run.DisbursedAt, &run.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return run, err
}

func (r *payrollRepo) ListPayrollRunsByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.PayrollRun, error) {
	query := `SELECT id, org_id, month, year, status, COALESCE(total_gross, 0), COALESCE(total_deductions, 0), COALESCE(total_net, 0), run_by, approved_by, approved_at, disbursed_at, created_at
		FROM payroll_runs WHERE org_id = $1 ORDER BY year DESC, month DESC`

	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []*models.PayrollRun
	for rows.Next() {
		run := &models.PayrollRun{}
		if err := rows.Scan(
			&run.ID, &run.OrgID, &run.Month, &run.Year, &run.Status, &run.TotalGross, &run.TotalDeductions, &run.TotalNet,
			&run.RunBy, &run.ApprovedBy, &run.ApprovedAt, &run.DisbursedAt, &run.CreatedAt,
		); err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, nil
}

func (r *payrollRepo) UpdateRunTotalsAndStatus(ctx context.Context, id uuid.UUID, gross, deductions, net float64, status string) error {
	query := `UPDATE payroll_runs SET total_gross=$1, total_deductions=$2, total_net=$3, status=$4 WHERE id=$5`
	_, err := r.db.Exec(ctx, query, gross, deductions, net, status, id)
	return err
}

// ── Payslips ─────────────────────────────────────────────────────────────────

func (r *payrollRepo) CreatePayslip(ctx context.Context, p *models.Payslip) error {
	query := `INSERT INTO payslips (emp_id, payroll_run_id, working_days, lop_days, paid_days, gross, total_deductions, net_pay, breakdown, payslip_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, generated_at`
	return r.db.QueryRow(ctx, query, p.EmpID, p.PayrollRunID, p.WorkingDays, p.LopDays, p.PaidDays, p.Gross, p.TotalDeductions, p.NetPay, p.Breakdown, p.PayslipURL).
		Scan(&p.ID, &p.GeneratedAt)
}

func (r *payrollRepo) GetPayslipByID(ctx context.Context, id uuid.UUID) (*models.Payslip, error) {
	query := `SELECT p.id, p.emp_id, p.payroll_run_id, p.working_days, p.lop_days, p.paid_days, p.gross, p.total_deductions, p.net_pay, p.breakdown, p.payslip_url, p.generated_at,
		e.first_name || ' ' || e.last_name as employee_name, e.emp_code
	FROM payslips p
	JOIN employees e ON e.id = p.emp_id
	WHERE p.id = $1`

	p := &models.Payslip{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.EmpID, &p.PayrollRunID, &p.WorkingDays, &p.LopDays, &p.PaidDays, &p.Gross, &p.TotalDeductions, &p.NetPay, &p.Breakdown, &p.PayslipURL, &p.GeneratedAt,
		&p.EmployeeName, &p.EmpCode,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *payrollRepo) ListPayslipsByRun(ctx context.Context, runID uuid.UUID) ([]*models.Payslip, error) {
	query := `SELECT p.id, p.emp_id, p.payroll_run_id, p.working_days, p.lop_days, p.paid_days, p.gross, p.total_deductions, p.net_pay, p.breakdown, p.payslip_url, p.generated_at,
		e.first_name || ' ' || e.last_name as employee_name, e.emp_code
	FROM payslips p
	JOIN employees e ON e.id = p.emp_id
	WHERE p.payroll_run_id = $1 ORDER BY e.emp_code ASC`

	rows, err := r.db.Query(ctx, query, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payslips []*models.Payslip
	for rows.Next() {
		p := &models.Payslip{}
		if err := rows.Scan(
			&p.ID, &p.EmpID, &p.PayrollRunID, &p.WorkingDays, &p.LopDays, &p.PaidDays, &p.Gross, &p.TotalDeductions, &p.NetPay, &p.Breakdown, &p.PayslipURL, &p.GeneratedAt,
			&p.EmployeeName, &p.EmpCode,
		); err != nil {
			return nil, err
		}
		payslips = append(payslips, p)
	}
	return payslips, nil
}

func (r *payrollRepo) ListPayslipsByEmp(ctx context.Context, empID uuid.UUID) ([]*models.Payslip, error) {
	query := `SELECT p.id, p.emp_id, p.payroll_run_id, p.working_days, p.lop_days, p.paid_days, p.gross, p.total_deductions, p.net_pay, p.breakdown, p.payslip_url, p.generated_at
	FROM payslips p WHERE p.emp_id = $1 ORDER BY p.generated_at DESC`

	rows, err := r.db.Query(ctx, query, empID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payslips []*models.Payslip
	for rows.Next() {
		p := &models.Payslip{}
		if err := rows.Scan(
			&p.ID, &p.EmpID, &p.PayrollRunID, &p.WorkingDays, &p.LopDays, &p.PaidDays, &p.Gross, &p.TotalDeductions, &p.NetPay, &p.Breakdown, &p.PayslipURL, &p.GeneratedAt,
		); err != nil {
			return nil, err
		}
		payslips = append(payslips, p)
	}
	return payslips, nil
}

// ── Reimbursements ───────────────────────────────────────────────────────────

func (r *payrollRepo) CreateReimbursement(ctx context.Context, rMB *models.Reimbursement) error {
	query := `INSERT INTO reimbursements (emp_id, category, amount, description, receipt_url)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, rMB.EmpID, rMB.Category, rMB.Amount, rMB.Description, rMB.ReceiptURL).
		Scan(&rMB.ID, &rMB.CreatedAt)
}

func (r *payrollRepo) ListReimbursementsByEmp(ctx context.Context, empID uuid.UUID) ([]*models.Reimbursement, error) {
	query := `SELECT id, emp_id, category, amount, description, receipt_url, status, approved_by, payroll_run_id, created_at
		FROM reimbursements WHERE emp_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, empID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Reimbursement
	for rows.Next() {
		rm := &models.Reimbursement{}
		if err := rows.Scan(&rm.ID, &rm.EmpID, &rm.Category, &rm.Amount, &rm.Description, &rm.ReceiptURL, &rm.Status, &rm.ApprovedBy, &rm.PayrollRunID, &rm.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, rm)
	}
	return list, nil
}

func (r *payrollRepo) ListReimbursementsByOrg(ctx context.Context, orgID uuid.UUID, status string) ([]*models.Reimbursement, error) {
	query := `SELECT r.id, r.emp_id, r.category, r.amount, r.description, r.receipt_url, r.status, r.approved_by, r.payroll_run_id, r.created_at,
		e.first_name || ' ' || e.last_name as employee_name
	FROM reimbursements r
	JOIN employees e ON e.id = r.emp_id
	WHERE e.org_id = $1`

	args := []interface{}{orgID}
	if status != "" {
		query += ` AND r.status = $2`
		args = append(args, status)
	}
	query += ` ORDER BY r.created_at DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Reimbursement
	for rows.Next() {
		rm := &models.Reimbursement{}
		if err := rows.Scan(&rm.ID, &rm.EmpID, &rm.Category, &rm.Amount, &rm.Description, &rm.ReceiptURL, &rm.Status, &rm.ApprovedBy, &rm.PayrollRunID, &rm.CreatedAt, &rm.EmployeeName); err != nil {
			return nil, err
		}
		list = append(list, rm)
	}
	return list, nil
}

func (r *payrollRepo) UpdateReimbursementStatus(ctx context.Context, id uuid.UUID, status string, approverID *uuid.UUID) error {
	query := `UPDATE reimbursements SET status=$1, approved_by=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, status, approverID, id)
	return err
}

// ── Tax Declarations ─────────────────────────────────────────────────────────

func (r *payrollRepo) SaveTaxDeclaration(ctx context.Context, t *models.TaxDeclaration) error {
	query := `INSERT INTO tax_declarations (emp_id, year, regime, sections)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (emp_id, year) DO UPDATE SET regime=$3, sections=$4, submitted_at=NOW()
		RETURNING id, submitted_at`
	return r.db.QueryRow(ctx, query, t.EmpID, t.Year, t.Regime, t.Sections).Scan(&t.ID, &t.SubmittedAt)
}

func (r *payrollRepo) GetTaxDeclaration(ctx context.Context, empID uuid.UUID, year int) (*models.TaxDeclaration, error) {
	query := `SELECT id, emp_id, year, regime, sections, submitted_at FROM tax_declarations WHERE emp_id = $1 AND year = $2`
	t := &models.TaxDeclaration{}
	err := r.db.QueryRow(ctx, query, empID, year).Scan(&t.ID, &t.EmpID, &t.Year, &t.Regime, &t.Sections, &t.SubmittedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return t, err
}