package apigrpchandler

import (
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/config"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/db"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/rabbitmq"
	rabbitmqwrapper "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/rabbitmq/wrapper"
	servicelogger "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/service-logger"
	svc "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/service"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/utils"
)

type GrpcHandler struct {
	appConfig *config.Config
	provider  *db.Provider
	svcConn   *svc.ServiceConnection
	logger    *servicelogger.AddonsLogrus
}

func New(appConfig *config.Config, provider *db.Provider, svcConn *svc.ServiceConnection, logger *servicelogger.AddonsLogrus) *GrpcHandler {
	return &GrpcHandler{
		appConfig: appConfig,
		provider:  provider,
		svcConn:   svcConn,
		logger:    logger,
	}
}

func (x *GrpcHandler) SetupRabbitMQConn(connectionName string) error {

	queueName := utils.GetEnv("queue-inquiry-status-swift", "queue-inquiry-status-swift")

	UserRabbit := utils.GetEnv("AMQP_HOST", "guest")
	PasswordRabbit := utils.GetEnv("AMQP_USER", "user")
	HostRabbit := utils.GetEnv("AMQP_PASSWORD", "127.0.0.1:5672")

	rabbitmqConfig := &rabbitmq.Config{
		User:           UserRabbit,
		Password:       PasswordRabbit,
		Host:           HostRabbit,
		ReconnectDelay: 120,
	}

	rabbitmqConfig.QueueConfig = &rabbitmq.QueueConfig{
		Durable:       true,
		AutoDelete:    false,
		Exclusive:     false,
		NoWait:        false,
		Args:          nil,
		PrefetchCount: 1,
		PrefetchSize:  0,
		Global:        false,
	}

	rabbit := rabbitmq.NewConnection(
		connectionName,
		queueName,
		rabbitmqConfig,
		&rabbitmqwrapper.RabbitMqWrapper{},
	)

	if err := rabbit.Connect(); err != nil {
		return err
	}

	if mqConChanErr := rabbit.InitChannel(); mqConChanErr != nil {
		return mqConChanErr
	}
	if mqConErr := rabbit.BindQueue(); mqConErr != nil {
		return mqConErr
	}

	//c.Logger.Info("Init rabbitmq consumer connection ", c.Name, " complete")

	return nil
}
