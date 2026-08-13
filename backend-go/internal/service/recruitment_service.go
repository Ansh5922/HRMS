package service

import (
"github.com/your-org/hrms-backend/internal/repository"
)

type RecruitmentService interface{}
type recruitmentService struct{ repo repository.RecruitmentRepository }

func NewRecruitmentService(repo repository.RecruitmentRepository) RecruitmentService { return &recruitmentService{repo: repo} }