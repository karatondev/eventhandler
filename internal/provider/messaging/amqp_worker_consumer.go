//go:build !mock
// +build !mock

package messaging

import (
	"eventhandler/internal/provider/amqpx"
	"eventhandler/internal/provider/messaging/consumer"
	"fmt"
	"sync"
	"sync/atomic"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	AMQPWorker interface {
		Handle(amqp.Delivery)
	}

	AMQPWorkerConsumer struct {
		opts          *consumer.AMQPConsumerOptions
		newAMQPWorker NewAMQPWorker
		wg            *sync.WaitGroup
		state         int32
		stop          chan struct{}
		channel       *amqpx.Channel
	}

	NewAMQPWorker func() AMQPWorker
)

func NewAMQPWorkerConsumer(conn amqpx.ChannelReader, newAMQPWorker NewAMQPWorker, options ...consumer.AMQPConsumerOption) *AMQPWorkerConsumer {
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	opts := &consumer.AMQPConsumerOptions{
		Concurrency: 1,
	}
	for _, option := range options {
		option(opts)
	}

	if opts.Qos == nil {
		opts.Qos = &consumer.QosOptions{
			PrefetchCount: 10,
			PrefetchSize:  0,
			Global:        false,
		}
	}

	return &AMQPWorkerConsumer{
		opts:          opts,
		newAMQPWorker: newAMQPWorker,
		wg:            &sync.WaitGroup{},
		channel:       channel,
	}
}

func (c *AMQPWorkerConsumer) Consume(queue string) error {
	if !atomic.CompareAndSwapInt32(&c.state, 0, 1) {
		return nil
	}

	_, _ = c.channel.QueueDeclare(queue, true, false, false, false, nil)

	c.stop = make(chan struct{})

	if c.opts.Qos != nil {
		err := c.channel.Qos(c.opts.Qos.PrefetchCount, c.opts.Qos.PrefetchSize, c.opts.Qos.Global)
		if err != nil {
			return err
		}
	}

	deliveries, err := c.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	concurrency := c.opts.Concurrency
	if c.opts.Qos.PrefetchCount < concurrency {
		concurrency = c.opts.Qos.PrefetchCount
	}

	for i := 0; i < concurrency; i++ {
		c.wg.Add(1)

		name := fmt.Sprintf("%s_%d", queue, i)
		worker := c.newAMQPWorker()
		go c.consume(name, c.wg, c.channel, worker.Handle, deliveries, c.stop)
	}

	return nil
}

func (c *AMQPWorkerConsumer) consume(name string, wg *sync.WaitGroup, channel *amqpx.Channel, handler func(delivery amqp.Delivery), deliveries <-chan amqp.Delivery, stop chan struct{}) {
	defer func(g *sync.WaitGroup, channel *amqpx.Channel) {
		if err := recover(); err != nil {
			c.consume(name, wg, channel, handler, deliveries, stop)
		} else {
			// if !channel.IsClosed() {
			// 	_ = channel.Close()
			// }
			wg.Done()
		}
	}(wg, channel)

	for {
		select {
		case delivery, ok := <-deliveries:
			if !ok || channel.IsClosed() {
				atomic.StoreInt32(&c.state, 0)
				return
			}

			handler(delivery)
		case <-stop:
			return
		}
	}
}

func (c *AMQPWorkerConsumer) Stop() {
	if !atomic.CompareAndSwapInt32(&c.state, 1, 0) {
		return
	}
	close(c.stop)

	c.wg.Wait()
}
