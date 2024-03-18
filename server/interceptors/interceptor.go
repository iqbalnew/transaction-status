package interceptors

import (
	"fmt"

	servicelogger "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/service-logger"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"

	"context"
	"time"
)

type Interceptor struct {
	logger *servicelogger.AddonsLogrus
}

func NewInterceptor(logger *servicelogger.AddonsLogrus) *Interceptor {
	return &Interceptor{
		logger: logger,
	}
}

// Interceptors implements the grpc.UnaryServerInteceptor function to add
// interceptors around all gRPC unary calls
func (i *Interceptor) UnaryInterceptors(
	authI *AuthInterceptor,
) grpc.UnaryServerInterceptor {
	return grpc_middleware.ChainUnaryServer(
		apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
		i.LoggingInterceptor,
		i.ErrorsInterceptor,
		authI.Unary(),
	)
}

// Interceptors implements the grpc.StreamServerInteceptor function to add
// interceptors around all gRPC stream calls
func (i *Interceptor) StreamInterceptors(
	authI *AuthInterceptor,
) grpc.StreamServerInterceptor {
	return grpc_middleware.ChainStreamServer(
		authI.Stream(),
	)
}

// ErrorsInterceptor adds error type checking to see if there are any known types
// what we return different grpc error codes for, for example: NotFound resources.
func (i *Interceptor) ErrorsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (out interface{}, err error) {
	out, err = handler(ctx, req)

	return out, err
}

// LoggingInterceptor adds logging around every gRPC call. It includes the method name and timing information.
// if the given handler raises an error, it also appends that to a key.
func (i *Interceptor) LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (out interface{}, err error) {

	if info.FullMethod == "/grpc.health.v1.Health/Check" {

		out, err = handler(ctx, req)
		return out, err

	} else {

		entry := i.logger.WithField(i.logger.GrpcMetadataKey, fmt.Sprintf("method: %s", info.FullMethod))
		start := time.Now()
		out, err = handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			entry = entry.WithError(err)
		}

		entry.WithField(i.logger.GrpcMetadataKey, fmt.Sprintf("duration: %s", duration.String())).Info("finished RPC")
		return out, err

	}

}
