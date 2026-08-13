package repository

import (
"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationRepository interface{}
type organizationRepo struct{ db *pgxpool.Pool }

func NewOrganizationRepository(db *pgxpool.Pool) OrganizationRepository { return &organizationRepo{db: db} }