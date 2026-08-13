package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/repository"
	"github.com/your-org/hrms-backend/pkg/utils"
)

type EmployeeService interface {
	// Employee Core
	CreateEmployee(ctx context.Context, emp *models.Employee) (*models.Employee, error)
	GetEmployeeByID(ctx context.Context, id uuid.UUID) (*models.Employee, error)
	GetEmployeeByUserID(ctx context.Context, userID uuid.UUID) (*models.Employee, error)
	ListEmployees(ctx context.Context, orgID uuid.UUID, search string, deptID, desigID *uuid.UUID, status string, page, limit int) ([]*models.Employee, int, error)
	UpdateEmployee(ctx context.Context, emp *models.Employee) error
	DeleteEmployee(ctx context.Context, id uuid.UUID) error

	// Department
	CreateDepartment(ctx context.Context, d *models.Department) (*models.Department, error)
	GetDepartment(ctx context.Context, id uuid.UUID) (*models.Department, error)
	ListDepartments(ctx context.Context, orgID uuid.UUID) ([]*models.Department, error)
	UpdateDepartment(ctx context.Context, d *models.Department) error
	DeleteDepartment(ctx context.Context, id uuid.UUID) error

	// Designation
	CreateDesignation(ctx context.Context, des *models.Designation) (*models.Designation, error)
	GetDesignation(ctx context.Context, id uuid.UUID) (*models.Designation, error)
	ListDesignations(ctx context.Context, orgID uuid.UUID) ([]*models.Designation, error)
	UpdateDesignation(ctx context.Context, des *models.Designation) error
	DeleteDesignation(ctx context.Context, id uuid.UUID) error

	// Bank Details
	SaveBankDetails(ctx context.Context, b *models.EmployeeBankDetails) (*models.EmployeeBankDetails, error)
	GetBankDetails(ctx context.Context, empID uuid.UUID) (*models.EmployeeBankDetails, error)

	// Emergency Contacts
	AddEmergencyContact(ctx context.Context, c *models.EmergencyContact) (*models.EmergencyContact, error)
	ListEmergencyContacts(ctx context.Context, empID uuid.UUID) ([]*models.EmergencyContact, error)
	DeleteEmergencyContact(ctx context.Context, id uuid.UUID) error

	// Documents
	AddDocument(ctx context.Context, doc *models.EmployeeDocument) (*models.EmployeeDocument, error)
	ListDocuments(ctx context.Context, empID uuid.UUID) ([]*models.EmployeeDocument, error)
	VerifyDocument(ctx context.Context, docID, verifierID uuid.UUID) error
	DeleteDocument(ctx context.Context, docID uuid.UUID) error

	// Onboarding
	CreateOnboardingTask(ctx context.Context, task *models.OnboardingTask) (*models.OnboardingTask, error)
	ListOnboardingTasks(ctx context.Context, empID uuid.UUID) ([]*models.OnboardingTask, error)
	UpdateOnboardingTask(ctx context.Context, taskID uuid.UUID, status string) error
}

type employeeService struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

// ── Employee Core ────────────────────────────────────────────────────────────

func (s *employeeService) CreateEmployee(ctx context.Context, emp *models.Employee) (*models.Employee, error) {
	if emp.EmpCode == "" {
		seq, err := s.repo.GetNextEmpSeq(ctx, emp.OrgID)
		if err != nil {
			return nil, err
		}
		emp.EmpCode = utils.EmpCodeFromName(emp.FirstName, emp.LastName, seq)
	}

	if emp.Status == "" {
		emp.Status = "active"
	}
	if emp.EmploymentType == "" {
		emp.EmploymentType = "full_time"
	}

	if err := s.repo.Create(ctx, emp); err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *employeeService) GetEmployeeByID(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	emp, err := s.repo.GetByID(ctx, id)
	if err != nil || emp == nil {
		return nil, errors.New("employee not found")
	}
	return emp, nil
}

func (s *employeeService) GetEmployeeByUserID(ctx context.Context, userID uuid.UUID) (*models.Employee, error) {
	emp, err := s.repo.GetByUserID(ctx, userID)
	if err != nil || emp == nil {
		return nil, errors.New("employee record not found for user")
	}
	return emp, nil
}

func (s *employeeService) ListEmployees(ctx context.Context, orgID uuid.UUID, search string, deptID, desigID *uuid.UUID, status string, page, limit int) ([]*models.Employee, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.repo.ListByOrg(ctx, orgID, search, deptID, desigID, status, limit, offset)
}

func (s *employeeService) UpdateEmployee(ctx context.Context, emp *models.Employee) error {
	return s.repo.Update(ctx, emp)
}

func (s *employeeService) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// ── Departments ──────────────────────────────────────────────────────────────

func (s *employeeService) CreateDepartment(ctx context.Context, d *models.Department) (*models.Department, error) {
	if strings.TrimSpace(d.Name) == "" {
		return nil, errors.New("department name is required")
	}
	if err := s.repo.CreateDepartment(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *employeeService) GetDepartment(ctx context.Context, id uuid.UUID) (*models.Department, error) {
	d, err := s.repo.GetDepartmentByID(ctx, id)
	if err != nil || d == nil {
		return nil, errors.New("department not found")
	}
	return d, nil
}

func (s *employeeService) ListDepartments(ctx context.Context, orgID uuid.UUID) ([]*models.Department, error) {
	return s.repo.ListDepartmentsByOrg(ctx, orgID)
}

func (s *employeeService) UpdateDepartment(ctx context.Context, d *models.Department) error {
	return s.repo.UpdateDepartment(ctx, d)
}

func (s *employeeService) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteDepartment(ctx, id)
}

// ── Designations ─────────────────────────────────────────────────────────────

func (s *employeeService) CreateDesignation(ctx context.Context, des *models.Designation) (*models.Designation, error) {
	if strings.TrimSpace(des.Title) == "" {
		return nil, errors.New("designation title is required")
	}
	if err := s.repo.CreateDesignation(ctx, des); err != nil {
		return nil, err
	}
	return des, nil
}

func (s *employeeService) GetDesignation(ctx context.Context, id uuid.UUID) (*models.Designation, error) {
	des, err := s.repo.GetDesignationByID(ctx, id)
	if err != nil || des == nil {
		return nil, errors.New("designation not found")
	}
	return des, nil
}

func (s *employeeService) ListDesignations(ctx context.Context, orgID uuid.UUID) ([]*models.Designation, error) {
	return s.repo.ListDesignationsByOrg(ctx, orgID)
}

func (s *employeeService) UpdateDesignation(ctx context.Context, des *models.Designation) error {
	return s.repo.UpdateDesignation(ctx, des)
}

func (s *employeeService) DeleteDesignation(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteDesignation(ctx, id)
}

// ── Bank Details ─────────────────────────────────────────────────────────────

func (s *employeeService) SaveBankDetails(ctx context.Context, b *models.EmployeeBankDetails) (*models.EmployeeBankDetails, error) {
	if b.BankName == "" || b.AccountNoEnc == "" || b.IFSCCode == "" {
		return nil, errors.New("bank name, account number, and IFSC code are required")
	}
	// Mask account number for response
	rawAcc := b.AccountNoEnc
	if err := s.repo.SaveBankDetails(ctx, b); err != nil {
		return nil, err
	}
	if len(rawAcc) > 4 {
		b.MaskedAccountNo = fmt.Sprintf("XXXX-XXXX-%s", rawAcc[len(rawAcc)-4:])
	} else {
		b.MaskedAccountNo = "****"
	}
	return b, nil
}

func (s *employeeService) GetBankDetails(ctx context.Context, empID uuid.UUID) (*models.EmployeeBankDetails, error) {
	b, err := s.repo.GetBankDetailsByEmpID(ctx, empID)
	if err != nil || b == nil {
		return nil, errors.New("bank details not found")
	}
	if len(b.AccountNoEnc) > 4 {
		b.MaskedAccountNo = fmt.Sprintf("XXXX-XXXX-%s", b.AccountNoEnc[len(b.AccountNoEnc)-4:])
	} else {
		b.MaskedAccountNo = "****"
	}
	return b, nil
}

// ── Emergency Contacts ───────────────────────────────────────────────────────

func (s *employeeService) AddEmergencyContact(ctx context.Context, c *models.EmergencyContact) (*models.EmergencyContact, error) {
	if c.Name == "" || c.Phone == "" {
		return nil, errors.New("contact name and phone are required")
	}
	if err := s.repo.AddEmergencyContact(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *employeeService) ListEmergencyContacts(ctx context.Context, empID uuid.UUID) ([]*models.EmergencyContact, error) {
	return s.repo.ListEmergencyContacts(ctx, empID)
}

func (s *employeeService) DeleteEmergencyContact(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteEmergencyContact(ctx, id)
}

// ── Documents ────────────────────────────────────────────────────────────────

func (s *employeeService) AddDocument(ctx context.Context, doc *models.EmployeeDocument) (*models.EmployeeDocument, error) {
	if doc.DocType == "" || doc.FileURL == "" {
		return nil, errors.New("document type and file URL are required")
	}
	if err := s.repo.AddDocument(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *employeeService) ListDocuments(ctx context.Context, empID uuid.UUID) ([]*models.EmployeeDocument, error) {
	return s.repo.ListDocuments(ctx, empID)
}

func (s *employeeService) VerifyDocument(ctx context.Context, docID, verifierID uuid.UUID) error {
	return s.repo.VerifyDocument(ctx, docID, verifierID)
}

func (s *employeeService) DeleteDocument(ctx context.Context, docID uuid.UUID) error {
	return s.repo.DeleteDocument(ctx, docID)
}

// ── Onboarding Tasks ─────────────────────────────────────────────────────────

func (s *employeeService) CreateOnboardingTask(ctx context.Context, task *models.OnboardingTask) (*models.OnboardingTask, error) {
	if task.TaskName == "" {
		return nil, errors.New("task name is required")
	}
	if task.Status == "" {
		task.Status = "pending"
	}
	if err := s.repo.CreateOnboardingTask(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *employeeService) ListOnboardingTasks(ctx context.Context, empID uuid.UUID) ([]*models.OnboardingTask, error) {
	return s.repo.ListOnboardingTasks(ctx, empID)
}

func (s *employeeService) UpdateOnboardingTask(ctx context.Context, taskID uuid.UUID, status string) error {
	return s.repo.UpdateOnboardingTask(ctx, taskID, status)
}