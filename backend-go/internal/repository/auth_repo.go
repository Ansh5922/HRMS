package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/hrms-backend/internal/models"
)

type AuthRepository interface {
	// Organization
	CreateOrganization(ctx context.Context, name string) (*models.Organization, error)

	// Users
	CreateUser(ctx context.Context, u *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListUsersByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.User, error)
	UpdateUser(ctx context.Context, u *models.User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, hash string) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error
	ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error
	LockUser(ctx context.Context, userID uuid.UUID, until time.Time) error
	DeactivateUser(ctx context.Context, userID uuid.UUID) error

	// Refresh Tokens
	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error

	// Password Reset
	SavePasswordResetToken(ctx context.Context, t *models.PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error)
	MarkResetTokenUsed(ctx context.Context, id uuid.UUID) error

	// Roles
	CreateRole(ctx context.Context, role *models.Role) error
	GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	ListRolesByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.Role, error)
	UpdateRole(ctx context.Context, role *models.Role) error
	DeleteRole(ctx context.Context, id uuid.UUID) error
	AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error

	// Permissions
	ListAllPermissions(ctx context.Context) ([]*models.Permission, error)
	GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]string, error)
	SetRolePermissions(ctx context.Context, roleID uuid.UUID, permIDs []uuid.UUID) error

	// Audit
	CreateAuditLog(ctx context.Context, log *models.AuditLog) error
}

type authRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepo{db: db}
}

// HashToken produces a SHA-256 hex digest used for refresh & reset tokens.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// ── Organization ─────────────────────────────────────────────────────────────

func (r *authRepo) CreateOrganization(ctx context.Context, name string) (*models.Organization, error) {
	query := `INSERT INTO organizations (name) VALUES ($1)
		RETURNING id, name, timezone, fin_year_start, is_active, created_at, updated_at`
	org := &models.Organization{}
	err := r.db.QueryRow(ctx, query, name).Scan(
		&org.ID, &org.Name, &org.Timezone, &org.FinYearStart, &org.IsActive, &org.CreatedAt, &org.UpdatedAt,
	)
	return org, err
}

// ── Users ────────────────────────────────────────────────────────────────────

func (r *authRepo) CreateUser(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (org_id, email, password_hash, role_id, is_active, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, u.OrgID, u.Email, u.PasswordHash, u.RoleID, u.IsActive, u.IsVerified).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT u.id, u.org_id, u.email, u.password_hash, u.role_id,
		u.is_active, u.is_verified, u.mfa_enabled, u.failed_attempts, u.locked_until,
		u.created_at, u.updated_at, COALESCE(r.name, '') as role_name
		FROM users u LEFT JOIN roles r ON r.id = u.role_id
		WHERE u.email = $1`
	u := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.OrgID, &u.Email, &u.PasswordHash, &u.RoleID,
		&u.IsActive, &u.IsVerified, &u.MFAEnabled, &u.FailedAttempts, &u.LockedUntil,
		&u.CreatedAt, &u.UpdatedAt, &u.RoleName,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *authRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT u.id, u.org_id, u.email, u.password_hash, u.role_id,
		u.is_active, u.is_verified, u.mfa_enabled, u.failed_attempts, u.locked_until,
		u.created_at, u.updated_at, COALESCE(r.name, '') as role_name
		FROM users u LEFT JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1`
	u := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.OrgID, &u.Email, &u.PasswordHash, &u.RoleID,
		&u.IsActive, &u.IsVerified, &u.MFAEnabled, &u.FailedAttempts, &u.LockedUntil,
		&u.CreatedAt, &u.UpdatedAt, &u.RoleName,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *authRepo) ListUsersByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.User, error) {
	query := `SELECT u.id, u.org_id, u.email, u.role_id, u.is_active, u.is_verified,
		u.created_at, u.updated_at, COALESCE(r.name, '') as role_name
		FROM users u LEFT JOIN roles r ON r.id = u.role_id
		WHERE u.org_id = $1 ORDER BY u.created_at DESC`
	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.OrgID, &u.Email, &u.RoleID, &u.IsActive, &u.IsVerified,
			&u.CreatedAt, &u.UpdatedAt, &u.RoleName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *authRepo) UpdateUser(ctx context.Context, u *models.User) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET email=$1, is_active=$2, is_verified=$3, role_id=$4, updated_at=NOW() WHERE id=$5`,
		u.Email, u.IsActive, u.IsVerified, u.RoleID, u.ID,
	)
	return err
}

func (r *authRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, hash string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET password_hash=$1, updated_at=NOW() WHERE id=$2`, hash, userID)
	return err
}

func (r *authRepo) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET last_login_at=NOW(), failed_attempts=0 WHERE id=$1`, userID)
	return err
}

func (r *authRepo) IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET failed_attempts = failed_attempts + 1 WHERE id=$1`, userID)
	return err
}

func (r *authRepo) ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET failed_attempts=0, locked_until=NULL WHERE id=$1`, userID)
	return err
}

func (r *authRepo) LockUser(ctx context.Context, userID uuid.UUID, until time.Time) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET locked_until=$1 WHERE id=$2`, until, userID)
	return err
}

func (r *authRepo) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET is_active=FALSE, updated_at=NOW() WHERE id=$1`, userID)
	return err
}

// ── Refresh Tokens ───────────────────────────────────────────────────────────

func (r *authRepo) SaveRefreshToken(ctx context.Context, t *models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, device_info, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, t.UserID, t.TokenHash, t.DeviceInfo, t.IPAddress, t.ExpiresAt).
		Scan(&t.ID, &t.CreatedAt)
}

func (r *authRepo) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, revoked, created_at
		FROM refresh_tokens WHERE token_hash=$1`
	t := &models.RefreshToken{}
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.Revoked, &t.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func (r *authRepo) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx, `UPDATE refresh_tokens SET revoked=TRUE WHERE token_hash=$1`, tokenHash)
	return err
}

func (r *authRepo) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE refresh_tokens SET revoked=TRUE WHERE user_id=$1 AND revoked=FALSE`, userID)
	return err
}

// ── Password Reset ───────────────────────────────────────────────────────────

func (r *authRepo) SavePasswordResetToken(ctx context.Context, t *models.PasswordResetToken) error {
	query := `INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, t.UserID, t.TokenHash, t.ExpiresAt).Scan(&t.ID, &t.CreatedAt)
}

func (r *authRepo) GetPasswordResetToken(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error) {
	query := `SELECT id, user_id, expires_at, used, created_at
		FROM password_reset_tokens WHERE token_hash=$1`
	t := &models.PasswordResetToken{}
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&t.ID, &t.UserID, &t.ExpiresAt, &t.Used, &t.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func (r *authRepo) MarkResetTokenUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE password_reset_tokens SET used=TRUE WHERE id=$1`, id)
	return err
}

// ── Roles ────────────────────────────────────────────────────────────────────

func (r *authRepo) CreateRole(ctx context.Context, role *models.Role) error {
	query := `INSERT INTO roles (org_id, name, description, is_system)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, role.OrgID, role.Name, role.Description, role.IsSystem).
		Scan(&role.ID, &role.CreatedAt)
}

func (r *authRepo) GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	query := `SELECT id, org_id, name, description, is_system, created_at FROM roles WHERE id=$1`
	role := &models.Role{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&role.ID, &role.OrgID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return role, err
}

func (r *authRepo) ListRolesByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.Role, error) {
	query := `SELECT id, org_id, name, description, is_system, created_at FROM roles WHERE org_id=$1 ORDER BY name`
	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		role := &models.Role{}
		if err := rows.Scan(&role.ID, &role.OrgID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *authRepo) UpdateRole(ctx context.Context, role *models.Role) error {
	_, err := r.db.Exec(ctx, `UPDATE roles SET name=$1, description=$2 WHERE id=$3`,
		role.Name, role.Description, role.ID)
	return err
}

func (r *authRepo) DeleteRole(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM roles WHERE id=$1 AND is_system=FALSE`, id)
	return err
}

func (r *authRepo) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET role_id=$1, updated_at=NOW() WHERE id=$2`, roleID, userID)
	return err
}

// ── Permissions ──────────────────────────────────────────────────────────────

func (r *authRepo) ListAllPermissions(ctx context.Context) ([]*models.Permission, error) {
	query := `SELECT id, resource, action FROM permissions ORDER BY resource, action`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []*models.Permission
	for rows.Next() {
		p := &models.Permission{}
		if err := rows.Scan(&p.ID, &p.Resource, &p.Action); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *authRepo) GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	query := `SELECT p.resource || ':' || p.action
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1`
	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	return perms, nil
}

func (r *authRepo) SetRolePermissions(ctx context.Context, roleID uuid.UUID, permIDs []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM role_permissions WHERE role_id=$1`, roleID); err != nil {
		return err
	}

	for _, pid := range permIDs {
		if _, err := tx.Exec(ctx,
			`INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			roleID, pid,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// ── Audit ────────────────────────────────────────────────────────────────────

func (r *authRepo) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	query := `INSERT INTO audit_logs (org_id, user_id, action, resource, resource_id, old_data, new_data, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(ctx, query,
		log.OrgID, log.UserID, log.Action, log.Resource, log.ResourceID,
		log.OldData, log.NewData, log.IPAddress, log.UserAgent,
	)
	return err
}

// FormatPermission builds the "resource:action" string used in RBAC checks.
func FormatPermission(resource, action string) string {
	return fmt.Sprintf("%s:%s", resource, action)
}