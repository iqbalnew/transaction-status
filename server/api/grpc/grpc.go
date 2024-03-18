package apigrpc

import (
	apigrpchandler "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/api/grpc/handler"
	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/config"
	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/db"
	servicelogger "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/service-logger"
	pb "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/pb"
	svc "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/service"
)

// Server represents the server implementation of the SW API.
type Server struct {
	appConfig   *config.Config
	logger      *servicelogger.AddonsLogrus
	provider    *db.Provider
	svcConn     *svc.ServiceConnection
	grpcHandler *apigrpchandler.GrpcHandler

	pb.TemplateServiceServer
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
