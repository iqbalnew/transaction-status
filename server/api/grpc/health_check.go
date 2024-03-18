package apigrpc

import (
	"context"

	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) HealthCheck(ctx context.Context, _ *emptypb.Empty) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Message: "API Running !"}, nil
}
