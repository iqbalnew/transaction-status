package grpcapimock

import (
	"context"

	authPB "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/service/stubs/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcApiMock struct {
	authPB.ApiServiceServer
}

func (gam *GrpcApiMock) VerifyToken(ctx context.Context, req *authPB.VerifyTokenReq) (*authPB.VerifyTokenRes, error) {
	if req.GetAccessToken() == "FailedAccessToken" {
		return nil, status.Error(codes.InvalidArgument, "token invalid")
	}

	return &authPB.VerifyTokenRes{
		IsValid: false,
	}, nil
}
