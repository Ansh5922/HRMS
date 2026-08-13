package handler

import (
"github.com/your-org/hrms-backend/internal/service"
)

type PerformanceHandler struct{ service service.PerformanceService }

func NewPerformanceHandler(svc service.PerformanceService) *PerformanceHandler { return &PerformanceHandler{service: svc} }