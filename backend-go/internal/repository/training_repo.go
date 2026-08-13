package repository

import (
"github.com/jackc/pgx/v5/pgxpool"
)

type TrainingRepository interface{}
type trainingRepo struct{ db *pgxpool.Pool }

func NewTrainingRepository(db *pgxpool.Pool) TrainingRepository { return &trainingRepo{db: db} }