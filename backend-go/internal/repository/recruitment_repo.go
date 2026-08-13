package repository

import (
"github.com/jackc/pgx/v5/pgxpool"
)

type RecruitmentRepository interface{}
type recruitmentRepo struct{ db *pgxpool.Pool }

func NewRecruitmentRepository(db *pgxpool.Pool) RecruitmentRepository { return &recruitmentRepo{db: db} }