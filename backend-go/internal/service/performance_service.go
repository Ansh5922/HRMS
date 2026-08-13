package service

import (
"github.com/your-org/hrms-backend/internal/repository"
)

type PerformanceService interface{}
type performanceService struct{ repo repository.PerformanceRepository }

func NewPerformanceService(repo repository.PerformanceRepository) PerformanceService { return &performanceService{repo: repo} }