package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/hrms-backend/internal/models"
)

type EmployeeRepository interface {
	// Employee Core
	Create(ctx context.Context, emp *models.Employee) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Employee, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Employee, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID, search string, deptID, desigID *uuid.UUID, status string, limit, offset int) ([]*models.Employee, int, error)
	Update(ctx context.Context, emp *models.Employee) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetNextEmpSeq(ctx context.Context, orgID uuid.UUID) (int, error)

	// Departments
	CreateDepartment(ctx context.Context, d *models.Department) error
	GetDepartmentByID(ctx context.Context, id uuid.UUID) (*models.Department, error)
	ListDepartmentsByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.Department, error)
	UpdateDepartment(ctx context.Context, d *models.Department) error
	DeleteDepartment(ctx context.Context, id uuid.UUID) error

	// Designations
	CreateDesignation(ctx context.Context, des *models.Designation) error
	GetDesignationByID(ctx context.Context, id uuid.UUID) (*models.Designation, error)
	ListDesignationsByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.Designation, error)
	UpdateDesignation(ctx context.Context, des *models.Designation) error
	DeleteDesignation(ctx context.Context, id uuid.UUID) error

	// Bank Details
	SaveBankDetails(ctx context.Context, b *models.EmployeeBankDetails) error
	GetBankDetailsByEmpID(ctx context.Context, empID uuid.UUID) (*models.EmployeeBankDetails, error)

	// Emergency Contacts
	AddEmergencyContact(ctx context.Context, c *models.EmergencyContact) error
	ListEmergencyContacts(ctx context.Context, empID uuid.UUID) ([]*models.EmergencyContact, error)
	DeleteEmergencyContact(ctx context.Context, id uuid.UUID) error

	// Documents
	AddDocument(ctx context.Context, doc *models.EmployeeDocument) error
	ListDocuments(ctx context.Context, empID uuid.UUID) ([]*models.EmployeeDocument, error)
	VerifyDocument(ctx context.Context, docID, verifierID uuid.UUID) error
	DeleteDocument(ctx context.Context, docID uuid.UUID) error

	// Onboarding Tasks
	CreateOnboardingTask(ctx context.Context, task *models.OnboardingTask) error
	ListOnboardingTasks(ctx context.Context, empID uuid.UUID) ([]*models.OnboardingTask, error)
	UpdateOnboardingTask(ctx context.Context, taskID uuid.UUID, status string) error
}

type employeeRepo struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) EmployeeRepository {
	return &employeeRepo{db: db}
}

// ── Employee Core ────────────────────────────────────────────────────────────

func (r *employeeRepo) Create(ctx context.Context, e *models.Employee) error {
	query := `INSERT INTO employees (
		org_id, user_id, emp_code, first_name, last_name, middle_name,
		dob, gender, blood_group, phone, personal_email, address, city, state, country, pincode, nationality,
		dept_id, designation_id, manager_id, employment_type, work_location, joining_date, status
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24)
	RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		e.OrgID, e.UserID, e.EmpCode, e.FirstName, e.LastName, e.MiddleName,
		e.DOB, e.Gender, e.BloodGroup, e.Phone, e.PersonalEmail, e.Address, e.City, e.State, e.Country, e.Pincode, e.Nationality,
		e.DeptID, e.DesignationID, e.ManagerID, e.EmploymentType, e.WorkLocation, e.JoiningDate, e.Status,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (r *employeeRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	query := `SELECT e.id, e.org_id, e.user_id, e.emp_code, e.first_name, e.last_name, e.middle_name,
		e.dob, e.gender, e.blood_group, e.phone, e.personal_email, e.address, e.city, e.state, e.country, e.pincode, e.nationality,
		e.dept_id, e.designation_id, e.manager_id, e.employment_type, e.work_location, e.joining_date, e.confirmation_date,
		e.status, e.exit_date, e.exit_reason, e.photo_url, e.created_at, e.updated_at,
		d.name as department_name, des.title as designation_name,
		COALESCE(m.first_name || ' ' || m.last_name, '') as manager_name
	FROM employees e
	LEFT JOIN departments d ON d.id = e.dept_id
	LEFT JOIN designations des ON des.id = e.designation_id
	LEFT JOIN employees m ON m.id = e.manager_id
	WHERE e.id = $1`

	e := &models.Employee{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.OrgID, &e.UserID, &e.EmpCode, &e.FirstName, &e.LastName, &e.MiddleName,
		&e.DOB, &e.Gender, &e.BloodGroup, &e.Phone, &e.PersonalEmail, &e.Address, &e.City, &e.State, &e.Country, &e.Pincode, &e.Nationality,
		&e.DeptID, &e.DesignationID, &e.ManagerID, &e.EmploymentType, &e.WorkLocation, &e.JoiningDate, &e.ConfirmationDate,
		&e.Status, &e.ExitDate, &e.ExitReason, &e.PhotoURL, &e.CreatedAt, &e.UpdatedAt,
		&e.DepartmentName, &e.DesignationName, &e.ManagerName,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return e, err
}

func (r *employeeRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Employee, error) {
	query := `SELECT e.id, e.org_id, e.user_id, e.emp_code, e.first_name, e.last_name,
		e.joining_date, e.status, e.created_at, e.updated_at,
		d.name as department_name, des.title as designation_name
	FROM employees e
	LEFT JOIN departments d ON d.id = e.dept_id
	LEFT JOIN designations des ON des.id = e.designation_id
	WHERE e.user_id = $1`

	e := &models.Employee{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&e.ID, &e.OrgID, &e.UserID, &e.EmpCode, &e.FirstName, &e.LastName,
		&e.JoiningDate, &e.Status, &e.CreatedAt, &e.UpdatedAt,
		&e.DepartmentName, &e.DesignationName,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return e, err
}

func (r *employeeRepo) ListByOrg(ctx context.Context, orgID uuid.UUID, search string, deptID, desigID *uuid.UUID, status string, limit, offset int) ([]*models.Employee, int, error) {
	countQuery := `SELECT COUNT(*) FROM employees e WHERE e.org_id = $1`
	var total int
	_ = r.db.QueryRow(ctx, countQuery, orgID).Scan(&total)

	query := `SELECT e.id, e.org_id, e.user_id, e.emp_code, e.first_name, e.last_name, e.phone, e.personal_email,
		e.joining_date, e.status, e.photo_url, e.created_at, e.updated_at,
		d.name as department_name, des.title as designation_name
	FROM employees e
	LEFT JOIN departments d ON d.id = e.dept_id
	LEFT JOIN designations des ON des.id = e.designation_id
	WHERE e.org_id = $1
	ORDER BY e.created_at DESC
	LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var emps []*models.Employee
	for rows.Next() {
		e := &models.Employee{}
		if err := rows.Scan(
			&e.ID, &e.OrgID, &e.UserID, &e.EmpCode, &e.FirstName, &e.LastName, &e.Phone, &e.PersonalEmail,
			&e.JoiningDate, &e.Status, &e.PhotoURL, &e.CreatedAt, &e.UpdatedAt,
			&e.DepartmentName, &e.DesignationName,
		); err != nil {
			return nil, 0, err
		}
		emps = append(emps, e)
	}
	return emps, total, nil
}

func (r *employeeRepo) Update(ctx context.Context, e *models.Employee) error {
	query := `UPDATE employees SET
		first_name=$1, last_name=$2, middle_name=$3, phone=$4, personal_email=$5, address=$6,
		city=$7, state=$8, country=$9, pincode=$10,
		dept_id=$11, designation_id=$12, manager_id=$13, status=$14, work_location=$15,
		updated_at=NOW()
	WHERE id=$16 AND org_id=$17`

	_, err := r.db.Exec(ctx, query,
		e.FirstName, e.LastName, e.MiddleName, e.Phone, e.PersonalEmail, e.Address,
		e.City, e.State, e.Country, e.Pincode,
		e.DeptID, e.DesignationID, e.ManagerID, e.Status, e.WorkLocation,
		e.ID, e.OrgID,
	)
	return err
}

func (r *employeeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM employees WHERE id=$1`, id)
	return err
}

func (r *employeeRepo) GetNextEmpSeq(ctx context.Context, orgID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM employees WHERE org_id=$1`, orgID).Scan(&count)
	return count + 1, err
}

// ── Departments ──────────────────────────────────────────────────────────────

func (r *employeeRepo) CreateDepartment(ctx context.Context, d *models.Department) error {
	query := `INSERT INTO departments (org_id, name, code, parent_dept_id, head_emp_id, description)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, d.OrgID, d.Name, d.Code, d.ParentDeptID, d.HeadEmpID, d.Description).
		Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

func (r *employeeRepo) GetDepartmentByID(ctx context.Context, id uuid.UUID) (*models.Department, error) {
	query := `SELECT d.id, d.org_id, d.name, d.code, d.parent_dept_id, d.head_emp_id, d.description, d.is_active, d.created_at, d.updated_at,
		COALESCE(e.first_name || ' ' || e.last_name, '') as head_emp_name,
		COALESCE(p.name, '') as parent_name
	FROM departments d
	LEFT JOIN employees e ON e.id = d.head_emp_id
	LEFT JOIN departments p ON p.id = d.parent_dept_id
	WHERE d.id = $1`

	d := &models.Department{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.OrgID, &d.Name, &d.Code, &d.ParentDeptID, &d.HeadEmpID, &d.Description, &d.IsActive, &d.CreatedAt, &d.UpdatedAt,
		&d.HeadEmpName, &d.ParentName,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return d, err
}

func (r *employeeRepo) ListDepartmentsByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.Department, error) {
	query := `SELECT d.id, d.org_id, d.name, d.code, d.parent_dept_id, d.head_emp_id, d.description, d.is_active, d.created_at, d.updated_at,
		COALESCE(e.first_name || ' ' || e.last_name, '') as head_emp_name,
		COALESCE(p.name, '') as parent_name
	FROM departments d
	LEFT JOIN employees e ON e.id = d.head_emp_id
	LEFT JOIN departments p ON p.id = d.parent_dept_id
	WHERE d.org_id = $1 ORDER BY d.name`

	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var depts []*models.Department
	for rows.Next() {
		d := &models.Department{}
		if err := rows.Scan(
			&d.ID, &d.OrgID, &d.Name, &d.Code, &d.ParentDeptID, &d.HeadEmpID, &d.Description, &d.IsActive, &d.CreatedAt, &d.UpdatedAt,
			&d.HeadEmpName, &d.ParentName,
		); err != nil {
			return nil, err
		}
		depts = append(depts, d)
	}
	return depts, nil
}

func (r *employeeRepo) UpdateDepartment(ctx context.Context, d *models.Department) error {
	query := `UPDATE departments SET name=$1, code=$2, parent_dept_id=$3, head_emp_id=$4, description=$5, is_active=$6, updated_at=NOW()
		WHERE id=$7 AND org_id=$8`
	_, err := r.db.Exec(ctx, query, d.Name, d.Code, d.ParentDeptID, d.HeadEmpID, d.Description, d.IsActive, d.ID, d.OrgID)
	return err
}

func (r *employeeRepo) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM departments WHERE id=$1`, id)
	return err
}

// ── Designations ─────────────────────────────────────────────────────────────

func (r *employeeRepo) CreateDesignation(ctx context.Context, des *models.Designation) error {
	query := `INSERT INTO designations (org_id, department_id, title, grade, level)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, des.OrgID, des.DepartmentID, des.Title, des.Grade, des.Level).
		Scan(&des.ID, &des.CreatedAt)
}

func (r *employeeRepo) GetDesignationByID(ctx context.Context, id uuid.UUID) (*models.Designation, error) {
	query := `SELECT des.id, des.org_id, des.department_id, des.title, des.grade, des.level, des.created_at,
		COALESCE(d.name, '') as department_name
	FROM designations des
	LEFT JOIN departments d ON d.id = des.department_id
	WHERE des.id = $1`

	des := &models.Designation{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&des.ID, &des.OrgID, &des.DepartmentID, &des.Title, &des.Grade, &des.Level, &des.CreatedAt, &des.DepartmentName,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return des, err
}

func (r *employeeRepo) ListDesignationsByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.Designation, error) {
	query := `SELECT des.id, des.org_id, des.department_id, des.title, des.grade, des.level, des.created_at,
		COALESCE(d.name, '') as department_name
	FROM designations des
	LEFT JOIN departments d ON d.id = des.department_id
	WHERE des.org_id = $1 ORDER BY des.title`

	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var desigs []*models.Designation
	for rows.Next() {
		des := &models.Designation{}
		if err := rows.Scan(
			&des.ID, &des.OrgID, &des.DepartmentID, &des.Title, &des.Grade, &des.Level, &des.CreatedAt, &des.DepartmentName,
		); err != nil {
			return nil, err
		}
		desigs = append(desigs, des)
	}
	return desigs, nil
}

func (r *employeeRepo) UpdateDesignation(ctx context.Context, des *models.Designation) error {
	query := `UPDATE designations SET title=$1, department_id=$2, grade=$3, level=$4 WHERE id=$5 AND org_id=$6`
	_, err := r.db.Exec(ctx, query, des.Title, des.DepartmentID, des.Grade, des.Level, des.ID, des.OrgID)
	return err
}

func (r *employeeRepo) DeleteDesignation(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM designations WHERE id=$1`, id)
	return err
}

// ── Bank Details ─────────────────────────────────────────────────────────────

func (r *employeeRepo) SaveBankDetails(ctx context.Context, b *models.EmployeeBankDetails) error {
	query := `INSERT INTO employee_bank_details (emp_id, bank_name, account_no_enc, ifsc_code, account_type, is_primary)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET bank_name=$2, account_no_enc=$3, ifsc_code=$4, account_type=$5, updated_at=NOW()
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, b.EmpID, b.BankName, b.AccountNoEnc, b.IFSCCode, b.AccountType, b.IsPrimary).
		Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *employeeRepo) GetBankDetailsByEmpID(ctx context.Context, empID uuid.UUID) (*models.EmployeeBankDetails, error) {
	query := `SELECT id, emp_id, bank_name, account_no_enc, ifsc_code, account_type, is_primary, created_at, updated_at
		FROM employee_bank_details WHERE emp_id = $1 AND is_primary = TRUE LIMIT 1`
	b := &models.EmployeeBankDetails{}
	err := r.db.QueryRow(ctx, query, empID).Scan(
		&b.ID, &b.EmpID, &b.BankName, &b.AccountNoEnc, &b.IFSCCode, &b.AccountType, &b.IsPrimary, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return b, err
}

// ── Emergency Contacts ───────────────────────────────────────────────────────

func (r *employeeRepo) AddEmergencyContact(ctx context.Context, c *models.EmergencyContact) error {
	query := `INSERT INTO emergency_contacts (emp_id, name, relation, phone, email, address)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRow(ctx, query, c.EmpID, c.Name, c.Relation, c.Phone, c.Email, c.Address).Scan(&c.ID)
}

func (r *employeeRepo) ListEmergencyContacts(ctx context.Context, empID uuid.UUID) ([]*models.EmergencyContact, error) {
	query := `SELECT id, emp_id, name, relation, phone, email, address FROM emergency_contacts WHERE emp_id = $1`
	rows, err := r.db.Query(ctx, query, empID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*models.EmergencyContact
	for rows.Next() {
		c := &models.EmergencyContact{}
		if err := rows.Scan(&c.ID, &c.EmpID, &c.Name, &c.Relation, &c.Phone, &c.Email, &c.Address); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *employeeRepo) DeleteEmergencyContact(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM emergency_contacts WHERE id=$1`, id)
	return err
}

// ── Documents ────────────────────────────────────────────────────────────────

func (r *employeeRepo) AddDocument(ctx context.Context, doc *models.EmployeeDocument) error {
	query := `INSERT INTO employee_documents (emp_id, doc_type, file_url, file_name)
		VALUES ($1, $2, $3, $4) RETURNING id, uploaded_at`
	return r.db.QueryRow(ctx, query, doc.EmpID, doc.DocType, doc.FileURL, doc.FileName).Scan(&doc.ID, &doc.UploadedAt)
}

func (r *employeeRepo) ListDocuments(ctx context.Context, empID uuid.UUID) ([]*models.EmployeeDocument, error) {
	query := `SELECT id, emp_id, doc_type, file_url, file_name, is_verified, verified_by, verified_at, uploaded_at
		FROM employee_documents WHERE emp_id = $1 ORDER BY uploaded_at DESC`
	rows, err := r.db.Query(ctx, query, empID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []*models.EmployeeDocument
	for rows.Next() {
		d := &models.EmployeeDocument{}
		if err := rows.Scan(&d.ID, &d.EmpID, &d.DocType, &d.FileURL, &d.FileName, &d.IsVerified, &d.VerifiedBy, &d.VerifiedAt, &d.UploadedAt); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

func (r *employeeRepo) VerifyDocument(ctx context.Context, docID, verifierID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE employee_documents SET is_verified=TRUE, verified_by=$1, verified_at=NOW() WHERE id=$2`, verifierID, docID)
	return err
}

func (r *employeeRepo) DeleteDocument(ctx context.Context, docID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM employee_documents WHERE id=$1`, docID)
	return err
}

// ── Onboarding Tasks ─────────────────────────────────────────────────────────

func (r *employeeRepo) CreateOnboardingTask(ctx context.Context, task *models.OnboardingTask) error {
	query := `INSERT INTO onboarding_tasks (emp_id, task_name, description, assigned_to, due_date)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(ctx, query, task.EmpID, task.TaskName, task.Description, task.AssignedTo, task.DueDate).Scan(&task.ID)
}

func (r *employeeRepo) ListOnboardingTasks(ctx context.Context, empID uuid.UUID) ([]*models.OnboardingTask, error) {
	query := `SELECT id, emp_id, task_name, description, assigned_to, due_date, completed_at, status
		FROM onboarding_tasks WHERE emp_id = $1 ORDER BY due_date ASC`
	rows, err := r.db.Query(ctx, query, empID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.OnboardingTask
	for rows.Next() {
		t := &models.OnboardingTask{}
		if err := rows.Scan(&t.ID, &t.EmpID, &t.TaskName, &t.Description, &t.AssignedTo, &t.DueDate, &t.CompletedAt, &t.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *employeeRepo) UpdateOnboardingTask(ctx context.Context, taskID uuid.UUID, status string) error {
	var completedAt *time.Time
	if status == "completed" {
		now := time.Now()
		completedAt = &now
	}
	_, err := r.db.Exec(ctx, `UPDATE onboarding_tasks SET status=$1, completed_at=$2 WHERE id=$3`, status, completedAt, taskID)
	return err
}