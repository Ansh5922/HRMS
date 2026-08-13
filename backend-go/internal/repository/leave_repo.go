package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/hrms-backend/internal/models"
)

type LeaveRepository interface {
	// Leave Types
	CreateLeaveType(ctx context.Context, lt *models.LeaveType) error
	GetLeaveTypeByID(ctx context.Context, id uuid.UUID) (*models.LeaveType, error)
	ListLeaveTypesByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.LeaveType, error)
	UpdateLeaveType(ctx context.Context, lt *models.LeaveType) error
	DeleteLeaveType(ctx context.Context, id uuid.UUID) error

	// Leave Policies
	CreateLeavePolicy(ctx context.Context, lp *models.LeavePolicy) error
	ListLeavePoliciesByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.LeavePolicy, error)
	DeleteLeavePolicy(ctx context.Context, id uuid.UUID) error

	// Leave Balances
	SetLeaveBalance(ctx context.Context, lb *models.LeaveBalance) error
	GetLeaveBalanceByEmp(ctx context.Context, empID uuid.UUID, year int) ([]*models.LeaveBalance, error)
	GetLeaveBalanceForType(ctx context.Context, empID, typeID uuid.UUID, year int) (*models.LeaveBalance, error)
	AdjustLeaveBalanceOnApply(ctx context.Context, empID, typeID uuid.UUID, year int, days float64) error
	AdjustLeaveBalanceOnApprove(ctx context.Context, empID, typeID uuid.UUID, year int, days float64) error
	AdjustLeaveBalanceOnReject(ctx context.Context, empID, typeID uuid.UUID, year int, days float64) error

	// Leave Applications
	Apply(ctx context.Context, app *models.LeaveApplication) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.LeaveApplication, error)
	ListByEmp(ctx context.Context, empID uuid.UUID) ([]*models.LeaveApplication, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID, status string) ([]*models.LeaveApplication, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error

	// Leave Approvals
	AddApprovalRecord(ctx context.Context, approval *models.LeaveApproval) error
	ListApprovals(ctx context.Context, appID uuid.UUID) ([]*models.LeaveApproval, error)
}

type leaveRepo struct {
	db *pgxpool.Pool
}

func NewLeaveRepository(db *pgxpool.Pool) LeaveRepository {
	return &leaveRepo{db: db}
}

// ── Leave Types ──────────────────────────────────────────────────────────────

func (r *leaveRepo) CreateLeaveType(ctx context.Context, lt *models.LeaveType) error {
	query := `INSERT INTO leave_types (
		org_id, name, code, max_days_per_year, carry_forward, carry_forward_max,
		is_paid, requires_doc, min_notice_days, gender_applicable, requires_approval
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		lt.OrgID, lt.Name, lt.Code, lt.MaxDaysPerYear, lt.CarryForward, lt.CarryForwardMax,
		lt.IsPaid, lt.RequiresDoc, lt.MinNoticeDays, lt.GenderApplicable, lt.RequiresApproval,
	).Scan(&lt.ID, &lt.CreatedAt)
}

func (r *leaveRepo) GetLeaveTypeByID(ctx context.Context, id uuid.UUID) (*models.LeaveType, error) {
	query := `SELECT id, org_id, name, code, max_days_per_year, carry_forward, carry_forward_max,
		is_paid, requires_doc, min_notice_days, gender_applicable, requires_approval, created_at
	FROM leave_types WHERE id = $1`

	lt := &models.LeaveType{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&lt.ID, &lt.OrgID, &lt.Name, &lt.Code, &lt.MaxDaysPerYear, &lt.CarryForward, &lt.CarryForwardMax,
		&lt.IsPaid, &lt.RequiresDoc, &lt.MinNoticeDays, &lt.GenderApplicable, &lt.RequiresApproval, &lt.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return lt, err
}

func (r *leaveRepo) ListLeaveTypesByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.LeaveType, error) {
	query := `SELECT id, org_id, name, code, max_days_per_year, carry_forward, carry_forward_max,
		is_paid, requires_doc, min_notice_days, gender_applicable, requires_approval, created_at
	FROM leave_types WHERE org_id = $1 ORDER BY name`

	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []*models.LeaveType
	for rows.Next() {
		lt := &models.LeaveType{}
		if err := rows.Scan(
			&lt.ID, &lt.OrgID, &lt.Name, &lt.Code, &lt.MaxDaysPerYear, &lt.CarryForward, &lt.CarryForwardMax,
			&lt.IsPaid, &lt.RequiresDoc, &lt.MinNoticeDays, &lt.GenderApplicable, &lt.RequiresApproval, &lt.CreatedAt,
		); err != nil {
			return nil, err
		}
		types = append(types, lt)
	}
	return types, nil
}

func (r *leaveRepo) UpdateLeaveType(ctx context.Context, lt *models.LeaveType) error {
	query := `UPDATE leave_types SET
		name=$1, code=$2, max_days_per_year=$3, carry_forward=$4, carry_forward_max=$5,
		is_paid=$6, requires_doc=$7, min_notice_days=$8, gender_applicable=$9, requires_approval=$10
	WHERE id=$11 AND org_id=$12`

	_, err := r.db.Exec(ctx, query,
		lt.Name, lt.Code, lt.MaxDaysPerYear, lt.CarryForward, lt.CarryForwardMax,
		lt.IsPaid, lt.RequiresDoc, lt.MinNoticeDays, lt.GenderApplicable, lt.RequiresApproval,
		lt.ID, lt.OrgID,
	)
	return err
}

func (r *leaveRepo) DeleteLeaveType(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM leave_types WHERE id=$1`, id)
	return err
}

// ── Leave Policies ───────────────────────────────────────────────────────────

func (r *leaveRepo) CreateLeavePolicy(ctx context.Context, lp *models.LeavePolicy) error {
	query := `INSERT INTO leave_policies (org_id, leave_type_id, designation_id, department_id, annual_quota, accrual_type)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, lp.OrgID, lp.LeaveTypeID, lp.DesignationID, lp.DepartmentID, lp.AnnualQuota, lp.AccrualType).
		Scan(&lp.ID, &lp.CreatedAt)
}

func (r *leaveRepo) ListLeavePoliciesByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.LeavePolicy, error) {
	query := `SELECT lp.id, lp.org_id, lp.leave_type_id, lp.designation_id, lp.department_id, lp.annual_quota, lp.accrual_type, lp.created_at,
		lt.name as leave_type_name
	FROM leave_policies lp
	JOIN leave_types lt ON lt.id = lp.leave_type_id
	WHERE lp.org_id = $1`

	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []*models.LeavePolicy
	for rows.Next() {
		lp := &models.LeavePolicy{}
		if err := rows.Scan(&lp.ID, &lp.OrgID, &lp.LeaveTypeID, &lp.DesignationID, &lp.DepartmentID, &lp.AnnualQuota, &lp.AccrualType, &lp.CreatedAt, &lp.LeaveTypeName); err != nil {
			return nil, err
		}
		policies = append(policies, lp)
	}
	return policies, nil
}

func (r *leaveRepo) DeleteLeavePolicy(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM leave_policies WHERE id=$1`, id)
	return err
}

// ── Leave Balances ───────────────────────────────────────────────────────────

func (r *leaveRepo) SetLeaveBalance(ctx context.Context, lb *models.LeaveBalance) error {
	query := `INSERT INTO leave_balances (emp_id, leave_type_id, year, total, used, pending)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (emp_id, leave_type_id, year)
		DO UPDATE SET total=$4, updated_at=NOW()
		RETURNING id, available, updated_at`

	return r.db.QueryRow(ctx, query, lb.EmpID, lb.LeaveTypeID, lb.Year, lb.Total, lb.Used, lb.Pending).
		Scan(&lb.ID, &lb.Available, &lb.UpdatedAt)
}

func (r *leaveRepo) GetLeaveBalanceByEmp(ctx context.Context, empID uuid.UUID, year int) ([]*models.LeaveBalance, error) {
	query := `SELECT lb.id, lb.emp_id, lb.leave_type_id, lb.year, lb.total, lb.used, lb.pending, (lb.total - lb.used - lb.pending) as available, lb.updated_at,
		lt.name as leave_type_name, lt.code as leave_type_code
	FROM leave_balances lb
	JOIN leave_types lt ON lt.id = lb.leave_type_id
	WHERE lb.emp_id = $1 AND lb.year = $2`

	rows, err := r.db.Query(ctx, query, empID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []*models.LeaveBalance
	for rows.Next() {
		lb := &models.LeaveBalance{}
		if err := rows.Scan(&lb.ID, &lb.EmpID, &lb.LeaveTypeID, &lb.Year, &lb.Total, &lb.Used, &lb.Pending, &lb.Available, &lb.UpdatedAt, &lb.LeaveTypeName, &lb.LeaveTypeCode); err != nil {
			return nil, err
		}
		balances = append(balances, lb)
	}
	return balances, nil
}

func (r *leaveRepo) GetLeaveBalanceForType(ctx context.Context, empID, typeID uuid.UUID, year int) (*models.LeaveBalance, error) {
	query := `SELECT lb.id, lb.emp_id, lb.leave_type_id, lb.year, lb.total, lb.used, lb.pending, (lb.total - lb.used - lb.pending) as available, lb.updated_at,
		lt.name as leave_type_name, lt.code as leave_type_code
	FROM leave_balances lb
	JOIN leave_types lt ON lt.id = lb.leave_type_id
	WHERE lb.emp_id = $1 AND lb.leave_type_id = $2 AND lb.year = $3`

	lb := &models.LeaveBalance{}
	err := r.db.QueryRow(ctx, query, empID, typeID, year).Scan(
		&lb.ID, &lb.EmpID, &lb.LeaveTypeID, &lb.Year, &lb.Total, &lb.Used, &lb.Pending, &lb.Available, &lb.UpdatedAt, &lb.LeaveTypeName, &lb.LeaveTypeCode,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return lb, err
}

func (r *leaveRepo) AdjustLeaveBalanceOnApply(ctx context.Context, empID, typeID uuid.UUID, year int, days float64) error {
	query := `UPDATE leave_balances SET pending = pending + $1, updated_at = NOW()
		WHERE emp_id = $2 AND leave_type_id = $3 AND year = $4`
	_, err := r.db.Exec(ctx, query, days, empID, typeID, year)
	return err
}

func (r *leaveRepo) AdjustLeaveBalanceOnApprove(ctx context.Context, empID, typeID uuid.UUID, year int, days float64) error {
	query := `UPDATE leave_balances SET pending = GREATEST(0, pending - $1), used = used + $1, updated_at = NOW()
		WHERE emp_id = $2 AND leave_type_id = $3 AND year = $4`
	_, err := r.db.Exec(ctx, query, days, empID, typeID, year)
	return err
}

func (r *leaveRepo) AdjustLeaveBalanceOnReject(ctx context.Context, empID, typeID uuid.UUID, year int, days float64) error {
	query := `UPDATE leave_balances SET pending = GREATEST(0, pending - $1), updated_at = NOW()
		WHERE emp_id = $2 AND leave_type_id = $3 AND year = $4`
	_, err := r.db.Exec(ctx, query, days, empID, typeID, year)
	return err
}

// ── Leave Applications ───────────────────────────────────────────────────────

func (r *leaveRepo) Apply(ctx context.Context, a *models.LeaveApplication) error {
	query := `INSERT INTO leave_applications (emp_id, leave_type_id, from_date, to_date, days, session, reason, doc_url, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		a.EmpID, a.LeaveTypeID, a.FromDate, a.ToDate, a.Days, a.Session, a.Reason, a.DocURL, a.Status,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *leaveRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.LeaveApplication, error) {
	query := `SELECT la.id, la.emp_id, la.leave_type_id, la.from_date, la.to_date, la.days, la.session, la.reason, la.doc_url, la.status, la.created_at, la.updated_at,
		e.first_name || ' ' || e.last_name as employee_name, e.emp_code,
		lt.name as leave_type_name, lt.code as leave_type_code
	FROM leave_applications la
	JOIN employees e ON e.id = la.emp_id
	JOIN leave_types lt ON lt.id = la.leave_type_id
	WHERE la.id = $1`

	a := &models.LeaveApplication{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.EmpID, &a.LeaveTypeID, &a.FromDate, &a.ToDate, &a.Days, &a.Session, &a.Reason, &a.DocURL, &a.Status, &a.CreatedAt, &a.UpdatedAt,
		&a.EmployeeName, &a.EmpCode, &a.LeaveTypeName, &a.LeaveTypeCode,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return a, err
}

func (r *leaveRepo) ListByEmp(ctx context.Context, empID uuid.UUID) ([]*models.LeaveApplication, error) {
	query := `SELECT la.id, la.emp_id, la.leave_type_id, la.from_date, la.to_date, la.days, la.session, la.reason, la.doc_url, la.status, la.created_at, la.updated_at,
		lt.name as leave_type_name, lt.code as leave_type_code
	FROM leave_applications la
	JOIN leave_types lt ON lt.id = la.leave_type_id
	WHERE la.emp_id = $1 ORDER BY la.created_at DESC`

	rows, err := r.db.Query(ctx, query, empID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*models.LeaveApplication
	for rows.Next() {
		a := &models.LeaveApplication{}
		if err := rows.Scan(
			&a.ID, &a.EmpID, &a.LeaveTypeID, &a.FromDate, &a.ToDate, &a.Days, &a.Session, &a.Reason, &a.DocURL, &a.Status, &a.CreatedAt, &a.UpdatedAt,
			&a.LeaveTypeName, &a.LeaveTypeCode,
		); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, nil
}

func (r *leaveRepo) ListByOrg(ctx context.Context, orgID uuid.UUID, status string) ([]*models.LeaveApplication, error) {
	query := `SELECT la.id, la.emp_id, la.leave_type_id, la.from_date, la.to_date, la.days, la.session, la.reason, la.doc_url, la.status, la.created_at, la.updated_at,
		e.first_name || ' ' || e.last_name as employee_name, e.emp_code,
		lt.name as leave_type_name, lt.code as leave_type_code
	FROM leave_applications la
	JOIN employees e ON e.id = la.emp_id
	JOIN leave_types lt ON lt.id = la.leave_type_id
	WHERE e.org_id = $1`

	args := []interface{}{orgID}
	if status != "" {
		query += ` AND la.status = $2`
		args = append(args, status)
	}
	query += ` ORDER BY la.created_at DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*models.LeaveApplication
	for rows.Next() {
		a := &models.LeaveApplication{}
		if err := rows.Scan(
			&a.ID, &a.EmpID, &a.LeaveTypeID, &a.FromDate, &a.ToDate, &a.Days, &a.Session, &a.Reason, &a.DocURL, &a.Status, &a.CreatedAt, &a.UpdatedAt,
			&a.EmployeeName, &a.EmpCode, &a.LeaveTypeName, &a.LeaveTypeCode,
		); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, nil
}

func (r *leaveRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE leave_applications SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

// ── Leave Approvals ──────────────────────────────────────────────────────────

func (r *leaveRepo) AddApprovalRecord(ctx context.Context, approval *models.LeaveApproval) error {
	query := `INSERT INTO leave_approvals (application_id, approver_id, level, action, comment)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, actioned_at`
	return r.db.QueryRow(ctx, query, approval.ApplicationID, approval.ApproverID, approval.Level, approval.Action, approval.Comment).
		Scan(&approval.ID, &approval.ActionedAt)
}

func (r *leaveRepo) ListApprovals(ctx context.Context, appID uuid.UUID) ([]*models.LeaveApproval, error) {
	query := `SELECT la.id, la.application_id, la.approver_id, la.level, la.action, la.comment, la.actioned_at,
		u.email as approver_name
	FROM leave_approvals la
	JOIN users u ON u.id = la.approver_id
	WHERE la.application_id = $1 ORDER BY la.level ASC`

	rows, err := r.db.Query(ctx, query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var approvals []*models.LeaveApproval
	for rows.Next() {
		app := &models.LeaveApproval{}
		if err := rows.Scan(&app.ID, &app.ApplicationID, &app.ApproverID, &app.Level, &app.Action, &app.Comment, &app.ActionedAt, &app.ApproverName); err != nil {
			return nil, err
		}
		approvals = append(approvals, app)
	}
	return approvals, nil
}
