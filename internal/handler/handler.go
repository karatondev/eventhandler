package handler

import (
	"context"
	"encoding/json"
	"eventhandler/internal/provider"
	"eventhandler/internal/service"
	"eventhandler/model"
	"eventhandler/model/constant"
	"fmt"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerHandler interface {
	Handle(data amqp.Delivery)
}

type consumerHandler struct {
	logger  provider.ILogger
	service service.Events
}

func NewConsumerHandler(logger provider.ILogger, service service.Events) ConsumerHandler {
	return &consumerHandler{
		logger:  logger,
		service: service,
	}
}

func (c *consumerHandler) Handle(data amqp.Delivery) {
	ctx := context.WithValue(context.Background(), constant.CtxReqIDKey, fmt.Sprintf("%s", uuid.New().String()))
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			c.logger.Errorfctx(provider.AppLog, ctx, true, "Unhandled panic: %v", r)
			data.Nack(false, true)
		}
	}(ctx)

	c.logger.Infofctx(provider.AppLog, ctx, "Received message: %s", string(data.Body))
	payload := model.QueueEvent{}
	if err := json.Unmarshal(data.Body, &payload); err != nil {
		c.logger.Errorfctx(provider.AppLog, ctx, false, "Failed to unmarshal message: %v", err)
		return
	}

	if err := c.service.HandleEvent(ctx, &payload); err != nil {
		c.logger.Errorfctx(provider.AppLog, ctx, false, "Failed to process queue: %v", err)
	}

	if err := data.Ack(false); err != nil {
		c.logger.Errorfctx(provider.AppLog, ctx, false, "Failed to acknowledge message: %v", err)
		return
	}

	c.logger.Infofctx(provider.AppLog, ctx, "Message acknowledged")
}
