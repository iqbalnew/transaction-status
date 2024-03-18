package apihttp

import servicelogger "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/service-logger"

type HttpHandler struct {
	Logger *servicelogger.AddonsLogrus
}

func New(
	logger *servicelogger.AddonsLogrus,
) *HttpHandler {
	return &HttpHandler{
		Logger: logger,
	}
}
