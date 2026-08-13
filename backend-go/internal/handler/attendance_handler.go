package handler

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/hrms-backend/internal/service"
	"github.com/your-org/hrms-backend/pkg/aiservice"
	"github.com/your-org/hrms-backend/pkg/response"
)

type AttendanceHandler struct {
	service  service.AttendanceService
	aiClient *aiservice.AIServiceClient
}

func NewAttendanceHandler(svc service.AttendanceService, aiClient *aiservice.AIServiceClient) *AttendanceHandler {
	return &AttendanceHandler{service: svc, aiClient: aiClient}
}

type CheckInReq struct {
	EmpID  string   `json:"emp_id" binding:"required"`
	Method string   `json:"method"` // 'manual', 'geo', 'face'
	Lat    *float64 `json:"lat"`
	Lng    *float64 `json:"lng"`
}

func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, _ := uuid.Parse(orgIDStr)

	var req CheckInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	empID, err := uuid.Parse(req.EmpID)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}

	method := req.Method
	if method == "" {
		method = "manual"
	}

	att, err := h.service.CheckIn(c.Request.Context(), orgID, empID, method, req.Lat, req.Lng)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, att, "check-in successful")
}

func (h *AttendanceHandler) FaceCheckIn(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	orgID, _ := uuid.Parse(orgIDStr)

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		response.BadRequest(c, "photo file is required")
		return
	}
	defer file.Close()

	imgBytes, err := io.ReadAll(file)
	if err != nil {
		response.InternalError(c, "failed to read photo")
		return
	}

	// Call FastAPI AI service to verify face
	result, err := h.aiClient.VerifyFace(imgBytes, header.Filename)
	if err != nil {
		response.InternalError(c, "face verification failed: "+err.Error())
		return
	}

	if !result.Matched || result.EmpID == nil {
		response.BadRequest(c, "face not recognized")
		return
	}

	empID, err := uuid.Parse(*result.EmpID)
	if err != nil {
		response.BadRequest(c, "invalid employee ID returned by AI service")
		return
	}

	method := "face"
	att, err := h.service.CheckIn(c.Request.Context(), orgID, empID, method, nil, nil)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, att, "face check-in successful")
}

func (h *AttendanceHandler) CheckOut(c *gin.Context) {
	empIDStr := c.Param("emp_id")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		response.BadRequest(c, "invalid emp_id")
		return
	}

	if err := h.service.CheckOut(c.Request.Context(), empID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil, "check-out successful")
}

func (h *AttendanceHandler) GetToday(c *gin.Context) {
	orgIDStr := c.GetString("org_id")
	dateStr := c.Query("date")

	records, err := h.service.GetToday(c.Request.Context(), orgIDStr, dateStr)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, records)
}