package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/config"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/repository"
	pkgjwt "github.com/your-org/hrms-backend/pkg/jwt"
	"github.com/your-org/hrms-backend/pkg/utils"
)

const (
	maxFailedAttempts = 5
	lockDuration      = 30 * time.Minute
	accessTokenTTL    = 15 * time.Minute
	refreshTokenTTL   = 7 * 24 * time.Hour
	resetTokenTTL     = 1 * time.Hour
)

// ── Interface ────────────────────────────────────────────────────────────────

type AuthService interface {
	// Auth flows
	RegisterOrgAndAdmin(ctx context.Context, orgName, email, password string) (*RegisterResponse, error)
	Login(ctx context.Context, email, password, ipAddress, userAgent string) (*LoginResponse, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*LoginResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID uuid.UUID) error

	// Profile
	GetMe(ctx context.Context, userID uuid.UUID) (*UserProfileResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPass, newPass string) error

	// Password reset
	ForgotPassword(ctx context.Context, email string) (*ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, token, newPassword string) error

	// User management (admin)
	CreateUser(ctx context.Context, orgID uuid.UUID, email, password string, roleID *uuid.UUID) (*models.User, error)
	ListUsers(ctx context.Context, orgID uuid.UUID) ([]*models.User, error)
	UpdateUser(ctx context.Context, u *models.User) error
	DeactivateUser(ctx context.Context, userID uuid.UUID) error

	// Role management
	CreateRole(ctx context.Context, role *models.Role) error
	ListRoles(ctx context.Context, orgID uuid.UUID) ([]*models.Role, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (*models.Role, error)
	UpdateRole(ctx context.Context, role *models.Role) error
	DeleteRole(ctx context.Context, roleID uuid.UUID) error
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error

	// Permission management
	ListPermissions(ctx context.Context) ([]*models.Permission, error)
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]string, error)
	SetRolePermissions(ctx context.Context, roleID uuid.UUID, permIDs []uuid.UUID) error
}

// ── Response types ───────────────────────────────────────────────────────────

type RegisterResponse struct {
	Organization *models.Organization `json:"organization"`
	AdminUser    *models.User         `json:"admin_user"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *models.User `json:"user"`
}

type UserProfileResponse struct {
	User        *models.User `json:"user"`
	Permissions []string     `json:"permissions"`
}

type ForgotPasswordResponse struct {
	Message    string `json:"message"`
	ResetToken string `json:"reset_token"` // In production, send via email instead
}

// ── Implementation ───────────────────────────────────────────────────────────

type authService struct {
	repo repository.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo repository.AuthRepository, cfg *config.Config) AuthService {
	return &authService{repo: repo, cfg: cfg}
}

// ── Auth flows ───────────────────────────────────────────────────────────────

func (s *authService) RegisterOrgAndAdmin(ctx context.Context, orgName, email, password string) (*RegisterResponse, error) {
	// Check duplicate
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create org
	org, err := s.repo.CreateOrganization(ctx, orgName)
	if err != nil {
		return nil, err
	}

	// Create super_admin role for the org
	superRole := &models.Role{
		OrgID:    org.ID,
		Name:     "super_admin",
		IsSystem: true,
	}
	if err := s.repo.CreateRole(ctx, superRole); err != nil {
		return nil, err
	}

	// Assign ALL permissions to super_admin
	allPerms, err := s.repo.ListAllPermissions(ctx)
	if err != nil {
		return nil, err
	}
	permIDs := make([]uuid.UUID, len(allPerms))
	for i, p := range allPerms {
		permIDs[i] = p.ID
	}
	if err := s.repo.SetRolePermissions(ctx, superRole.ID, permIDs); err != nil {
		return nil, err
	}

	// Hash password
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create admin user with super_admin role
	user := &models.User{
		OrgID:        org.ID,
		Email:        email,
		PasswordHash: hash,
		RoleID:       &superRole.ID,
		IsActive:     true,
		IsVerified:   true,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Get permissions for token
	permissions, _ := s.repo.GetPermissionsByRoleID(ctx, superRole.ID)

	// Generate tokens
	accessToken, err := pkgjwt.GenerateToken(
		user.ID.String(), org.ID.String(),
		superRole.ID.String(), superRole.Name,
		permissions, s.cfg.AppSecret, accessTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := s.createRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	user.RoleName = "super_admin"
	user.Permissions = permissions

	return &RegisterResponse{
		Organization: org,
		AdminUser:    user,
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (s *authService) Login(ctx context.Context, email, password, ipAddress, userAgent string) (*LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if account is locked
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return nil, errors.New("account is temporarily locked, try again later")
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Verify password
	if !utils.CheckPassword(user.PasswordHash, password) {
		// Increment failed attempts
		_ = s.repo.IncrementFailedAttempts(ctx, user.ID)
		if user.FailedAttempts+1 >= maxFailedAttempts {
			_ = s.repo.LockUser(ctx, user.ID, time.Now().Add(lockDuration))
		}
		return nil, errors.New("invalid email or password")
	}

	// Reset failed attempts & update last login
	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	// Load permissions
	var permissions []string
	roleID := ""
	roleName := user.RoleName
	if user.RoleID != nil {
		roleID = user.RoleID.String()
		permissions, _ = s.repo.GetPermissionsByRoleID(ctx, *user.RoleID)
	}

	// Generate access token
	accessToken, err := pkgjwt.GenerateToken(
		user.ID.String(), user.OrgID.String(),
		roleID, roleName, permissions,
		s.cfg.AppSecret, accessTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenStr, err := s.createRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// Audit log
	_ = s.repo.CreateAuditLog(ctx, &models.AuditLog{
		OrgID:     user.OrgID,
		UserID:    &user.ID,
		Action:    "login",
		Resource:  "auth",
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
	})

	user.Permissions = permissions

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		User:         user,
	}, nil
}

func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	tokenHash := repository.HashToken(refreshToken)
	storedToken, err := s.repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if storedToken == nil {
		return nil, errors.New("invalid refresh token")
	}
	if storedToken.Revoked {
		// Possible token reuse attack — revoke all tokens for this user
		_ = s.repo.RevokeAllUserTokens(ctx, storedToken.UserID)
		return nil, errors.New("refresh token has been revoked")
	}
	if storedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	// Revoke old token (rotation)
	_ = s.repo.RevokeRefreshToken(ctx, tokenHash)

	// Load user with role
	user, err := s.repo.GetUserByID(ctx, storedToken.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Load permissions
	var permissions []string
	roleID := ""
	roleName := user.RoleName
	if user.RoleID != nil {
		roleID = user.RoleID.String()
		permissions, _ = s.repo.GetPermissionsByRoleID(ctx, *user.RoleID)
	}

	// New access token
	accessToken, err := pkgjwt.GenerateToken(
		user.ID.String(), user.OrgID.String(),
		roleID, roleName, permissions,
		s.cfg.AppSecret, accessTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	// New refresh token
	newRefreshStr, err := s.createRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	user.Permissions = permissions

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshStr,
		User:         user,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := repository.HashToken(refreshToken)
	return s.repo.RevokeRefreshToken(ctx, tokenHash)
}

func (s *authService) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return s.repo.RevokeAllUserTokens(ctx, userID)
}

// ── Profile ──────────────────────────────────────────────────────────────────

func (s *authService) GetMe(ctx context.Context, userID uuid.UUID) (*UserProfileResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	var permissions []string
	if user.RoleID != nil {
		permissions, _ = s.repo.GetPermissionsByRoleID(ctx, *user.RoleID)
	}

	return &UserProfileResponse{
		User:        user,
		Permissions: permissions,
	}, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPass, newPass string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	if !utils.CheckPassword(user.PasswordHash, oldPass) {
		return errors.New("current password is incorrect")
	}

	hash, err := utils.HashPassword(newPass)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, userID, hash); err != nil {
		return err
	}

	// Revoke all refresh tokens (force re-login on other devices)
	_ = s.repo.RevokeAllUserTokens(ctx, userID)
	return nil
}

// ── Password Reset ───────────────────────────────────────────────────────────

func (s *authService) ForgotPassword(ctx context.Context, email string) (*ForgotPasswordResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	// Always return success to prevent email enumeration
	if user == nil {
		return &ForgotPasswordResponse{
			Message: "If that email exists, a reset link has been sent",
		}, nil
	}

	rawToken := utils.GenerateToken(32)
	tokenHash := repository.HashToken(rawToken)

	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(resetTokenTTL),
	}
	if err := s.repo.SavePasswordResetToken(ctx, resetToken); err != nil {
		return nil, err
	}

	// In production: send email with rawToken embedded in a URL
	// For dev: return raw token directly
	return &ForgotPasswordResponse{
		Message:    "If that email exists, a reset link has been sent",
		ResetToken: rawToken,
	}, nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	tokenHash := repository.HashToken(token)
	resetToken, err := s.repo.GetPasswordResetToken(ctx, tokenHash)
	if err != nil {
		return err
	}
	if resetToken == nil {
		return errors.New("invalid or expired reset token")
	}
	if resetToken.Used {
		return errors.New("reset token has already been used")
	}
	if resetToken.ExpiresAt.Before(time.Now()) {
		return errors.New("reset token has expired")
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, resetToken.UserID, hash); err != nil {
		return err
	}

	_ = s.repo.MarkResetTokenUsed(ctx, resetToken.ID)
	_ = s.repo.RevokeAllUserTokens(ctx, resetToken.UserID)
	_ = s.repo.ResetFailedAttempts(ctx, resetToken.UserID)

	return nil
}

// ── User management ──────────────────────────────────────────────────────────

func (s *authService) CreateUser(ctx context.Context, orgID uuid.UUID, email, password string, roleID *uuid.UUID) (*models.User, error) {
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		OrgID:        orgID,
		Email:        email,
		PasswordHash: hash,
		RoleID:       roleID,
		IsActive:     true,
		IsVerified:   false,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) ListUsers(ctx context.Context, orgID uuid.UUID) ([]*models.User, error) {
	return s.repo.ListUsersByOrg(ctx, orgID)
}

func (s *authService) UpdateUser(ctx context.Context, u *models.User) error {
	return s.repo.UpdateUser(ctx, u)
}

func (s *authService) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	_ = s.repo.RevokeAllUserTokens(ctx, userID)
	return s.repo.DeactivateUser(ctx, userID)
}

// ── Role management ──────────────────────────────────────────────────────────

func (s *authService) CreateRole(ctx context.Context, role *models.Role) error {
	return s.repo.CreateRole(ctx, role)
}

func (s *authService) ListRoles(ctx context.Context, orgID uuid.UUID) ([]*models.Role, error) {
	roles, err := s.repo.ListRolesByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}
	// Attach permissions to each role
	for _, role := range roles {
		perms, _ := s.repo.GetPermissionsByRoleID(ctx, role.ID)
		for _, p := range perms {
			role.Permissions = append(role.Permissions, models.Permission{Resource: p})
		}
	}
	return roles, nil
}

func (s *authService) GetRole(ctx context.Context, roleID uuid.UUID) (*models.Role, error) {
	role, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil || role == nil {
		return nil, errors.New("role not found")
	}
	permStrs, _ := s.repo.GetPermissionsByRoleID(ctx, roleID)
	for _, p := range permStrs {
		role.Permissions = append(role.Permissions, models.Permission{Resource: p})
	}
	return role, nil
}

func (s *authService) UpdateRole(ctx context.Context, role *models.Role) error {
	return s.repo.UpdateRole(ctx, role)
}

func (s *authService) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	return s.repo.DeleteRole(ctx, roleID)
}

func (s *authService) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.repo.AssignRoleToUser(ctx, userID, roleID)
}

// ── Permission management ────────────────────────────────────────────────────

func (s *authService) ListPermissions(ctx context.Context) ([]*models.Permission, error) {
	return s.repo.ListAllPermissions(ctx)
}

func (s *authService) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	return s.repo.GetPermissionsByRoleID(ctx, roleID)
}

func (s *authService) SetRolePermissions(ctx context.Context, roleID uuid.UUID, permIDs []uuid.UUID) error {
	return s.repo.SetRolePermissions(ctx, roleID, permIDs)
}

// ── Internal helpers ─────────────────────────────────────────────────────────

func (s *authService) createRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	rawToken := utils.GenerateToken(32)
	tokenHash := repository.HashToken(rawToken)

	refModel := &models.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}
	if err := s.repo.SaveRefreshToken(ctx, refModel); err != nil {
		return "", err
	}
	return rawToken, nil
}