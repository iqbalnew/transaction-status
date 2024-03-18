package apigrpchandler

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (gh *GrpcHandler) NotImplementedError() error {
	st := status.New(codes.Unimplemented, "Not implemented yet")
	return st.Err()
}

func (gh *GrpcHandler) UnauthorizedError() error {
	st := status.New(codes.Unauthenticated, "Unauthorized")
	return st.Err()
}

func (gh *GrpcHandler) ServerError() error {
	st := status.New(codes.Internal, "Internal Error")
	return st.Err()
}

func (gh *GrpcHandler) ServerErrorWithDetails(err error) error {
	st := status.New(codes.Internal, err.Error())
	return st.Err()
}

func (gh *GrpcHandler) InvalidArgument() error {
	st := status.New(codes.InvalidArgument, "Bad Request")
	return st.Err()
}

func (gh *GrpcHandler) NotFoundError() error {
	st := status.New(codes.NotFound, "Not Found")
	return st.Err()
}

func (gh *GrpcHandler) PermissionDeniedError() error {
	st := status.New(codes.PermissionDenied, "Permission Denied")
	return st.Err()
}
