package rabbitmq

import (
	"context"
	"errors"
	"testing"
	"time"

	rabbitmqmock "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/rabbitmq/mock"
	rabbitmqwrapper "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/rabbitmq/wrapper"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/suite"
)

type RabbitmqTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *RabbitmqTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func TestInitRabbitmq(t *testing.T) {
	suite.Run(t, new(RabbitmqTestSuite))
}

func (s *RabbitmqTestSuite) TestRabbitMq_NewConnection() {
	type expectation struct {
		out *Connection
		err error
	}

	tests := map[string]struct {
		name     string
		queue    string
		config   *Config
		expected expectation
	}{
		"SuccessCreateNewConnection": {
			name:   "Connection",
			queue:  "connection",
			config: &Config{},
			expected: expectation{
				out: &Connection{
					name:   "Connection",
					queue:  "connection",
					config: &Config{},
				},
				err: nil,
			},
		},
		"ConnectionAlreadyExist": {
			name:   "Connection",
			queue:  "connection",
			config: &Config{},
			expected: expectation{
				out: &Connection{
					name:   "Connection",
					queue:  "connection",
					config: &Config{},
				},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			out := NewConnection(tt.name, tt.queue, tt.config, nil)

			if tt.expected.out.name != out.name ||
				tt.expected.out.queue != out.queue {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, out)
			}
		})
	}
}

func (s *RabbitmqTestSuite) TestRabbitMq_GetConnection() {
	type expectation struct {
		out *Connection
		err error
	}

	tests := map[string]struct {
		name     string
		queue    string
		config   *Config
		expected expectation
	}{
		"SuccessGetConnection": {
			name:   "ExistConnection",
			queue:  "exist-connection",
			config: &Config{},
			expected: expectation{
				out: &Connection{
					name:   "ExistConnection",
					queue:  "exist-connection",
					config: &Config{},
				},
				err: nil,
			},
		},
		"ConnectionNotFound": {
			name: "",
			expected: expectation{
				out: nil,
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			var c *Connection
			if tt.name != "" {
				NewConnection(tt.name, tt.queue, tt.config, nil)
			}
			c = GetConnection(tt.name)

			if c != nil {
				if tt.expected.out.name != c.name {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out.name, c.name)
				}
			} else {
				if tt.expected.out != c {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, c)
				}
			}
		})
	}
}

func (s *RabbitmqTestSuite) TestRabbitMq_GetName() {
	type expectation struct {
		out string
		err error
	}

	tests := map[string]struct {
		name     string
		expected expectation
	}{
		"SuccessGetConnectionName": {
			name: "ConnectionAddons",
			expected: expectation{
				out: "ConnectionAddons",
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			c := NewConnection(tt.name, "", &Config{}, nil)
			out := c.GetName()

			if tt.expected.out != out {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, out)
			}
		})
	}
}

func (s *RabbitmqTestSuite) TestRabbitMq_CheckConnection() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		expected expectation
	}{
		"SuccessCheckConnection": {
			expected: expectation{
				err: errors.New("Test Error Channel"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			c := NewConnection("TestErrChannel", "", &Config{}, nil)
			go func() {
				time.Sleep(500 * time.Millisecond)
				c.err <- errors.New("Test Error Channel")
			}()
			out := <-c.CheckConnection()

			if tt.expected.err.Error() != out.Error() {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, out)
			}
		})
	}
}

func (s *RabbitmqTestSuite) TestRabbitMq_Connect() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		name     string
		queue    string
		config   *Config
		expected expectation
	}{
		"SuccessConnect": {
			name:  "Connect",
			queue: "connect-mq",
			config: &Config{
				Host:     "localhost",
				User:     "guest",
				Password: "guest",
			},
			expected: expectation{
				err: nil,
			},
		},
		"FailedConnect": {
			name:  "Connect-Failed",
			queue: "connect-mq-failed",
			config: &Config{
				Host: "",
			},
			expected: expectation{
				err: errors.New("URL cannot be empty"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {

			c := NewConnection(tt.name, tt.queue, tt.config, &rabbitmqmock.RabbitMqMock{})
			err := c.Connect()

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *RabbitmqTestSuite) TestRabbitMq_InitChannel() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		name     string
		queue    string
		config   *Config
		amqpw    rabbitmqwrapper.RabbitMqConnection
		expected expectation
	}{
		"SuccessInitChannel": {
			name:  "InitChannel",
			queue: "init-channel-mq",
			config: &Config{
				Host:     "localhost",
				User:     "guest",
				Password: "guest",
			},
			amqpw: &rabbitmqmock.RabbitMqMock{
				Conn: &amqp.Connection{},
			},
			expected: expectation{
				err: nil,
			},
		},
		"FailedInitChannel": {
			name:  "FailedInitChannel",
			queue: "init-mq-channel-failed",
			config: &Config{
				Host:     "localhost",
				User:     "guest",
				Password: "guest",
			},
			amqpw: &rabbitmqmock.RabbitMqMock{},
			expected: expectation{
				err: errors.New("please create connection first"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			c := NewConnection(tt.name, tt.queue, tt.config, &rabbitmqmock.RabbitMqMock{})
			c.conn = tt.amqpw
			err := c.InitChannel()

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *RabbitmqTestSuite) TestRabbitMq_CloseChannel() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		name     string
		queue    string
		config   *Config
		amqpw    rabbitmqwrapper.RabbitMqChannel
		expected expectation
	}{
		"SuccessCloseChannel": {
			name:  "CloseChannel",
			queue: "close-channel-mq",
			config: &Config{
				Host:     "localhost",
				User:     "guest",
				Password: "guest",
			},
			amqpw: &rabbitmqmock.RabbitMqMock{
				Chann: &amqp.Channel{},
				Conn:  &amqp.Connection{},
			},
			expected: expectation{
				err: nil,
			},
		},
		"FailedCloseChannel": {
			name:  "FailedCloseChannel",
			queue: "close-mq-channel-failed",
			config: &Config{
				Host:     "localhost",
				User:     "guest",
				Password: "guest",
			},
			amqpw: &rabbitmqmock.RabbitMqMock{},
			expected: expectation{
				err: errors.New("channel already closed"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			c := NewConnection(tt.name, tt.queue, tt.config, &rabbitmqmock.RabbitMqMock{})
			c.Channel = tt.amqpw
			err := c.CloseChannel()

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *RabbitmqTestSuite) TestRabbitMq_BindQueue() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		name     string
		queue    string
		config   *Config
		amqpw    rabbitmqwrapper.RabbitMqChannel
		expected expectation
	}{
		"SuccessBindQueue": {
			name:  "BindQueue",
			queue: "bind-queue-mq",
			config: &Config{
				Host:     "localhost",
				User:     "guest",
				Password: "guest",
				QueueConfig: &QueueConfig{
					Durable:       true,
					AutoDelete:    false,
					Exclusive:     false,
					NoWait:        false,
					Args:          nil,
					PrefetchCount: 1,
					PrefetchSize:  0,
					Global:        false,
				},
			},
			amqpw: &rabbitmqmock.RabbitMqMock{
				Chann: &amqp.Channel{},
			},
			expected: expectation{
				err: nil,
			},
		},
		"FailedBindQueue": {
			name:  "FailedBindQueue",
			queue: "bind-queue-failed",
			config: &Config{
				Host:     "localhost",
				User:     "guest",
				Password: "guest",
				QueueConfig: &QueueConfig{
					Durable:       true,
					AutoDelete:    false,
					Exclusive:     false,
					NoWait:        false,
					Args:          nil,
					PrefetchCount: 1,
					PrefetchSize:  0,
					Global:        false,
				},
			},
			amqpw: &rabbitmqmock.RabbitMqMock{},
			expected: expectation{
				err: errors.New("please create channel first"),
			},
		},
		"SuccessBindQueueQos": {
			name:  "BindQueueQos",
			queue: "bind-queue-qos-mq",
			config: &Config{
				Host:     "localhost",
				User:     "guest",
				Password: "guest",
				QueueConfig: &QueueConfig{
					Durable:       true,
					AutoDelete:    false,
					Exclusive:     false,
					NoWait:        false,
					Args:          nil,
					PrefetchCount: -1,
					PrefetchSize:  0,
					Global:        false,
				},
			},
			amqpw: &rabbitmqmock.RabbitMqMock{
				Chann: &amqp.Channel{},
			},
			expected: expectation{
				err: errors.New("prefetchCount cannot less than 1"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			c := NewConnection(tt.name, tt.queue, tt.config, &rabbitmqmock.RabbitMqMock{})
			c.Channel = tt.amqpw
			err := c.BindQueue()

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *RabbitmqTestSuite) TestRabbitMq_Consume() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		name     string
		queue    string
		config   *Config
		amqpw    rabbitmqwrapper.RabbitMqChannel
		expected expectation
	}{
		"SuccessConsume": {
			name:  "Consume",
			queue: "consume-mq",
			amqpw: &rabbitmqmock.RabbitMqMock{
				Chann: &amqp.Channel{},
			},
			expected: expectation{
				err: nil,
			},
		},
		"FailedConsume": {
			name:  "FailedConsume",
			queue: "consume-failed",
			amqpw: nil,
			expected: expectation{
				err: errors.New("consumer channel is nil, please init consumer first"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			c := NewConnection(tt.name, tt.queue, tt.config, &rabbitmqmock.RabbitMqMock{})
			c.Channel = tt.amqpw
			_, err := c.Consume()

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *RabbitmqTestSuite) TestRabbitMq_Publish() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		name     string
		queue    string
		config   *Config
		amqpw    rabbitmqwrapper.RabbitMqChannel
		expected expectation
	}{
		"SuccessPublish": {
			name:  "SuccessPublish",
			queue: "publish-success",
			amqpw: &rabbitmqmock.RabbitMqMock{
				Chann: &amqp.Channel{},
			},
			expected: expectation{
				err: nil,
			},
		},
		"FailedPublishChannelNil": {
			name:  "FailedPublishChannelNil",
			queue: "publish-failed-channel-nil",
			amqpw: nil,
			expected: expectation{
				err: errors.New("producer channel is nil, please init producer first"),
			},
		},
		"FailedPublish": {
			name:  "FailedPublish",
			queue: "publish-failed",
			amqpw: &rabbitmqmock.RabbitMqMock{
				Chann: &amqp.Channel{},
			},
			expected: expectation{
				err: errors.New("publish failed"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			c := NewConnection(tt.name, tt.queue, tt.config, &rabbitmqmock.RabbitMqMock{})
			c.Channel = tt.amqpw
			var err error

			go func() {
				c.err <- errors.New("masuk error")
			}()
			err = c.Publish(s.ctx, "application/json", []byte("unit-test"), 0)

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *RabbitmqTestSuite) TestRabbitMq_CleanUp() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		amqpwCo  rabbitmqwrapper.RabbitMqConnection
		amqpwCh  rabbitmqwrapper.RabbitMqChannel
		expected expectation
	}{
		"SuccessCleanUp": {
			amqpwCo: &rabbitmqmock.RabbitMqMock{
				Conn:  &amqp.Connection{},
				Chann: &amqp.Channel{},
			},
			amqpwCh: &rabbitmqmock.RabbitMqMock{
				Chann: &amqp.Channel{},
				Conn:  &amqp.Connection{},
			},
			expected: expectation{
				err: nil,
			},
		},
		"FailedCleanUpChannelClose": {
			amqpwCh: &rabbitmqmock.RabbitMqMock{},
			amqpwCo: &rabbitmqmock.RabbitMqMock{},
			expected: expectation{
				err: errors.New("channel already closed"),
			},
		},
		"FailedCleanUpConnectionClose": {
			amqpwCh: &rabbitmqmock.RabbitMqMock{
				Conn:  &amqp.Connection{},
				Chann: &amqp.Channel{},
			},
			amqpwCo: &rabbitmqmock.RabbitMqMock{
				Chann: &amqp.Channel{},
			},
			expected: expectation{
				err: errors.New("connection already closed"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			c := &Connection{
				Channel: tt.amqpwCh,
				conn:    tt.amqpwCo,
			}
			err := c.CleanUp()

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}
