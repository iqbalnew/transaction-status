package rabbitmqmock

import (
	"context"
	"errors"
	"strings"

	rabbitmqwrapper "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/rabbitmq/wrapper"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqMock struct {
	Conn  *amqp.Connection
	Chann *amqp.Channel
	rabbitmqwrapper.RabbitMqInterface
	rabbitmqwrapper.RabbitMqConnection
	rabbitmqwrapper.RabbitMqChannel
}

func (rm *RabbitMqMock) Dial(url string) (*amqp.Connection, error) {
	if url == "amqp://:@/" {
		return nil, errors.New("URL cannot be empty")
	}

	return &amqp.Connection{}, nil
}

func (rm *RabbitMqMock) Channel() (*amqp.Channel, error) {
	if rm.Conn == nil {
		return nil, errors.New("please create connection first")
	}
	return &amqp.Channel{}, nil
}

func (rm *RabbitMqMock) Close() error {
	if rm.Chann == nil {
		return errors.New("channel already closed")
	}

	if rm.Conn == nil {
		return errors.New("connection already closed")
	}

	return nil
}

func (rm *RabbitMqMock) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	if rm.Chann == nil {
		return amqp.Queue{}, errors.New("please create channel first")
	}
	return amqp.Queue{}, nil
}

func (rm *RabbitMqMock) Qos(prefetchCount int, prefetchSize int, global bool) error {
	if prefetchCount < 1 {
		return errors.New("prefetchCount cannot less than 1")
	}
	return nil
}

func (rm *RabbitMqMock) Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return nil, nil
}

func (rm *RabbitMqMock) PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	if strings.EqualFold("publish-failed", key) {
		return errors.New("publish failed")
	}

	return nil
}
