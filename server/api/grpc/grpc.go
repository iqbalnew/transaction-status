package apigrpc

import (
	"context"
	"time"

	apigrpchandler "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/api/grpc/handler"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/config"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/db"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/rabbitmq"
	rabbitmqwrapper "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/rabbitmq/wrapper"
	servicelogger "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/service-logger"
	pb "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/pb"
	svc "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/service"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/utils"
	"google.golang.org/protobuf/encoding/protojson"
)

// Server represents the server implementation of the SW API.
type Server struct {
	appConfig   *config.Config
	logger      *servicelogger.AddonsLogrus
	provider    *db.Provider
	svcConn     *svc.ServiceConnection
	grpcHandler *apigrpchandler.GrpcHandler

	pb.TransactionStatusServiceServer
}

func New(
	appConfig *config.Config,
	logger *servicelogger.AddonsLogrus,
	provider *db.Provider,
	svcConn *svc.ServiceConnection,
) (*Server, error) {
	grpcHandler := apigrpchandler.New(appConfig, provider, svcConn, logger)

	server := &Server{
		appConfig:   appConfig,
		provider:    provider,
		svcConn:     svcConn,
		logger:      logger,
		grpcHandler: grpcHandler,
	}

	return server, nil
}
func (x *Server) SetupRabbitMQConn(connectionName string, queueName string) error {

	UserRabbit := utils.GetEnv("USER-RABBIT", "guest")
	PasswordRabbit := utils.GetEnv("PASSWORD-RABBIT", "guest")
	HostRabbit := utils.GetEnv("HOST-RABBIT", "127.0.0.1:5672")

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

	x.logger.Info("Init rabbitmq consumer connection ", rabbit.GetName(), " complete")

	return nil
}

func (x *Server) JobPending() {
	minute := 3
	for {
		x.logger.Println("masuk ")
		PendingJobResult, err := x.provider.GetAllJobTransactionsPending(
			context.Background(),
			&pb.Pagination{
				Limit:  100,
				Page:   1,
				Sort:   "updated_at",
				Dir:    pb.Direction_ASC,
				Filter: " status = 'NEW' and status = 'IN_PROGRESS'",
			})

		if err != nil {
			x.logger.Errorln("failedGet Data", err)
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
			//return nil, err
		}

		if len(PendingJobResult) == 0 {
			x.logger.Errorln("data pending not found wait ", minute, " minute to listing and push to rabbit")
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
		}

		// process data
		x.ProcessJobPending(PendingJobResult)

	}

}

func (x *Server) ProcessJobPending(PendingJobResult []*pb.JobTransactionStatusPending) {

	for _, v := range PendingJobResult {
		x.PublishJob(v)
	}

}

func (x *Server) PublishJob(data *pb.JobTransactionStatusPending) {

	if data.GetId() == 0 {
		x.logger.Error("Id Job 0")
		return
	}

	if data.GetTaskId() == 0 {
		x.logger.Error("Id Job 0")
		return
	}

	var conn *rabbitmq.Connection

	switch data.GetType() {
	case "Kliring":
		conn = rabbitmq.GetConnection("rabbit-conn-publisher-kliring")
	case "Swift":
		conn = rabbitmq.GetConnection("rabbit-conn-publisher-swift")
	default:
		x.logger.Errorln("connection Not Found")
		return
	}

	payload := pb.MassageRabbitPublish{
		IdJob:  data.GetId(),
		TaskId: data.GetTaskId(),
	}

	message, parseErr := protojson.Marshal(&payload)
	if parseErr != nil || data == nil {
		x.logger.Error("Failed parse request:", parseErr)
		return
	}

	err := conn.Publish(context.Background(), "application/json", message, 0)
	if err != nil {
		x.logger.Errorln("failed publish erroro: ")
		return
	}

	//update
	data.Status = pb.StatusInquiryJob_IN_QUEUE
	_, err = x.provider.UpdateJobTransactionPending(context.Background(), data)
	if err != nil {
		x.logger.Errorln("failed updated UpdateJobTransactionPending error :", err)
		return
	}

}
