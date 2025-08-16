package service

import (
	"context"
	"eventhandler/internal/provider"
	"eventhandler/internal/repository"
	"eventhandler/model"

	"github.com/redis/go-redis/v9"
)

type Events interface {
	HandleEvent(ctx context.Context, data *model.QueueEvent) error
}

type service struct {
	logger provider.ILogger
	repo   repository.EventRepository
	redis  *redis.Client
}

func NewService(logger provider.ILogger, repo repository.EventRepository, redis *redis.Client) Events {
	return &service{
		logger: logger,
		repo:   repo,
		redis:  redis,
	}
}
