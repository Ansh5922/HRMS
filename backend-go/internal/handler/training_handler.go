package handler

import (
"github.com/your-org/hrms-backend/internal/service"
)

type TrainingHandler struct{ service service.TrainingService }

func NewTrainingHandler(svc service.TrainingService) *TrainingHandler { return &TrainingHandler{service: svc} }