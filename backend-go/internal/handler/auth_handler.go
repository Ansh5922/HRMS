package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/models"
	"github.com/your-org/hrms-backend/internal/service"
	"github.com/your-org/hrms-backend/pkg/response"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{service: svc}
}

// ── Request DTOs ─────────────────────────────────────────────────────────────

type RegisterReq struct {
	OrgName  string `json:"org_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ForgotPasswordReq struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordReq struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type CreateUserReq struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=8"`
	RoleID   *string `json:"role_id"`
}

type UpdateUserReq struct {
	Email      *string `json:"email"`
	IsActive   *bool   `json:"is_active"`
	IsVerified *bool   `json:"is_verified"`
	RoleID     *string `json:"role_id"`
}

type CreateRoleReq struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type UpdateRoleReq struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type AssignRoleReq struct {
	UserID string `json:"user_id" binding:"required"`
	RoleID string `json:"role_id" binding:"required"`
}

type SetPermissionsReq struct {
	PermissionIDs []string `json:"permission_ids" binding:"required"`
}

// ── Public Auth Endpoints ────────────────────────────────────────────────────

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	res, err := h.service.RegisterOrgAndAdmin(c.Request.Context(), req.OrgName, req.Email, req.Password)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	res, err := h.service.Login(c.Request.Context(), req.Email, req.Password, ipAddress, userAgent)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	response.OK(c, res, "login successful")
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	res, err := h.service.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	response.OK(c, res, "token refreshed")
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	res, err := h.service.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, res)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, nil, "password reset successful")
}

// ── Protected Profile Endpoints ──────────────────────────────────────────────

func (h *AuthHandler) GetMe(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	res, err := h.service.GetMe(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "user")
		return
	}

	response.OK(c, res)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	var req ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, nil, "password changed successfully")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "logged out successfully")
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	if err := h.service.LogoutAll(c.Request.Context(), userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "logged out from all devices")
}

// ── User Management (Admin) ─────────────────────────────────────────────────

func (h *AuthHandler) CreateUser(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var roleID *uuid.UUID
	if req.RoleID != nil {
		parsed, err := uuid.Parse(*req.RoleID)
		if err != nil {
			response.BadRequest(c, "invalid role_id")
			return
		}
		roleID = &parsed
	}

	user, err := h.service.CreateUser(c.Request.Context(), orgID, req.Email, req.Password, roleID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, user)
}

func (h *AuthHandler) ListUsers(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	users, err := h.service.ListUsers(c.Request.Context(), orgID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, users)
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user := &models.User{ID: id}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.IsVerified != nil {
		user.IsVerified = *req.IsVerified
	}
	if req.RoleID != nil {
		parsed, _ := uuid.Parse(*req.RoleID)
		user.RoleID = &parsed
	}

	if err := h.service.UpdateUser(c.Request.Context(), user); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "user updated")
}

func (h *AuthHandler) DeactivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	if err := h.service.DeactivateUser(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "user deactivated")
}

// ── Role Management ──────────────────────────────────────────────────────────

func (h *AuthHandler) CreateRole(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	var req CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	role := &models.Role{
		OrgID:       orgID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.CreateRole(c.Request.Context(), role); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, role)
}

func (h *AuthHandler) ListRoles(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	roles, err := h.service.ListRoles(c.Request.Context(), orgID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, roles)
}

func (h *AuthHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid role ID")
		return
	}

	role, err := h.service.GetRole(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "role")
		return
	}

	response.OK(c, role)
}

func (h *AuthHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid role ID")
		return
	}

	var req UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	role := &models.Role{ID: id}
	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = req.Description
	}

	if err := h.service.UpdateRole(c.Request.Context(), role); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "role updated")
}

func (h *AuthHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid role ID")
		return
	}

	if err := h.service.DeleteRole(c.Request.Context(), id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "role deleted")
}

func (h *AuthHandler) AssignRole(c *gin.Context) {
	var req AssignRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		response.BadRequest(c, "invalid user_id")
		return
	}
	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		response.BadRequest(c, "invalid role_id")
		return
	}

	if err := h.service.AssignRole(c.Request.Context(), userID, roleID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "role assigned to user")
}

// ── Permission Management ────────────────────────────────────────────────────

func (h *AuthHandler) ListPermissions(c *gin.Context) {
	perms, err := h.service.ListPermissions(c.Request.Context())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, perms)
}

func (h *AuthHandler) GetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid role ID")
		return
	}

	perms, err := h.service.GetRolePermissions(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, perms)
}

func (h *AuthHandler) SetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid role ID")
		return
	}

	var req SetPermissionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	permIDs := make([]uuid.UUID, 0, len(req.PermissionIDs))
	for _, pidStr := range req.PermissionIDs {
		pid, err := uuid.Parse(pidStr)
		if err != nil {
			response.BadRequest(c, "invalid permission ID: "+pidStr)
			return
		}
		permIDs = append(permIDs, pid)
	}

	if err := h.service.SetRolePermissions(c.Request.Context(), roleID, permIDs); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "role permissions updated")
}