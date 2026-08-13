package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	ContentURL  *string   `json:"content_url"`
	CreatedAt   time.Time `json:"created_at"`
}