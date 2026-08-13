package repository

import (
"github.com/jackc/pgx/v5/pgxpool"
)

type PerformanceRepository interface{}
type performanceRepo struct{ db *pgxpool.Pool }

func NewPerformanceRepository(db *pgxpool.Pool) PerformanceRepository { return &performanceRepo{db: db} }