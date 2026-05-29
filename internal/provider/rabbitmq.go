package provider

import (
	"fmt"

	"zaplio/shared/pkg/amqpx"

	"eventhandler/util"
)

func NewAMQPConn() (amqpx.ChannelReaderCloser, error) {
	cfg := util.Configuration.AMQP
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d", cfg.Scheme, cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	return amqpx.Dial(dsn)
}
