package repository

import (
	"context"
	"eventhandler/entity"
	"eventhandler/internal/provider"
	"eventhandler/util"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SaveEvent saves an inbound event to the database
type EventInboundRepository interface {
	SaveEvent(ctx context.Context, req *entity.CreateAccountEventRequest) error
	SaveMessageInbound(ctx context.Context, req *entity.CreateMessageInboundRequest) error
}

type eventInboundRepo struct {
	logger provider.ILogger
	pool   *pgxpool.Pool
}

func NewEventInboundRepository(logger provider.ILogger, pool *pgxpool.Pool) EventInboundRepository {
	return &eventInboundRepo{
		logger: logger,
		pool:   pool,
	}
}

const (
	insertAccountEventQuery = `
		INSERT INTO whatsapp_web.account_events 
			(event_id, account_id, event_type, timestamp, data, created_at, updated_at) 
		VALUES (@event_id, @account_id, @event_type, @timestamp, @data, @created_at, @updated_at)`

	insertMessageInboundQuery = `
		INSERT INTO whatsapp_web.message_inbounds 
			(account_id, from_me, message_id, sender, message_type, received_at, data, created_at, updated_at) 
		VALUES (@account_id, @from_me, @message_id, @sender, @message_type, @received_at, @data, @created_at, @updated_at)`
)

func (r *eventInboundRepo) SaveEvent(ctx context.Context, req *entity.CreateAccountEventRequest) error {
	params := pgx.NamedArgs{
		"event_id":   req.EventID,
		"account_id": req.AccountID,
		"event_type": req.EventType,
		"timestamp":  req.Timestamp,
		"data":       req.Data,
		"created_at": util.NowUTC(),
		"updated_at": util.NowUTC(),
	}

	_, err := r.pool.Exec(ctx, insertAccountEventQuery, params)
	return err
}

func (r *eventInboundRepo) SaveMessageInbound(ctx context.Context, req *entity.CreateMessageInboundRequest) error {
	params := pgx.NamedArgs{
		"event_id":     req.EventID,
		"account_id":   req.AccountID,
		"from_me":      req.FromMe,
		"message_id":   req.MessageID,
		"sender":       req.Sender,
		"message_type": req.MessageType,
		"received_at":  req.ReceivedAt,
		"data":         req.Data,
		"created_at":   util.NowUTC(),
		"updated_at":   util.NowUTC(),
	}

	_, err := r.pool.Exec(ctx, insertMessageInboundQuery, params)
	return err
}
