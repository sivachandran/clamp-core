package executors

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
)

type AMQPVal struct {
	ConnectionURL string `json:"connection_url" binding:"required"`
	QueueName     string `json:"queue_name"`
	ExchangeName  string `json:"exchange_name"`
	RoutingKey    string `json:"routing_key"`
	ContentType   string `json:"content_type"`
}

func (val AMQPVal) DoExecute(requestBody interface{}) (interface{}, error) {
	prefix := log.Prefix()
	log.SetPrefix("")
	log.Printf("%s AMQP Executor: Executing amqp %s body:%v", prefix, val.getName(), requestBody)

	conn, err := amqp.Dial(val.ConnectionURL)
	if err != nil {
		log.Printf("%s AMQP Error: %s", prefix, err.Error())
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%s AMQP Error: %s", prefix, err.Error())
		return nil, err
	}
	defer ch.Close()

	if val.ExchangeName != "" {
		return sendMessageToExchange(ch, val, requestBody, prefix)
	} else if val.QueueName != "" {
		return sendMessageToQueue(ch, val, requestBody, prefix)
	} else {
		return nil, errors.New("AMQP - queue/exchange name not specified")
	}
}

func sendMessageToQueue(ch *amqp.Channel, val AMQPVal, body interface{}, prefix string) (interface{}, error) {
	bytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	err = ch.Publish(
		"",
		val.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: val.ContentType,
			Body:        bytes,
		})
	if err != nil {
		return nil, err
	} else {
		log.Printf("%s AMQP Executor: pushed message successfully", prefix)
	}
	return nil, nil
}

func sendMessageToExchange(ch *amqp.Channel, val AMQPVal, body interface{}, prefix string) (interface{}, error) {
	bytes, err := json.Marshal(body)
	err = ch.Publish(
		val.ExchangeName,
		val.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: val.ContentType,
			Body:        bytes,
		})
	if err != nil {
		return nil, err
	} else {
		log.Printf("%s AMQP Executor: pushed message successfully", prefix)
	}
	return nil, nil
}

func (val AMQPVal) getName() string {
	if val.QueueName != "" {
		return val.QueueName
	}
	return val.ExchangeName
}
