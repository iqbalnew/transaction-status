package interceptors

import (
	"context"
	"strings"

	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/constant"
	svc "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/service"
	authPB "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/service/stubs/auth"
	"golang.org/x/exp/slices"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type AuthInterceptor struct {
	accessibleRoles map[string][]string
	svcConn         *svc.ServiceConnection
}

func NewAuthInterceptor(apiServicePath string, svcConn *svc.ServiceConnection) *AuthInterceptor {
	return &AuthInterceptor{
		svcConn:         svcConn,
		accessibleRoles: accessibleRoles(apiServicePath),
	}
}

// Filter access by role
func accessibleRoles(apiServicePath string) map[string][]string {

	// restricted api
	return map[string][]string{
		apiServicePath + "GetAllTemplates":   {},
		apiServicePath + "GetTemplateDetail": {},
		apiServicePath + "SaveTemplate":      {},
		apiServicePath + "UpdateTemplate":    {},
		apiServicePath + "DeleteTemplate":    {},
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if !interceptor.isRestricted(info.FullMethod) {
			return handler(ctx, req)
		}

		claims, err := interceptor.claimsToken(ctx)
		if err != nil {
			return nil, err
		}

		err = interceptor.authorize(ctx, claims, info.FullMethod)
		if err != nil {
			return nil, err
		}

		jsonValue, jsonErr := protojson.Marshal(claims)
		if jsonErr != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, constant.CtxTokenKey, jsonValue)

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {

		if !interceptor.isRestricted(info.FullMethod) {
			return handler(srv, stream)
		}

		claims, err := interceptor.claimsToken(stream.Context())
		if err != nil {
			return err
		}

		err = interceptor.authorize(stream.Context(), claims, info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (interceptor *AuthInterceptor) isRestricted(method string) bool {
	_, restricted := interceptor.accessibleRoles[method]
	return restricted
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, claims *authPB.VerifyTokenRes, method string) error {
	featureRoles := []string{}
	for _, v := range claims.GetProductRoles() {
		if strings.EqualFold(constant.ProductName, v.GetProductName()) {
			featureRoles = v.GetAuthorities()
			break
		}
	}

	accessibleRoles, ok := interceptor.accessibleRoles[method]
	if !ok {
		// everyone can access
		return nil
	}

	if len(accessibleRoles) < 1 {
		return nil
	}

	for _, role := range accessibleRoles {
		if slices.Contains(featureRoles, role) {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "Access denied")
}

func (interceptor *AuthInterceptor) claimsToken(ctx context.Context) (*authPB.VerifyTokenRes, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	split := strings.Split(values[0], " ")
	accessToken := split[0]
	if len(split) > 1 {
		accessToken = split[1]
	}
	claims, err := interceptor.GetMeFromAuthService(ctx, accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return claims, nil
}

func (interceptor *AuthInterceptor) GetMeFromAuthService(ctx context.Context, accessToken string) (*authPB.VerifyTokenRes, error) {
	authClient := interceptor.svcConn.AuthServiceClient()

	dataUser, err := authClient.VerifyToken(ctx, &authPB.VerifyTokenReq{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}
	if dataUser == nil {
		return nil, status.Errorf(codes.Aborted, "Failed To Get Data User")
	}

	return dataUser, nil
}
