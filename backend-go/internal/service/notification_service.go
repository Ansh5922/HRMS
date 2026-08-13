package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/repository"
)

type NotificationService interface {
	// Notifications
	SendNotification(ctx context.Context, n *models.Notification) (*models.Notification, error)
	ListMyNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	MarkAsRead(ctx context.Context, id, userID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	DeleteNotification(ctx context.Context, id, userID uuid.UUID) error

	// Announcements
	CreateAnnouncement(ctx context.Context, a *models.Announcement) (*models.Announcement, error)
	ListAnnouncements(ctx context.Context, orgID uuid.UUID) ([]*models.Announcement, error)
	GetAnnouncementByID(ctx context.Context, id uuid.UUID) (*models.Announcement, error)
	DeleteAnnouncement(ctx context.Context, id uuid.UUID) error

	// Preferences
	SavePreference(ctx context.Context, pref *models.NotificationPreference) error
	GetPreferences(ctx context.Context, userID uuid.UUID) ([]*models.NotificationPreference, error)

	// Push Tokens
	RegisterPushToken(ctx context.Context, pt *models.PushToken) error
	RemovePushToken(ctx context.Context, userID uuid.UUID, token string) error
}

type notificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

// ── Notifications ────────────────────────────────────────────────────────────

func (s *notificationService) SendNotification(ctx context.Context, n *models.Notification) (*models.Notification, error) {
	if strings.TrimSpace(n.Title) == "" {
		return nil, errors.New("notification title is required")
	}
	if n.Type == "" {
		n.Type = "general"
	}
	if err := s.repo.CreateNotification(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *notificationService) ListMyNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	return s.repo.ListNotificationsByUser(ctx, userID, limit, offset)
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

func (s *notificationService) MarkAsRead(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.MarkAsRead(ctx, id, userID)
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *notificationService) DeleteNotification(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.DeleteNotification(ctx, id, userID)
}

// ── Announcements ────────────────────────────────────────────────────────────

func (s *notificationService) CreateAnnouncement(ctx context.Context, a *models.Announcement) (*models.Announcement, error) {
	if strings.TrimSpace(a.Title) == "" || strings.TrimSpace(a.Content) == "" {
		return nil, errors.New("announcement title and content are required")
	}
	if err := s.repo.CreateAnnouncement(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *notificationService) ListAnnouncements(ctx context.Context, orgID uuid.UUID) ([]*models.Announcement, error) {
	return s.repo.ListAnnouncementsByOrg(ctx, orgID)
}

func (s *notificationService) GetAnnouncementByID(ctx context.Context, id uuid.UUID) (*models.Announcement, error) {
	a, err := s.repo.GetAnnouncementByID(ctx, id)
	if err != nil || a == nil {
		return nil, errors.New("announcement not found")
	}
	return a, nil
}

func (s *notificationService) DeleteAnnouncement(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteAnnouncement(ctx, id)
}

// ── Preferences ──────────────────────────────────────────────────────────────

func (s *notificationService) SavePreference(ctx context.Context, pref *models.NotificationPreference) error {
	if pref.EventType == "" {
		return errors.New("event_type is required")
	}
	return s.repo.SavePreference(ctx, pref)
}

func (s *notificationService) GetPreferences(ctx context.Context, userID uuid.UUID) ([]*models.NotificationPreference, error) {
	return s.repo.GetPreferencesByUser(ctx, userID)
}

// ── Push Tokens ──────────────────────────────────────────────────────────────

func (s *notificationService) RegisterPushToken(ctx context.Context, pt *models.PushToken) error {
	if pt.Token == "" {
		return errors.New("push token is required")
	}
	if pt.Platform == "" {
		pt.Platform = "web"
	}
	return s.repo.RegisterPushToken(ctx, pt)
}

func (s *notificationService) RemovePushToken(ctx context.Context, userID uuid.UUID, token string) error {
	return s.repo.RemovePushToken(ctx, userID, token)
}