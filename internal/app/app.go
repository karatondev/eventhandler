package app

import (
	"context"
	"eventhandler/internal/handler"
	"eventhandler/internal/provider"
	"eventhandler/internal/provider/amqpx"
	"eventhandler/internal/provider/messaging"
	"eventhandler/internal/provider/messaging/consumer"
	"eventhandler/internal/repository"
	"eventhandler/internal/service"
	"eventhandler/model/constant"
	"eventhandler/util"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Run(cfg *util.Config) {
	ctx := context.WithValue(context.Background(), constant.CtxReqIDKey, "MAIN")

	logger := provider.NewLogger()

	db, err := provider.NewPostgresConnection(ctx)
	if err != nil {
		logger.Errorfctx(provider.AppLog, ctx, false, "Failed connect to PostgreSQL: %v", err)
		return
	}

	redis, err := provider.NewRedisConnection(ctx)
	if err != nil {
		logger.Errorfctx(provider.AppLog, ctx, false, "Failed connect to Redis: %v", err)
		return
	}

	conn, err := provider.NewAMQPConn()
	if err != nil {
		logger.Errorfctx(provider.AppLog, ctx, false, "Failed connect to AMQP:", err)
	}
	logger.Infofctx(provider.AppLog, ctx, "Application started")

	repo := repository.NewEventRepository(logger, db)
	svc := service.NewService(logger, repo, redis)
	consumerHandler := handler.NewConsumerHandler(logger, svc)

	consumers := []*messaging.AMQPWorkerConsumer{}
	go func() {
		for _, queue := range util.Configuration.QueuesList() {
			consumer := messaging.NewAMQPWorkerConsumer(
				conn,
				func() messaging.AMQPWorker { return consumerHandler },
				consumer.Qos(util.Configuration.AMQP.PrefetchCount, util.Configuration.AMQP.PrefetchSize, util.Configuration.AMQP.Global),
				consumer.Concurrency(util.Configuration.AMQP.Concurrency),
			)

			if err := consumer.Consume(queue); err != nil {
				logger.Errorfctx(provider.AppLog, ctx, false, "Failed to create consumer for queue %s: %v", queue, err)
			}

			logger.Infofctx(provider.AppLog, ctx, "Successfully creating consumer for queue: %s", queue)
			consumers = append(consumers, consumer)
		}
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownCh
	logger.Infofctx(provider.AppLog, ctx, "Receiving signal: %s", sig)

	func(consumer *[]*messaging.AMQPWorkerConsumer, amqConn amqpx.ChannelReaderCloser) {
		wg := sync.WaitGroup{}
		if consumer != nil {
			for _, v := range *consumer {
				wg.Add(1)
				go func(v *messaging.AMQPWorkerConsumer) {
					defer wg.Done()
					v.Stop()
				}(v)
			}
		}
		wg.Wait()

		if amqConn != nil {
			if err := amqConn.Close(); err != nil {
				logger.Errorfctx(provider.AppLog, ctx, false, "Failed to close AMQP connection: %v", err)
			}
		}

		logger.Infofctx(provider.AppLog, ctx, "Successfully stop Application.")

	}(&consumers, conn)

}
