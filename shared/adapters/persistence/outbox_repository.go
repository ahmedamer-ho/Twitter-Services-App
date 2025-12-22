package persistence

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OutboxRepository struct {
	db *pgx.Conn
}

func NewOutboxRepository(db *pgx.Conn) *OutboxRepository {
	return &OutboxRepository{db: db}
}
