package service

import (
	"context"
	"eventhandler/internal/repository"
	"eventhandler/model/entity"

	sharedmodel "zaplio/shared/model"
	"zaplio/shared/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type Events interface {
	HandleEvent(ctx context.Context, data *sharedmodel.QueueEvent) error
	GetAccountBySenderJID(ctx context.Context, senderJID string) (*entity.WhatsAppAccount, error)
}

type service struct {
	logger          logger.ILogger
	redis           *redis.Client
	inboundOutbound repository.InboundOutboundRepository
	whatsappRepo    repository.WhatsAppAccountRepository
}

func NewService(log logger.ILogger, redis *redis.Client, inboundOutbound repository.InboundOutboundRepository, whatsappRepo repository.WhatsAppAccountRepository) Events {
	return &service{
		logger:          log,
		redis:           redis,
		inboundOutbound: inboundOutbound,
		whatsappRepo:    whatsappRepo,
	}
}
