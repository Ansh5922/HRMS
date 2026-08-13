package models

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID               uuid.UUID  `json:"id"`
	OrgID            uuid.UUID  `json:"org_id"`
	UserID           *uuid.UUID `json:"user_id"`
	EmpCode          string     `json:"emp_code"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	MiddleName       *string    `json:"middle_name"`
	DOB              *time.Time `json:"dob"`
	Gender           *string    `json:"gender"`
	BloodGroup       *string    `json:"blood_group"`
	Phone            *string    `json:"phone"`
	PersonalEmail    *string    `json:"personal_email"`
	Address          *string    `json:"address"`
	City             *string    `json:"city"`
	State            *string    `json:"state"`
	Country          *string    `json:"country"`
	Pincode          *string    `json:"pincode"`
	Nationality      *string    `json:"nationality"`
	AadhaarNoEnc     *string    `json:"-"`
	PanNoEnc         *string    `json:"-"`
	PassportNo       *string    `json:"passport_no"`
	DeptID           *uuid.UUID `json:"dept_id"`
	DesignationID    *uuid.UUID `json:"designation_id"`
	ManagerID        *uuid.UUID `json:"manager_id"`
	EmploymentType   string     `json:"employment_type"`
	WorkLocation     *string    `json:"work_location"`
	JoiningDate      time.Time  `json:"joining_date"`
	ConfirmationDate *time.Time `json:"confirmation_date"`
	Status           string     `json:"status"`
	ExitDate         *time.Time `json:"exit_date"`
	ExitReason       *string    `json:"exit_reason"`
	PhotoURL         *string    `json:"photo_url"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Joined fields
	DepartmentName  *string `json:"department_name,omitempty"`
	DesignationName *string `json:"designation_name,omitempty"`
	ManagerName     *string `json:"manager_name,omitempty"`
}

type Department struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       uuid.UUID  `json:"org_id"`
	Name        string     `json:"name"`
	Code        *string    `json:"code"`
	ParentDeptID *uuid.UUID `json:"parent_dept_id"`
	HeadEmpID   *uuid.UUID `json:"head_emp_id"`
	Description *string    `json:"description"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Joined fields
	HeadEmpName *string `json:"head_emp_name,omitempty"`
	ParentName  *string `json:"parent_name,omitempty"`
}

type Designation struct {
	ID           uuid.UUID  `json:"id"`
	OrgID        uuid.UUID  `json:"org_id"`
	DepartmentID *uuid.UUID `json:"department_id"`
	Title        string     `json:"title"`
	Grade        *string    `json:"grade"`
	Level        *int       `json:"level"`
	CreatedAt    time.Time  `json:"created_at"`

	// Joined
	DepartmentName *string `json:"department_name,omitempty"`
}

type EmployeeBankDetails struct {
	ID           uuid.UUID `json:"id"`
	EmpID        uuid.UUID `json:"emp_id"`
	BankName     string    `json:"bank_name"`
	AccountNoEnc string    `json:"-"`
	IFSCCode     string    `json:"ifsc_code"`
	AccountType  *string   `json:"account_type"`
	IsPrimary    bool      `json:"is_primary"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Display field (masked account number for response)
	MaskedAccountNo string `json:"account_number,omitempty"`
}

type EmergencyContact struct {
	ID       uuid.UUID `json:"id"`
	EmpID    uuid.UUID `json:"emp_id"`
	Name     string    `json:"name"`
	Relation *string   `json:"relation"`
	Phone    string    `json:"phone"`
	Email    *string   `json:"email"`
	Address  *string   `json:"address"`
}

type EmployeeDocument struct {
	ID         uuid.UUID  `json:"id"`
	EmpID      uuid.UUID  `json:"emp_id"`
	DocType    string     `json:"doc_type"`
	FileURL    string     `json:"file_url"`
	FileName   *string    `json:"file_name"`
	IsVerified bool       `json:"is_verified"`
	VerifiedBy *uuid.UUID `json:"verified_by"`
	VerifiedAt *time.Time `json:"verified_at"`
	UploadedAt time.Time  `json:"uploaded_at"`
}

type OnboardingTask struct {
	ID          uuid.UUID  `json:"id"`
	EmpID       uuid.UUID  `json:"emp_id"`
	TaskName    string     `json:"task_name"`
	Description *string    `json:"description"`
	AssignedTo  *uuid.UUID `json:"assigned_to"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	Status      string     `json:"status"`
}