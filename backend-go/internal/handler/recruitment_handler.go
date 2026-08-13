package handler

import (
"github.com/your-org/hrms-backend/internal/service"
)

type RecruitmentHandler struct{ service service.RecruitmentService }

func NewRecruitmentHandler(svc service.RecruitmentService) *RecruitmentHandler { return &RecruitmentHandler{service: svc} }