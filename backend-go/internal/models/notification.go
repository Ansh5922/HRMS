package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID  `json:"id"`
	OrgID     uuid.UUID  `json:"org_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Message   *string    `json:"message"`
	ActionURL *string    `json:"action_url"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type Announcement struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       uuid.UUID  `json:"org_id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	AuthorID    *uuid.UUID `json:"author_id"`
	TargetDept  *uuid.UUID `json:"target_dept"`
	PublishedAt time.Time  `json:"published_at"`
	ExpiresAt   *time.Time `json:"expires_at"`

	// Joined
	AuthorName     *string `json:"author_name,omitempty"`
	TargetDeptName *string `json:"target_dept_name,omitempty"`
}

type NotificationPreference struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	EventType string    `json:"event_type"`
	Email     bool      `json:"email"`
	Push      bool      `json:"push"`
	SMS       bool      `json:"sms"`
}

type PushToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"created_at"`
}