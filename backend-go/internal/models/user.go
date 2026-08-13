package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID  `json:"id"`
	OrgID          uuid.UUID  `json:"org_id"`
	Email          string     `json:"email"`
	PasswordHash   string     `json:"-"`
	RoleID         *uuid.UUID `json:"role_id"`
	IsActive       bool       `json:"is_active"`
	IsVerified     bool       `json:"is_verified"`
	MFAEnabled     bool       `json:"mfa_enabled"`
	MFASecret      *string    `json:"-"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	FailedAttempts int        `json:"-"`
	LockedUntil    *time.Time `json:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Joined/computed fields
	RoleName    string   `json:"role_name,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type Role struct {
	ID          uuid.UUID    `json:"id"`
	OrgID       uuid.UUID    `json:"org_id"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	IsSystem    bool         `json:"is_system"`
	CreatedAt   time.Time    `json:"created_at"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	ID       uuid.UUID `json:"id"`
	Resource string    `json:"resource"`
	Action   string    `json:"action"`
}

type RefreshToken struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	TokenHash  string    `json:"-"`
	DeviceInfo *string   `json:"device_info"`
	IPAddress  *string   `json:"ip_address"`
	ExpiresAt  time.Time `json:"expires_at"`
	Revoked    bool      `json:"revoked"`
	CreatedAt  time.Time `json:"created_at"`
}

type PasswordResetToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditLog struct {
	ID         int64      `json:"id"`
	OrgID      uuid.UUID  `json:"org_id"`
	UserID     *uuid.UUID `json:"user_id"`
	Action     string     `json:"action"`
	Resource   string     `json:"resource"`
	ResourceID *uuid.UUID `json:"resource_id"`
	OldData    *string    `json:"old_data"`
	NewData    *string    `json:"new_data"`
	IPAddress  *string    `json:"ip_address"`
	UserAgent  *string    `json:"user_agent"`
	CreatedAt  time.Time  `json:"created_at"`
}