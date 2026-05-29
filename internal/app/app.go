package app

import (
	"context"
	"eventhandler/internal/handler"
	"eventhandler/internal/provider"
	"eventhandler/internal/repository"
	"eventhandler/internal/service"
	"eventhandler/util"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"zaplio/shared/constant"
	"zaplio/shared/pkg/amqpx"
	"zaplio/shared/pkg/logger"
	"zaplio/shared/pkg/messaging"
	"zaplio/shared/pkg/messaging/consumer"
)

func Run(cfg *util.Config) {
	ctx := context.WithValue(context.Background(), constant.CtxReqIDKey, "MAIN")

	log := logger.NewLogger(logger.LoggerConfig{
		Dir:        util.Configuration.Logger.Dir,
		FileName:   util.Configuration.Logger.FileName,
		MaxBackups: util.Configuration.Logger.MaxBackups,
		MaxSize:    util.Configuration.Logger.MaxSize,
		MaxAge:     util.Configuration.Logger.MaxAge,
		Compress:   util.Configuration.Logger.Compress,
		LocalTime:  util.Configuration.Logger.LocalTime,
		Level:      util.Configuration.Logger.Level,
	})

	db, err := provider.NewPostgresConnection(ctx)
	if err != nil {
		log.Errorfctx(logger.AppLog, ctx, false, "Failed connect to PostgreSQL: %v", err)
		return
	}

	redis, err := provider.NewRedisConnection(ctx)
	if err != nil {
		log.Errorfctx(logger.AppLog, ctx, false, "Failed connect to Redis: %v", err)
		return
	}

	conn, err := provider.NewAMQPConn()
	if err != nil {
		log.Errorfctx(logger.AppLog, ctx, false, "Failed connect to AMQP:", err)
	}
	log.Infofctx(logger.AppLog, ctx, "Application started")

	repo := repository.NewInboundOutboundRepository(log, db)
	whatsappRepo := repository.NewWhatsAppAccountRepository(db)
	svc := service.NewService(log, redis, repo, whatsappRepo)
	consumerHandler := handler.NewConsumerHandler(log, svc)

	consumers := []*messaging.AMQPWorkerConsumer{}
	go func() {
		for _, queue := range util.Configuration.QueuesList() {
			c := messaging.NewAMQPWorkerConsumer(
				conn,
				func() messaging.AMQPWorker { return consumerHandler },
				consumer.Qos(util.Configuration.AMQP.PrefetchCount, util.Configuration.AMQP.PrefetchSize, util.Configuration.AMQP.Global),
				consumer.Concurrency(util.Configuration.AMQP.Concurrency),
			)

			if err := c.Consume(queue); err != nil {
				log.Errorfctx(logger.AppLog, ctx, false, "Failed to create consumer for queue %s: %v", queue, err)
			}

			log.Infofctx(logger.AppLog, ctx, "Successfully creating consumer for queue: %s", queue)
			consumers = append(consumers, c)
		}
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownCh
	log.Infofctx(logger.AppLog, ctx, "Receiving signal: %s", sig)

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
				log.Errorfctx(logger.AppLog, ctx, false, "Failed to close AMQP connection: %v", err)
			}
		}

		log.Infofctx(logger.AppLog, ctx, "Successfully stop Application.")

	}(&consumers, conn)

}
