package repository

import (
	"context"
	"eventhandler/internal/provider"
	"eventhandler/model/entity"
	"eventhandler/util"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SaveEvent saves an inbound event to the database
type InboundOutboundRepository interface {
	SaveEvent(ctx context.Context, req *entity.CreateAccountEventRequest) error
	SaveMessageInbound(ctx context.Context, req *entity.CreateMessageInboundRequest) error
	SaveMessageOutbound(ctx context.Context, req *entity.CreateMessageOutboundRequest) error
	SaveMessageReceipt(ctx context.Context, req *entity.CreateMessageReceiptRequest) error
}

type inboundOutboundRepo struct {
	logger provider.ILogger
	pool   *pgxpool.Pool
}

func NewInboundOutboundRepository(logger provider.ILogger, pool *pgxpool.Pool) InboundOutboundRepository {
	return &inboundOutboundRepo{
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
			(inbound_id, account_id, from_me, message_id, sender, message_type, received_at, data, created_at, updated_at) 
		VALUES (@inbound_id, @account_id, @from_me, @message_id, @sender, @message_type, @received_at, @data, @created_at, @updated_at)`

	insertMessageOutboundQuery = `
		INSERT INTO whatsapp_web.message_outbounds
			(outbound_id, account_id, message_id, recipient, message_type, sent_at, data, created_at, updated_at)
		VALUES (@outbound_id, @account_id, @message_id, @recipient, @message_type, @sent_at, @data, @created_at, @updated_at)`

	insertMessageReceiptQuery = `
		INSERT INTO whatsapp_web.message_receipts
			(receipt_id, account_id, message_id, timestamp, status, created_at, updated_at)
		VALUES (@receipt_id, @account_id, @message_id, @timestamp, @status, @created_at, @updated_at)`
)

func (r *inboundOutboundRepo) SaveEvent(ctx context.Context, req *entity.CreateAccountEventRequest) error {
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

func (r *inboundOutboundRepo) SaveMessageInbound(ctx context.Context, req *entity.CreateMessageInboundRequest) error {
	params := pgx.NamedArgs{
		"inbound_id":   req.EventID,
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

func (r *inboundOutboundRepo) SaveMessageOutbound(ctx context.Context, req *entity.CreateMessageOutboundRequest) error {
	params := pgx.NamedArgs{
		"outbound_id":  req.EventID,
		"account_id":   req.AccountID,
		"message_id":   req.MessageID,
		"recipient":    req.Recipient,
		"message_type": req.MessageType,
		"sent_at":      req.SentAt,
		"data":         req.Data,
		"created_at":   util.NowUTC(),
		"updated_at":   util.NowUTC(),
	}

	_, err := r.pool.Exec(ctx, insertMessageOutboundQuery, params)
	return err
}

func (r *inboundOutboundRepo) SaveMessageReceipt(ctx context.Context, req *entity.CreateMessageReceiptRequest) error {
	params := pgx.NamedArgs{
		"receipt_id": req.ReceiptID,
		"account_id": req.AccountID,
		"message_id": req.MessageID,
		"timestamp":  req.Timestamp,
		"status":     req.Status,
		"created_at": util.NowUTC(),
		"updated_at": util.NowUTC(),
	}

	_, err := r.pool.Exec(ctx, insertMessageReceiptQuery, params)
	return err
}
