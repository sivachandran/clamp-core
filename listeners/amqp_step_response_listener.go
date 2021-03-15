package listeners

import "clamp-core/config"

type AMQPStepResponseListenerInterface interface {
	Listen()
}

var AMQPStepResponseListener AMQPStepResponseListenerInterface

func init() {
	if config.ENV.QueueDriver == "amqp" {
		AMQPStepResponseListener = &amqpListener{}
	}
}
