package listeners

import "clamp-core/config"

type KafkaStepResponseListenerInterface interface {
	Listen()
}

var KafkaStepResponseListener KafkaStepResponseListenerInterface

func init() {
	if config.ENV.KafkaDriver == "kafka" {
		KafkaStepResponseListener = &Consumer{}
	}
}
