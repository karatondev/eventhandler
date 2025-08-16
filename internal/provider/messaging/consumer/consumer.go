package consumer

type (
	QosOptions struct {
		PrefetchCount int
		PrefetchSize  int
		Global        bool
	}

	AMQPConsumerOptions struct {
		Qos         *QosOptions
		Concurrency int
	}

	AMQPConsumerOption func(options *AMQPConsumerOptions)
)

func Qos(prefetchCount, prefetchSize int, global bool) AMQPConsumerOption {
	return func(options *AMQPConsumerOptions) {
		options.Qos = &QosOptions{
			PrefetchCount: prefetchCount,
			PrefetchSize:  prefetchSize,
			Global:        global,
		}
	}
}

func Concurrency(concurrency int) AMQPConsumerOption {
	return func(options *AMQPConsumerOptions) {
		options.Concurrency = concurrency
	}
}
