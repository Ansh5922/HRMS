package models

import (
	"time"

	"github.com/google/uuid"
)

type JobPosting struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       uuid.UUID  `json:"org_id"`
	Title       string     `json:"title"`
	DeptID      *uuid.UUID `json:"dept_id"`
	Location    *string    `json:"location"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Candidate struct {
	ID        uuid.UUID `json:"id"`
	OrgID     uuid.UUID `json:"org_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone"`
	ResumeURL *string   `json:"resume_url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}