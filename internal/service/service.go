package service

import (
	"context"
	"eventhandler/internal/provider"
	"eventhandler/internal/repository"
	"eventhandler/model"
	"eventhandler/model/entity"

	"github.com/redis/go-redis/v9"
)

type Events interface {
	HandleEvent(ctx context.Context, data *model.QueueEvent) error
	GetAccountBySenderJID(ctx context.Context, senderJID string) (*entity.WhatsAppAccount, error)
}

type service struct {
	logger          provider.ILogger
	redis           *redis.Client
	inboundOutbound repository.InboundOutboundRepository
	whatsappRepo    repository.WhatsAppAccountRepository
}

func NewService(logger provider.ILogger, redis *redis.Client, inboundOutbound repository.InboundOutboundRepository, whatsappRepo repository.WhatsAppAccountRepository) Events {
	return &service{
		logger:          logger,
		redis:           redis,
		inboundOutbound: inboundOutbound,
		whatsappRepo:    whatsappRepo,
	}
}
