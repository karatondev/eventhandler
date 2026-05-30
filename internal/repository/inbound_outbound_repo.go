package repository

import (
	"context"
	"eventhandler/model/entity"
	"eventhandler/util"

	"zaplio/shared/pkg/logger"

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
	logger logger.ILogger
	pool   *pgxpool.Pool
}

func NewInboundOutboundRepository(log logger.ILogger, pool *pgxpool.Pool) InboundOutboundRepository {
	return &inboundOutboundRepo{
		logger: log,
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
		VALUES (@receipt_id, @account_id, @message_id, @timestamp, @status, @created_at, @updated_at)
		ON CONFLICT (receipt_id) DO NOTHING`
)

func (r *inboundOutboundRepo) SaveEvent(ctx context.Context, req *entity.CreateAccountEventRequest) error {
	r.logger.Infofctx(logger.AppLog, ctx, "DB insert account_event: event_id=%s, account_id=%s, event_type=%s, data=%s", req.EventID, req.AccountID, req.EventType, string(req.Data))
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
	if err != nil {
		return err
	}
	r.logger.Infofctx(logger.AppLog, ctx, "DB insert account_event success: event_id=%s", req.EventID)
	return nil
}

func (r *inboundOutboundRepo) SaveMessageInbound(ctx context.Context, req *entity.CreateMessageInboundRequest) error {
	r.logger.Infofctx(logger.AppLog, ctx, "DB insert message_inbound: inbound_id=%s, account_id=%s, message_id=%s, sender=%s, message_type=%s", req.EventID, req.AccountID, req.MessageID, req.Sender, req.MessageType)
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
	if err != nil {
		return err
	}
	r.logger.Infofctx(logger.AppLog, ctx, "DB insert message_inbound success: inbound_id=%s", req.EventID)
	return nil
}

func (r *inboundOutboundRepo) SaveMessageOutbound(ctx context.Context, req *entity.CreateMessageOutboundRequest) error {
	r.logger.Infofctx(logger.AppLog, ctx, "DB insert message_outbound: outbound_id=%s, account_id=%s, message_id=%s, recipient=%s, message_type=%s", req.EventID, req.AccountID, req.MessageID, req.Recipient, req.MessageType)
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
	if err != nil {
		return err
	}
	r.logger.Infofctx(logger.AppLog, ctx, "DB insert message_outbound success: outbound_id=%s", req.EventID)
	return nil
}

func (r *inboundOutboundRepo) SaveMessageReceipt(ctx context.Context, req *entity.CreateMessageReceiptRequest) error {
	r.logger.Infofctx(logger.AppLog, ctx, "DB insert message_receipt: receipt_id=%s, account_id=%s, message_id=%s, status=%s", req.ReceiptID, req.AccountID, req.MessageID, req.Status)
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
	if err != nil {
		return err
	}
	r.logger.Infofctx(logger.AppLog, ctx, "DB insert message_receipt success: receipt_id=%s", req.ReceiptID)
	return nil
}
