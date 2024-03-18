package apihttp

import (
	"context"
	"encoding/json"
	"net/http"

	pb "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (hh *HttpHandler) CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {

	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Set("Content-type", "application/json")
	grpcErrorCode := status.Code(err)
	httpErrorCode := runtime.HTTPStatusFromCode(grpcErrorCode)
	msg := status.Convert(err).Message()
	w.WriteHeader(httpErrorCode)

	switch grpcErrorCode {
	case codes.InvalidArgument:
		msg = "Invalid argument"
	case codes.Unauthenticated:
		msg = "Permission Denied"
	default:
		msg = "Internal server error"
	}

	hh.Logger.Infof("[DEBUG] Error Gateway: %s, GRPC Error Code: %v, HTTP Error Code: %v, error message: %s", msg, grpcErrorCode, httpErrorCode, err.Error())

	body := &pb.GeneralBodyResponse{
		Error:   true,
		Code:    uint32(httpErrorCode),
		Message: msg,
	}

	jErr := json.NewEncoder(w).Encode(body)

	if jErr != nil {
		w.Write([]byte(fallback))
	}

}

func (hh *HttpHandler) HttpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get("file-download"); len(vals) > 0 {

		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "file-download")

		w.Header().Set("Content-Disposition", md.HeaderMD.Get("Content-Disposition")[0])
		w.Header().Set("Content-Length", md.HeaderMD.Get("Content-Length")[0])

	}

	return nil

}
