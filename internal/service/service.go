package service

import (
	"context"
	"eventhandler/entity"
	"eventhandler/internal/provider"
	"eventhandler/internal/repository"
	"eventhandler/model"

	"github.com/redis/go-redis/v9"
)

type Events interface {
	HandleEvent(ctx context.Context, data *model.QueueEvent) error
	GetAccountBySenderJID(ctx context.Context, senderJID string) (*entity.WhatsAppAccount, error)
}

type service struct {
	logger           provider.ILogger
	redis            *redis.Client
	eventInboundRepo repository.EventInboundRepository
	whatsappRepo     repository.WhatsAppAccountRepository
}

func NewService(logger provider.ILogger, redis *redis.Client, eventInboundRepo repository.EventInboundRepository, whatsappRepo repository.WhatsAppAccountRepository) Events {
	return &service{
		logger:           logger,
		redis:            redis,
		eventInboundRepo: eventInboundRepo,
		whatsappRepo:     whatsappRepo,
	}
}
