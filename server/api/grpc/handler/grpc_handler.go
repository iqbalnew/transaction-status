package apigrpchandler

import (
	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/config"
	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/db"
	servicelogger "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/service-logger"
	svc "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/service"
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
