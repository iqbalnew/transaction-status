package rabbitmq

import (
	"context"
	"errors"
	"fmt"

	rabbitmqwrapper "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/rabbitmq/wrapper"
	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueConfig struct {
	// #region Queue Declare Config
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
	// #endregion

	// #region Qos Config
	PrefetchCount int
	PrefetchSize  int
	Global        bool
	// #endregion
}
type Config struct {
	User           string
	Password       string
	Host           string
	ReconnectDelay int
	*QueueConfig
}

type Connection struct {
	name    string
	queue   string
	config  *Config
	err     chan error
	amqpw   rabbitmqwrapper.RabbitMqInterface
	Channel rabbitmqwrapper.RabbitMqChannel
	conn    rabbitmqwrapper.RabbitMqConnection
}

var (
	connectionPool = make(map[string]*Connection)
)

func NewConnection(name string, queue string, config *Config, amqpw rabbitmqwrapper.RabbitMqInterface) *Connection {
	if c, ok := connectionPool[name]; ok {
		return c
	}

	c := &Connection{
		name:   name,
		queue:  queue,
		config: config,
		err:    make(chan error),
		amqpw:  amqpw,
	}
	connectionPool[name] = c

	return c
}

func GetConnection(name string) *Connection {
	return connectionPool[name]
}

func (c *Connection) GetName() string {
	return c.name
}

func (c *Connection) CheckConnection() chan error {
	return c.err
}

func (c *Connection) Connect() error {
	var err error
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s/",
		c.config.User,
		c.config.Password,
		c.config.Host)
	c.conn, err = c.amqpw.Dial(rabbitUrl)
	if err != nil {
		return err
	}

	go func() {
		<-c.conn.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		c.err <- errors.New("Connection Closed")
	}()

	return nil
}

func (c *Connection) InitChannel() error {
	var channelMqErr error
	c.Channel, channelMqErr = c.conn.Channel()
	if channelMqErr != nil {
		return channelMqErr
	}

	return nil
}

func (c *Connection) CloseChannel() error {
	err := c.Channel.Close()
	if err != nil {
		return errors.New("channel already closed")
	}

	return nil
}

func (c *Connection) BindQueue() error {
	if _, err := c.Channel.QueueDeclare(c.queue, c.config.Durable, c.config.AutoDelete, c.config.Exclusive, c.config.NoWait, c.config.Args); err != nil {
		return err
	}

	if err := c.Channel.Qos(c.config.PrefetchCount, c.config.PrefetchSize, c.config.Global); err != nil {
		return err
	}

	return nil
}

func (c *Connection) Reconnect() error {
	if cErr := c.Connect(); cErr != nil {
		return cErr
	}
	if icErr := c.InitChannel(); icErr != nil {
		return icErr
	}
	if bqErr := c.BindQueue(); bqErr != nil {
		return bqErr
	}
	return nil
}

func (c *Connection) Consume() (<-chan amqp.Delivery, error) {
	if c.Channel == nil {
		return nil, errors.New("consumer channel is nil, please init consumer first")
	}
	return c.Channel.Consume(
		c.queue,
		"",
		false,
		false,
		false,
		false,
		c.config.Args,
	)
}

func (c *Connection) Publish(ctx context.Context, contentType string, qmsg []byte, priority uint8) error {
	if c.Channel == nil {
		return errors.New("producer channel is nil, please init producer first")
	}

	select { //non blocking channel - if there is no error will go to default where we do nothing
	case err := <-c.err:
		if err != nil {
			if reconErr := c.Reconnect(); reconErr != nil {
				return reconErr
			}
		}
	default:
	}

	proErr := c.Channel.PublishWithContext(ctx,
		"",
		c.queue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  contentType,
			Body:         qmsg,
			Priority:     priority,
		})

	if proErr != nil {
		return proErr
	}

	return nil
}

func (c *Connection) CleanUp() error {
	if chErr := c.CloseChannel(); chErr != nil {
		return chErr
	}

	if coErr := c.conn.Close(); coErr != nil {
		return errors.New("connection already closed")
	}
	delete(connectionPool, c.name)
	return nil
}
