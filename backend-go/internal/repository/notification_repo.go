package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/hrms-backend/internal/models"
)

type NotificationRepository interface {
	// Notifications
	CreateNotification(ctx context.Context, n *models.Notification) error
	ListNotificationsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	MarkAsRead(ctx context.Context, id, userID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	DeleteNotification(ctx context.Context, id, userID uuid.UUID) error

	// Announcements
	CreateAnnouncement(ctx context.Context, a *models.Announcement) error
	ListAnnouncementsByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.Announcement, error)
	GetAnnouncementByID(ctx context.Context, id uuid.UUID) (*models.Announcement, error)
	DeleteAnnouncement(ctx context.Context, id uuid.UUID) error

	// Preferences
	SavePreference(ctx context.Context, pref *models.NotificationPreference) error
	GetPreferencesByUser(ctx context.Context, userID uuid.UUID) ([]*models.NotificationPreference, error)

	// Push Tokens
	RegisterPushToken(ctx context.Context, pt *models.PushToken) error
	ListPushTokensByUser(ctx context.Context, userID uuid.UUID) ([]*models.PushToken, error)
	RemovePushToken(ctx context.Context, userID uuid.UUID, token string) error
}

type notificationRepo struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) NotificationRepository {
	return &notificationRepo{db: db}
}

// ── Notifications ────────────────────────────────────────────────────────────

func (r *notificationRepo) CreateNotification(ctx context.Context, n *models.Notification) error {
	query := `INSERT INTO notifications (org_id, user_id, type, title, message, action_url)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, n.OrgID, n.UserID, n.Type, n.Title, n.Message, n.ActionURL).
		Scan(&n.ID, &n.CreatedAt)
}

func (r *notificationRepo) ListNotificationsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `SELECT id, org_id, user_id, type, title, message, action_url, is_read, read_at, created_at
		FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Notification
	for rows.Next() {
		n := &models.Notification{}
		if err := rows.Scan(&n.ID, &n.OrgID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.ActionURL, &n.IsRead, &n.ReadAt, &n.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, nil
}

func (r *notificationRepo) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *notificationRepo) MarkAsRead(ctx context.Context, id, userID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true, read_at = NOW() WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, id, userID)
	return err
}

func (r *notificationRepo) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true, read_at = NOW() WHERE user_id = $1 AND is_read = false`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *notificationRepo) DeleteNotification(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, id, userID)
	return err
}

// ── Announcements ────────────────────────────────────────────────────────────

func (r *notificationRepo) CreateAnnouncement(ctx context.Context, a *models.Announcement) error {
	query := `INSERT INTO announcements (org_id, title, content, author_id, target_dept, published_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, COALESCE($6, NOW()), $7) RETURNING id, published_at`
	return r.db.QueryRow(ctx, query, a.OrgID, a.Title, a.Content, a.AuthorID, a.TargetDept, a.PublishedAt, a.ExpiresAt).
		Scan(&a.ID, &a.PublishedAt)
}

func (r *notificationRepo) ListAnnouncementsByOrg(ctx context.Context, orgID uuid.UUID) ([]*models.Announcement, error) {
	query := `SELECT a.id, a.org_id, a.title, a.content, a.author_id, a.target_dept, a.published_at, a.expires_at,
		u.email as author_name, d.name as target_dept_name
	FROM announcements a
	LEFT JOIN users u ON u.id = a.author_id
	LEFT JOIN departments d ON d.id = a.target_dept
	WHERE a.org_id = $1 AND (a.expires_at IS NULL OR a.expires_at > NOW())
	ORDER BY a.published_at DESC`

	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Announcement
	for rows.Next() {
		a := &models.Announcement{}
		if err := rows.Scan(&a.ID, &a.OrgID, &a.Title, &a.Content, &a.AuthorID, &a.TargetDept, &a.PublishedAt, &a.ExpiresAt, &a.AuthorName, &a.TargetDeptName); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *notificationRepo) GetAnnouncementByID(ctx context.Context, id uuid.UUID) (*models.Announcement, error) {
	query := `SELECT a.id, a.org_id, a.title, a.content, a.author_id, a.target_dept, a.published_at, a.expires_at,
		u.email as author_name, d.name as target_dept_name
	FROM announcements a
	LEFT JOIN users u ON u.id = a.author_id
	LEFT JOIN departments d ON d.id = a.target_dept
	WHERE a.id = $1`

	a := &models.Announcement{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.OrgID, &a.Title, &a.Content, &a.AuthorID, &a.TargetDept, &a.PublishedAt, &a.ExpiresAt, &a.AuthorName, &a.TargetDeptName,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return a, err
}

func (r *notificationRepo) DeleteAnnouncement(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM announcements WHERE id = $1`, id)
	return err
}

// ── Preferences ──────────────────────────────────────────────────────────────

func (r *notificationRepo) SavePreference(ctx context.Context, pref *models.NotificationPreference) error {
	query := `INSERT INTO notification_preferences (user_id, event_type, email, push, sms)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, event_type) DO UPDATE SET email=$3, push=$4, sms=$5
		RETURNING id`
	return r.db.QueryRow(ctx, query, pref.UserID, pref.EventType, pref.Email, pref.Push, pref.SMS).Scan(&pref.ID)
}

func (r *notificationRepo) GetPreferencesByUser(ctx context.Context, userID uuid.UUID) ([]*models.NotificationPreference, error) {
	query := `SELECT id, user_id, event_type, email, push, sms FROM notification_preferences WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prefs []*models.NotificationPreference
	for rows.Next() {
		p := &models.NotificationPreference{}
		if err := rows.Scan(&p.ID, &p.UserID, &p.EventType, &p.Email, &p.Push, &p.SMS); err != nil {
			return nil, err
		}
		prefs = append(prefs, p)
	}
	return prefs, nil
}

// ── Push Tokens ──────────────────────────────────────────────────────────────

func (r *notificationRepo) RegisterPushToken(ctx context.Context, pt *models.PushToken) error {
	query := `INSERT INTO push_tokens (user_id, token, platform)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, token) DO UPDATE SET platform=$3, created_at=NOW()
		RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, pt.UserID, pt.Token, pt.Platform).Scan(&pt.ID, &pt.CreatedAt)
}

func (r *notificationRepo) ListPushTokensByUser(ctx context.Context, userID uuid.UUID) ([]*models.PushToken, error) {
	query := `SELECT id, user_id, token, platform, created_at FROM push_tokens WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*models.PushToken
	for rows.Next() {
		pt := &models.PushToken{}
		if err := rows.Scan(&pt.ID, &pt.UserID, &pt.Token, &pt.Platform, &pt.CreatedAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, pt)
	}
	return tokens, nil
}

func (r *notificationRepo) RemovePushToken(ctx context.Context, userID uuid.UUID, token string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM push_tokens WHERE user_id = $1 AND token = $2`, userID, token)
	return err
}