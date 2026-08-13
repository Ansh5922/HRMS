package service

import (
"github.com/your-org/hrms-backend/internal/repository"
)

type OrganizationService interface{}
type organizationService struct{ repo repository.OrganizationRepository }

func NewOrganizationService(repo repository.OrganizationRepository) OrganizationService { return &organizationService{repo: repo} }