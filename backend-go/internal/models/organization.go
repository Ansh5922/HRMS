package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	LogoURL      *string   `json:"logo_url"`
	Address      *string   `json:"address"`
	Timezone     string    `json:"timezone"`
	FinYearStart int       `json:"fin_year_start"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}