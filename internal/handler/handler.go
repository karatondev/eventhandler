package handler

import (
	"context"
	"encoding/json"
	"eventhandler/internal/service"
	"fmt"

	"zaplio/shared/constant"
	sharedmodel "zaplio/shared/model"
	"zaplio/shared/pkg/logger"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerHandler interface {
	Handle(data amqp.Delivery)
}

type consumerHandler struct {
	logger  logger.ILogger
	service service.Events
}

func NewConsumerHandler(log logger.ILogger, service service.Events) ConsumerHandler {
	return &consumerHandler{
		logger:  log,
		service: service,
	}
}

func (c *consumerHandler) Handle(data amqp.Delivery) {
	ctx := context.WithValue(context.Background(), constant.CtxReqIDKey, fmt.Sprintf("%s", uuid.New().String()))
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			c.logger.Errorfctx(logger.AppLog, ctx, true, "Unhandled panic: %v", r)
			data.Nack(false, true)
		}
	}(ctx)

	c.logger.Infofctx(logger.AppLog, ctx, "Received message: %s", string(data.Body))
	payload := sharedmodel.QueueEvent{}
	if err := json.Unmarshal(data.Body, &payload); err != nil {
		c.logger.Errorfctx(logger.AppLog, ctx, false, "Failed to unmarshal message: %v", err)
		data.Nack(false, false)
		return
	}
	c.logger.Infofctx(logger.AppLog, ctx, "Queue message parsed: event_type=%s, event_id=%s, sender_jid=%s", payload.EventType, payload.EventID, payload.SenderJID)

	if err := c.service.HandleEvent(ctx, &payload); err != nil {
		c.logger.Errorfctx(logger.AppLog, ctx, false, "Failed to process queue: %v", err)
		data.Nack(false, false)
		return
	}

	if err := data.Ack(false); err != nil {
		c.logger.Errorfctx(logger.AppLog, ctx, false, "Failed to acknowledge message: %v", err)
		return
	}

	c.logger.Infofctx(logger.AppLog, ctx, "Message acknowledged")
}
