package repository

import (
	"context"
	"eventhandler/internal/provider"
	"eventhandler/model"

	"github.com/jackc/pgx/v5"
)

type EventRepository interface {
	SaveInbound(ctx context.Context, event *model.QueueEvent) error
	SaveEvent(ctx context.Context, event *model.QueueEvent) error
}

type repo struct {
	logger provider.ILogger
	conn   *pgx.Conn
}

func NewEventRepository(logger provider.ILogger, conn *pgx.Conn) EventRepository {
	return &repo{
		logger: logger,
		conn:   conn,
	}
}
