package models

import (
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
	ID            uuid.UUID  `json:"id"`
	OrgID         uuid.UUID  `json:"org_id"`
	EmpID         uuid.UUID  `json:"emp_id"`
	Date          time.Time  `json:"date"`
	CheckIn       *time.Time `json:"check_in"`
	CheckOut      *time.Time `json:"check_out"`
	Method        *string    `json:"method"`
	Status        *string    `json:"status"`
	DurationMins  *int       `json:"duration_mins"`
	OvertimeMins  int        `json:"overtime_mins"`
	LocationLat   *float64   `json:"location_lat"`
	LocationLng   *float64   `json:"location_lng"`
	FaceConfScore *float64   `json:"face_conf_score"`
	IsRegularized bool       `json:"is_regularized"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Shift struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Name        string    `json:"name"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	GraceMins   int       `json:"grace_mins"`
	IsOvernight bool      `json:"is_overnight"`
	CreatedAt   time.Time `json:"created_at"`
}