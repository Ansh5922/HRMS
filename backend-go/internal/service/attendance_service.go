package service

import (
"context"
"time"

"github.com/google/uuid"
"github.com/your-org/hrms-backend/internal/models"
"github.com/your-org/hrms-backend/internal/repository"
)

type AttendanceService interface {
CheckIn(ctx context.Context, orgID, empID uuid.UUID, method string, lat, lng *float64) (*models.Attendance, error)
CheckOut(ctx context.Context, empID uuid.UUID) error
GetToday(ctx context.Context, orgID, dateStr string) ([]*models.Attendance, error)
}

type attendanceService struct {
repo repository.AttendanceRepository
}

func NewAttendanceService(repo repository.AttendanceRepository) AttendanceService {
return &attendanceService{repo: repo}
}

func (s *attendanceService) CheckIn(ctx context.Context, orgID, empID uuid.UUID, method string, lat, lng *float64) (*models.Attendance, error) {
now := time.Now()
today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
status := "present"

att := &models.Attendance{
OrgID:       orgID,
EmpID:       empID,
Date:        today,
CheckIn:     &now,
Method:      &method,
Status:      &status,
LocationLat: lat,
LocationLng: lng,
}

if err := s.repo.CheckIn(ctx, att); err != nil {
return nil, err
}
return att, nil
}

func (s *attendanceService) CheckOut(ctx context.Context, empID uuid.UUID) error {
now := time.Now()
today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
return s.repo.CheckOut(ctx, empID, today, now)
}

func (s *attendanceService) GetToday(ctx context.Context, orgIDStr, dateStr string) ([]*models.Attendance, error) {
orgID, err := uuid.Parse(orgIDStr)
if err != nil {
return nil, err
}
targetDate := time.Now()
if dateStr != "" {
targetDate, _ = time.Parse("2006-01-02", dateStr)
}
today := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.Local)
return s.repo.ListByOrgDate(ctx, orgID, today)
}
