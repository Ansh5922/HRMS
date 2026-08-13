package service

import (
"github.com/your-org/hrms-backend/internal/repository"
)

type TrainingService interface{}
type trainingService struct{ repo repository.TrainingRepository }

func NewTrainingService(repo repository.TrainingRepository) TrainingService { return &trainingService{repo: repo} }