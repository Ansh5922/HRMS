package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/repository"
)

type LeaveService interface {
	// Leave Types
	CreateLeaveType(ctx context.Context, lt *models.LeaveType) (*models.LeaveType, error)
	ListLeaveTypes(ctx context.Context, orgID uuid.UUID) ([]*models.LeaveType, error)
	UpdateLeaveType(ctx context.Context, lt *models.LeaveType) error
	DeleteLeaveType(ctx context.Context, id uuid.UUID) error

	// Leave Policies
	CreateLeavePolicy(ctx context.Context, lp *models.LeavePolicy) (*models.LeavePolicy, error)
	ListLeavePolicies(ctx context.Context, orgID uuid.UUID) ([]*models.LeavePolicy, error)
	DeleteLeavePolicy(ctx context.Context, id uuid.UUID) error

	// Leave Balances
	SetLeaveBalance(ctx context.Context, lb *models.LeaveBalance) (*models.LeaveBalance, error)
	GetLeaveBalances(ctx context.Context, empID uuid.UUID, year int) ([]*models.LeaveBalance, error)

	// Leave Applications
	ApplyLeave(ctx context.Context, app *models.LeaveApplication) (*models.LeaveApplication, error)
	GetLeaveByID(ctx context.Context, id uuid.UUID) (*models.LeaveApplication, error)
	ListLeavesByEmp(ctx context.Context, empID uuid.UUID) ([]*models.LeaveApplication, error)
	ListLeavesByOrg(ctx context.Context, orgID uuid.UUID, status string) ([]*models.LeaveApplication, error)
	ApproveLeave(ctx context.Context, appID, approverID uuid.UUID, comment *string) error
	RejectLeave(ctx context.Context, appID, approverID uuid.UUID, comment *string) error
	CancelLeave(ctx context.Context, appID, empID uuid.UUID) error
}

type leaveService struct {
	repo repository.LeaveRepository
}

func NewLeaveService(repo repository.LeaveRepository) LeaveService {
	return &leaveService{repo: repo}
}

// ── Leave Types ──────────────────────────────────────────────────────────────

func (s *leaveService) CreateLeaveType(ctx context.Context, lt *models.LeaveType) (*models.LeaveType, error) {
	if strings.TrimSpace(lt.Name) == "" || strings.TrimSpace(lt.Code) == "" {
		return nil, errors.New("leave type name and code are required")
	}
	lt.Code = strings.ToUpper(strings.TrimSpace(lt.Code))
	if lt.GenderApplicable == "" {
		lt.GenderApplicable = "all"
	}

	if err := s.repo.CreateLeaveType(ctx, lt); err != nil {
		return nil, err
	}
	return lt, nil
}

func (s *leaveService) ListLeaveTypes(ctx context.Context, orgID uuid.UUID) ([]*models.LeaveType, error) {
	return s.repo.ListLeaveTypesByOrg(ctx, orgID)
}

func (s *leaveService) UpdateLeaveType(ctx context.Context, lt *models.LeaveType) error {
	return s.repo.UpdateLeaveType(ctx, lt)
}

func (s *leaveService) DeleteLeaveType(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteLeaveType(ctx, id)
}

// ── Leave Policies ───────────────────────────────────────────────────────────

func (s *leaveService) CreateLeavePolicy(ctx context.Context, lp *models.LeavePolicy) (*models.LeavePolicy, error) {
	if lp.AnnualQuota <= 0 {
		return nil, errors.New("annual quota must be greater than 0")
	}
	if lp.AccrualType == "" {
		lp.AccrualType = "upfront"
	}
	if err := s.repo.CreateLeavePolicy(ctx, lp); err != nil {
		return nil, err
	}
	return lp, nil
}

func (s *leaveService) ListLeavePolicies(ctx context.Context, orgID uuid.UUID) ([]*models.LeavePolicy, error) {
	return s.repo.ListLeavePoliciesByOrg(ctx, orgID)
}

func (s *leaveService) DeleteLeavePolicy(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteLeavePolicy(ctx, id)
}

// ── Leave Balances ───────────────────────────────────────────────────────────

func (s *leaveService) SetLeaveBalance(ctx context.Context, lb *models.LeaveBalance) (*models.LeaveBalance, error) {
	if lb.Year == 0 {
		lb.Year = time.Now().Year()
	}
	if err := s.repo.SetLeaveBalance(ctx, lb); err != nil {
		return nil, err
	}
	return lb, nil
}

func (s *leaveService) GetLeaveBalances(ctx context.Context, empID uuid.UUID, year int) ([]*models.LeaveBalance, error) {
	if year == 0 {
		year = time.Now().Year()
	}
	return s.repo.GetLeaveBalanceByEmp(ctx, empID, year)
}

// ── Leave Applications ───────────────────────────────────────────────────────

func (s *leaveService) ApplyLeave(ctx context.Context, app *models.LeaveApplication) (*models.LeaveApplication, error) {
	if app.FromDate.After(app.ToDate) {
		return nil, errors.New("from_date cannot be after to_date")
	}

	// Calculate leave days
	if app.Session == "first_half" || app.Session == "second_half" {
		app.Days = 0.5
	} else {
		app.Session = "full"
		diff := app.ToDate.Sub(app.FromDate)
		app.Days = float64(int(diff.Hours()/24) + 1)
	}

	// Get leave type configuration
	lt, err := s.repo.GetLeaveTypeByID(ctx, app.LeaveTypeID)
	if err != nil || lt == nil {
		return nil, errors.New("invalid leave type")
	}

	// Check if document is required
	if lt.RequiresDoc && (app.DocURL == nil || *app.DocURL == "") {
		return nil, fmt.Errorf("leave type '%s' requires supporting document upload", lt.Name)
	}

	// Check Leave Balance
	currentYear := app.FromDate.Year()
	balance, err := s.repo.GetLeaveBalanceForType(ctx, app.EmpID, app.LeaveTypeID, currentYear)
	if err != nil {
		return nil, err
	}
	if balance != nil {
		available := balance.Total - balance.Used - balance.Pending
		if available < app.Days {
			return nil, fmt.Errorf("insufficient leave balance: requested %.1f days, but only %.1f available", app.Days, available)
		}
	}

	// Auto-approve if leave type doesn't require approval
	if !lt.RequiresApproval {
		app.Status = "approved"
	} else {
		app.Status = "pending"
	}

	if err := s.repo.Apply(ctx, app); err != nil {
		return nil, err
	}

	// Adjust balance
	if app.Status == "pending" && balance != nil {
		_ = s.repo.AdjustLeaveBalanceOnApply(ctx, app.EmpID, app.LeaveTypeID, currentYear, app.Days)
	} else if app.Status == "approved" && balance != nil {
		_ = s.repo.AdjustLeaveBalanceOnApprove(ctx, app.EmpID, app.LeaveTypeID, currentYear, app.Days)
	}

	return app, nil
}

func (s *leaveService) GetLeaveByID(ctx context.Context, id uuid.UUID) (*models.LeaveApplication, error) {
	app, err := s.repo.GetByID(ctx, id)
	if err != nil || app == nil {
		return nil, errors.New("leave application not found")
	}
	return app, nil
}

func (s *leaveService) ListLeavesByEmp(ctx context.Context, empID uuid.UUID) ([]*models.LeaveApplication, error) {
	return s.repo.ListByEmp(ctx, empID)
}

func (s *leaveService) ListLeavesByOrg(ctx context.Context, orgID uuid.UUID, status string) ([]*models.LeaveApplication, error) {
	return s.repo.ListByOrg(ctx, orgID, status)
}

func (s *leaveService) ApproveLeave(ctx context.Context, appID, approverID uuid.UUID, comment *string) error {
	app, err := s.repo.GetByID(ctx, appID)
	if err != nil || app == nil {
		return errors.New("leave application not found")
	}
	if app.Status != "pending" {
		return fmt.Errorf("cannot approve leave with status '%s'", app.Status)
	}

	// Record approval audit
	approval := &models.LeaveApproval{
		ApplicationID: appID,
		ApproverID:    approverID,
		Level:         1,
		Action:        "approved",
		Comment:       comment,
	}
	if err := s.repo.AddApprovalRecord(ctx, approval); err != nil {
		return err
	}

	// Update status to approved
	if err := s.repo.UpdateStatus(ctx, appID, "approved"); err != nil {
		return err
	}

	// Adjust balance (deduct pending, add to used)
	_ = s.repo.AdjustLeaveBalanceOnApprove(ctx, app.EmpID, app.LeaveTypeID, app.FromDate.Year(), app.Days)
	return nil
}

func (s *leaveService) RejectLeave(ctx context.Context, appID, approverID uuid.UUID, comment *string) error {
	app, err := s.repo.GetByID(ctx, appID)
	if err != nil || app == nil {
		return errors.New("leave application not found")
	}
	if app.Status != "pending" {
		return fmt.Errorf("cannot reject leave with status '%s'", app.Status)
	}

	approval := &models.LeaveApproval{
		ApplicationID: appID,
		ApproverID:    approverID,
		Level:         1,
		Action:        "rejected",
		Comment:       comment,
	}
	if err := s.repo.AddApprovalRecord(ctx, approval); err != nil {
		return err
	}

	if err := s.repo.UpdateStatus(ctx, appID, "rejected"); err != nil {
		return err
	}

	// Restore pending balance
	_ = s.repo.AdjustLeaveBalanceOnReject(ctx, app.EmpID, app.LeaveTypeID, app.FromDate.Year(), app.Days)
	return nil
}

func (s *leaveService) CancelLeave(ctx context.Context, appID, empID uuid.UUID) error {
	app, err := s.repo.GetByID(ctx, appID)
	if err != nil || app == nil {
		return errors.New("leave application not found")
	}
	if app.EmpID != empID {
		return errors.New("unauthorized to cancel this leave application")
	}
	if app.Status != "pending" {
		return errors.New("only pending leave applications can be cancelled")
	}

	if err := s.repo.UpdateStatus(ctx, appID, "cancelled"); err != nil {
		return err
	}

	// Restore pending balance
	_ = s.repo.AdjustLeaveBalanceOnReject(ctx, app.EmpID, app.LeaveTypeID, app.FromDate.Year(), app.Days)
	return nil
}
