package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/hrms-backend/internal/models"
)

type AttendanceRepository interface {
	CheckIn(ctx context.Context, att *models.Attendance) error
	CheckOut(ctx context.Context, empID uuid.UUID, date time.Time, checkOut time.Time) error
	GetTodayRecord(ctx context.Context, empID uuid.UUID, date time.Time) (*models.Attendance, error)
	ListByOrgDate(ctx context.Context, orgID uuid.UUID, date time.Time) ([]*models.Attendance, error)
}

type attendanceRepo struct {
	db *pgxpool.Pool
}

func NewAttendanceRepository(db *pgxpool.Pool) AttendanceRepository {
	return &attendanceRepo{db: db}
}

func (r *attendanceRepo) CheckIn(ctx context.Context, a *models.Attendance) error {
	query := `INSERT INTO attendance_records (org_id, emp_id, date, check_in, method, status, location_lat, location_lng, face_conf_score)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (emp_id, date) DO UPDATE SET check_in = EXCLUDED.check_in, method = EXCLUDED.method, status = EXCLUDED.status
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		a.OrgID, a.EmpID, a.Date, a.CheckIn, a.Method, a.Status, a.LocationLat, a.LocationLng, a.FaceConfScore,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *attendanceRepo) CheckOut(ctx context.Context, empID uuid.UUID, date time.Time, checkOut time.Time) error {
	query := `UPDATE attendance_records
		SET check_out = $1, duration_mins = EXTRACT(EPOCH FROM ($1 - check_in))/60, updated_at = NOW()
		WHERE emp_id = $2 AND date = $3`

	_, err := r.db.Exec(ctx, query, checkOut, empID, date)
	return err
}

func (r *attendanceRepo) GetTodayRecord(ctx context.Context, empID uuid.UUID, date time.Time) (*models.Attendance, error) {
	query := `SELECT id, org_id, emp_id, date, check_in, check_out, method, status FROM attendance_records WHERE emp_id = $1 AND date = $2`
	a := &models.Attendance{}
	err := r.db.QueryRow(ctx, query, empID, date).Scan(
		&a.ID, &a.OrgID, &a.EmpID, &a.Date, &a.CheckIn, &a.CheckOut, &a.Method, &a.Status,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return a, err
}

func (r *attendanceRepo) ListByOrgDate(ctx context.Context, orgID uuid.UUID, date time.Time) ([]*models.Attendance, error) {
	query := `SELECT id, org_id, emp_id, date, check_in, check_out, method, status FROM attendance_records WHERE org_id = $1 AND date = $2`
	rows, err := r.db.Query(ctx, query, orgID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.Attendance
	for rows.Next() {
		a := &models.Attendance{}
		if err := rows.Scan(&a.ID, &a.OrgID, &a.EmpID, &a.Date, &a.CheckIn, &a.CheckOut, &a.Method, &a.Status); err != nil {
			return nil, err
		}
		records = append(records, a)
	}
	return records, nil
}
