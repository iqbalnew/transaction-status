package services

import (
	"fmt"
	"log"

	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	servicelogger "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/service-logger"
	authPB "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/service/stubs/auth"
)

type ServiceConnection struct {
	AuthService *grpc.ClientConn
}

func InitServicesConn(
	logger *servicelogger.AddonsLogrus,
	certFile string,
	authAddress string,
) *ServiceConnection {

	var err error
	var creds credentials.TransportCredentials

	if certFile != "" {
		creds, err = credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Panic(err)
		}
	} else {
		creds = insecure.NewCredentials()
	}

	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(apmgrpc.NewStreamClientInterceptor()),
	)

	services := &ServiceConnection{}

	// Auth Service
	services.AuthService, err = initGrpcClientConn(authAddress, "Auth Service", opts...)
	if err != nil {
		// log.Fatalf("%v", err)
		// os.Exit(1)
		return nil
	}

	return services

}

func initGrpcClientConn(address string, name string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {

	if address == "" {
		// return nil, fmt.Errorf("[service - connection] %s address is empty", name)
		return nil, nil
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("[service - connection] Failed to establish a connection to %s: %v", name, err.Error())
	}

	return conn, nil

}

func (s *ServiceConnection) AuthServiceClient() authPB.ApiServiceClient {
	return authPB.NewApiServiceClient(s.AuthService)
}

func (s *ServiceConnection) CloseAllServicesConn() {
	s.CloseAuthServiceConn()
}

func (s *ServiceConnection) CloseAuthServiceConn() error {
	return s.AuthService.Close()
}
