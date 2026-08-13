package handler

import (
"github.com/your-org/hrms-backend/internal/service"
)

type OrganizationHandler struct{ service service.OrganizationService }

func NewOrganizationHandler(svc service.OrganizationService) *OrganizationHandler { return &OrganizationHandler{service: svc} }