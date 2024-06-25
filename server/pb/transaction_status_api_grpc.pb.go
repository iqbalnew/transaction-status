// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.2
// source: transaction_status_api.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	TransactionStatusService_HealthCheck_FullMethodName       = "/transaction_status.service.v1.TransactionStatusService/HealthCheck"
	TransactionStatusService_GetAllTemplates_FullMethodName   = "/transaction_status.service.v1.TransactionStatusService/GetAllTemplates"
	TransactionStatusService_GetTemplateDetail_FullMethodName = "/transaction_status.service.v1.TransactionStatusService/GetTemplateDetail"
	TransactionStatusService_SaveTemplate_FullMethodName      = "/transaction_status.service.v1.TransactionStatusService/SaveTemplate"
	TransactionStatusService_UpdateTemplate_FullMethodName    = "/transaction_status.service.v1.TransactionStatusService/UpdateTemplate"
	TransactionStatusService_DeleteTemplate_FullMethodName    = "/transaction_status.service.v1.TransactionStatusService/DeleteTemplate"
)

// TransactionStatusServiceClient is the client API for TransactionStatusService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransactionStatusServiceClient interface {
	HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthCheckResponse, error)
	GetAllTemplates(ctx context.Context, in *GetAllTemplatesRequest, opts ...grpc.CallOption) (*GetAllTemplatesResponse, error)
	GetTemplateDetail(ctx context.Context, in *GetTemplateDetailRequest, opts ...grpc.CallOption) (*GetTemplateDetailResponse, error)
	SaveTemplate(ctx context.Context, in *SaveTemplateRequest, opts ...grpc.CallOption) (*GeneralBodyResponse, error)
	UpdateTemplate(ctx context.Context, in *UpdateTemplateRequest, opts ...grpc.CallOption) (*GetTemplateDetailResponse, error)
	DeleteTemplate(ctx context.Context, in *DeleteTemplateRequest, opts ...grpc.CallOption) (*GeneralBodyResponse, error)
}

type transactionStatusServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionStatusServiceClient(cc grpc.ClientConnInterface) TransactionStatusServiceClient {
	return &transactionStatusServiceClient{cc}
}

func (c *transactionStatusServiceClient) HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, TransactionStatusService_HealthCheck_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionStatusServiceClient) GetAllTemplates(ctx context.Context, in *GetAllTemplatesRequest, opts ...grpc.CallOption) (*GetAllTemplatesResponse, error) {
	out := new(GetAllTemplatesResponse)
	err := c.cc.Invoke(ctx, TransactionStatusService_GetAllTemplates_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionStatusServiceClient) GetTemplateDetail(ctx context.Context, in *GetTemplateDetailRequest, opts ...grpc.CallOption) (*GetTemplateDetailResponse, error) {
	out := new(GetTemplateDetailResponse)
	err := c.cc.Invoke(ctx, TransactionStatusService_GetTemplateDetail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionStatusServiceClient) SaveTemplate(ctx context.Context, in *SaveTemplateRequest, opts ...grpc.CallOption) (*GeneralBodyResponse, error) {
	out := new(GeneralBodyResponse)
	err := c.cc.Invoke(ctx, TransactionStatusService_SaveTemplate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionStatusServiceClient) UpdateTemplate(ctx context.Context, in *UpdateTemplateRequest, opts ...grpc.CallOption) (*GetTemplateDetailResponse, error) {
	out := new(GetTemplateDetailResponse)
	err := c.cc.Invoke(ctx, TransactionStatusService_UpdateTemplate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionStatusServiceClient) DeleteTemplate(ctx context.Context, in *DeleteTemplateRequest, opts ...grpc.CallOption) (*GeneralBodyResponse, error) {
	out := new(GeneralBodyResponse)
	err := c.cc.Invoke(ctx, TransactionStatusService_DeleteTemplate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionStatusServiceServer is the server API for TransactionStatusService service.
// All implementations must embed UnimplementedTransactionStatusServiceServer
// for forward compatibility
type TransactionStatusServiceServer interface {
	HealthCheck(context.Context, *emptypb.Empty) (*HealthCheckResponse, error)
	GetAllTemplates(context.Context, *GetAllTemplatesRequest) (*GetAllTemplatesResponse, error)
	GetTemplateDetail(context.Context, *GetTemplateDetailRequest) (*GetTemplateDetailResponse, error)
	SaveTemplate(context.Context, *SaveTemplateRequest) (*GeneralBodyResponse, error)
	UpdateTemplate(context.Context, *UpdateTemplateRequest) (*GetTemplateDetailResponse, error)
	DeleteTemplate(context.Context, *DeleteTemplateRequest) (*GeneralBodyResponse, error)
	mustEmbedUnimplementedTransactionStatusServiceServer()
}

// UnimplementedTransactionStatusServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTransactionStatusServiceServer struct {
}

func (UnimplementedTransactionStatusServiceServer) HealthCheck(context.Context, *emptypb.Empty) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedTransactionStatusServiceServer) GetAllTemplates(context.Context, *GetAllTemplatesRequest) (*GetAllTemplatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllTemplates not implemented")
}
func (UnimplementedTransactionStatusServiceServer) GetTemplateDetail(context.Context, *GetTemplateDetailRequest) (*GetTemplateDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTemplateDetail not implemented")
}
func (UnimplementedTransactionStatusServiceServer) SaveTemplate(context.Context, *SaveTemplateRequest) (*GeneralBodyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveTemplate not implemented")
}
func (UnimplementedTransactionStatusServiceServer) UpdateTemplate(context.Context, *UpdateTemplateRequest) (*GetTemplateDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTemplate not implemented")
}
func (UnimplementedTransactionStatusServiceServer) DeleteTemplate(context.Context, *DeleteTemplateRequest) (*GeneralBodyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTemplate not implemented")
}
func (UnimplementedTransactionStatusServiceServer) mustEmbedUnimplementedTransactionStatusServiceServer() {
}

// UnsafeTransactionStatusServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransactionStatusServiceServer will
// result in compilation errors.
type UnsafeTransactionStatusServiceServer interface {
	mustEmbedUnimplementedTransactionStatusServiceServer()
}

func RegisterTransactionStatusServiceServer(s grpc.ServiceRegistrar, srv TransactionStatusServiceServer) {
	s.RegisterService(&TransactionStatusService_ServiceDesc, srv)
}

func _TransactionStatusService_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionStatusServiceServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransactionStatusService_HealthCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionStatusServiceServer).HealthCheck(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionStatusService_GetAllTemplates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllTemplatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionStatusServiceServer).GetAllTemplates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransactionStatusService_GetAllTemplates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionStatusServiceServer).GetAllTemplates(ctx, req.(*GetAllTemplatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionStatusService_GetTemplateDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTemplateDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionStatusServiceServer).GetTemplateDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransactionStatusService_GetTemplateDetail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionStatusServiceServer).GetTemplateDetail(ctx, req.(*GetTemplateDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionStatusService_SaveTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveTemplateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionStatusServiceServer).SaveTemplate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransactionStatusService_SaveTemplate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionStatusServiceServer).SaveTemplate(ctx, req.(*SaveTemplateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionStatusService_UpdateTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTemplateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionStatusServiceServer).UpdateTemplate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransactionStatusService_UpdateTemplate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionStatusServiceServer).UpdateTemplate(ctx, req.(*UpdateTemplateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionStatusService_DeleteTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTemplateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionStatusServiceServer).DeleteTemplate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransactionStatusService_DeleteTemplate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionStatusServiceServer).DeleteTemplate(ctx, req.(*DeleteTemplateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TransactionStatusService_ServiceDesc is the grpc.ServiceDesc for TransactionStatusService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TransactionStatusService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "transaction_status.service.v1.TransactionStatusService",
	HandlerType: (*TransactionStatusServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _TransactionStatusService_HealthCheck_Handler,
		},
		{
			MethodName: "GetAllTemplates",
			Handler:    _TransactionStatusService_GetAllTemplates_Handler,
		},
		{
			MethodName: "GetTemplateDetail",
			Handler:    _TransactionStatusService_GetTemplateDetail_Handler,
		},
		{
			MethodName: "SaveTemplate",
			Handler:    _TransactionStatusService_SaveTemplate_Handler,
		},
		{
			MethodName: "UpdateTemplate",
			Handler:    _TransactionStatusService_UpdateTemplate_Handler,
		},
		{
			MethodName: "DeleteTemplate",
			Handler:    _TransactionStatusService_DeleteTemplate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "transaction_status_api.proto",
}