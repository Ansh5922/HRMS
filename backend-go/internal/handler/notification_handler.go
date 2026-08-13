package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/service"
	"github.com/your-org/hrms-backend/pkg/response"
)

type NotificationHandler struct {
	service service.NotificationService
}

func NewNotificationHandler(svc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: svc}
}

// ── Request DTOs ─────────────────────────────────────────────────────────────

type CreateAnnouncementReq struct {
	Title      string  `json:"title" binding:"required"`
	Content    string  `json:"content" binding:"required"`
	TargetDept *string `json:"target_dept"`
	ExpiresAt  *string `json:"expires_at"` // RFC3339 or YYYY-MM-DD
}

type SavePreferenceReq struct {
	EventType string `json:"event_type" binding:"required"`
	Email     bool   `json:"email"`
	Push      bool   `json:"push"`
	SMS       bool   `json:"sms"`
}

type RegisterPushTokenReq struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform"` // android, ios, web
}

// ── Notifications Endpoints ──────────────────────────────────────────────────

func (h *NotificationHandler) ListMyNotifications(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	list, err := h.service.ListMyNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, list)
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	count, err := h.service.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"unread_count": count})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid notification ID")
		return
	}

	if err := h.service.MarkAsRead(c.Request.Context(), id, userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "notification marked as read")
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	if err := h.service.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "all notifications marked as read")
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid notification ID")
		return
	}

	if err := h.service.DeleteNotification(c.Request.Context(), id, userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "notification deleted")
}

// ── Announcements Endpoints ──────────────────────────────────────────────────

func (h *NotificationHandler) CreateAnnouncement(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	var req CreateAnnouncementReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	a := &models.Announcement{
		OrgID:    orgID,
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: &userID,
	}

	if req.TargetDept != nil {
		if dID, err := uuid.Parse(*req.TargetDept); err == nil {
			a.TargetDept = &dID
		}
	}
	if req.ExpiresAt != nil {
		if exp, err := time.Parse(time.RFC3339, *req.ExpiresAt); err == nil {
			a.ExpiresAt = &exp
		}
	}

	res, err := h.service.CreateAnnouncement(c.Request.Context(), a)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *NotificationHandler) ListAnnouncements(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	list, err := h.service.ListAnnouncements(c.Request.Context(), orgID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, list)
}

func (h *NotificationHandler) GetAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid announcement ID")
		return
	}

	a, err := h.service.GetAnnouncementByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "announcement")
		return
	}

	response.OK(c, a)
}

func (h *NotificationHandler) DeleteAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid announcement ID")
		return
	}

	if err := h.service.DeleteAnnouncement(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "announcement deleted")
}

// ── Preferences Endpoints ────────────────────────────────────────────────────

func (h *NotificationHandler) SavePreference(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	var req SavePreferenceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pref := &models.NotificationPreference{
		UserID:    userID,
		EventType: req.EventType,
		Email:     req.Email,
		Push:      req.Push,
		SMS:       req.SMS,
	}

	if err := h.service.SavePreference(c.Request.Context(), pref); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, pref, "notification preference saved")
}

func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	prefs, err := h.service.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, prefs)
}

// ── Push Token Endpoints ─────────────────────────────────────────────────────

func (h *NotificationHandler) RegisterPushToken(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	var req RegisterPushTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pt := &models.PushToken{
		UserID:   userID,
		Token:    req.Token,
		Platform: req.Platform,
	}

	if err := h.service.RegisterPushToken(c.Request.Context(), pt); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, pt, "push token registered")
}

func (h *NotificationHandler) RemovePushToken(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		response.BadRequest(c, "token query param is required")
		return
	}

	if err := h.service.RemovePushToken(c.Request.Context(), userID, token); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "push token removed")
}