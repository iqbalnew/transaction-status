package apigrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	pb "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/pb"
	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/utils"
	"go.elastic.co/apm/v2"
)

func (s *Server) GetAllTemplates(ctx context.Context, req *pb.GetAllTemplatesRequest) (*pb.GetAllTemplatesResponse, error) {
	span, _ := apm.StartSpan(ctx, "GetAllTemplates", "process")
	span.Action = "execute"
	defer span.End()

	logger := s.logger.WithTaskID(req.GetProcessId())
	logger.Debug("[GetAllTemplates] Start processing request")
	logger.WithField(s.logger.GrpcMetadataKey, req.String()).Info("[GetAllTemplates] Receive new request")

	templateList, pagination, getAllErr := s.provider.GetAllTemplates(ctx, req.GetPagination())
	if getAllErr != nil {
		return nil, getAllErr
	}

	return &pb.GetAllTemplatesResponse{
		Error:      false,
		Code:       http.StatusOK,
		Message:    "Get template data success",
		Data:       templateList,
		Pagination: pagination,
	}, nil
}

func (s *Server) GetTemplateDetail(ctx context.Context, req *pb.GetTemplateDetailRequest) (*pb.GetTemplateDetailResponse, error) {
	span, _ := apm.StartSpan(ctx, "GetTemplateDetail", "process")
	span.Action = "execute"
	defer span.End()

	logger := s.logger.WithTaskID(req.GetProcessId())
	logger.Debug("[GetTemplateDetail] Start processing request")
	logger.WithField(s.logger.GrpcMetadataKey, req.String()).Info("[GetTemplateDetail] Receive new request")

	templateData, getTemplateErr := s.provider.GetTemplateById(ctx, req.GetTemplateId())
	if getTemplateErr != nil {
		return nil, getTemplateErr
	}

	result := &pb.GetTemplateDetailResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Get template data success",
		Data:    templateData,
	}

	if templateData.GetTemplateId() <= 0 {
		result.Error = true
		result.Message = "Template data not found"
		result.Code = http.StatusNotFound
		result.Data = nil
	}

	return result, nil
}

func (s *Server) SaveTemplate(ctx context.Context, req *pb.SaveTemplateRequest) (*pb.GeneralBodyResponse, error) {
	span, _ := apm.StartSpan(ctx, "SaveTemplate", "process")
	span.Action = "execute"
	defer span.End()

	logger := s.logger.WithTaskID(req.GetProcessId())
	logger.Debug("[SaveTemplate] Start processing request")
	logger.WithField(s.logger.GrpcMetadataKey, req.String()).Info("[SaveTemplate] Receive new request")

	logger.Debug("[BatchFileUploadCheckDuplicate] Start get user from context")
	userData, _ := utils.GetUserFromContext(ctx)
	logger.Debug("[BatchFileUploadCheckDuplicate] Finish get user from context")

	userDataJson, marshalErr := json.Marshal(userData)
	if marshalErr != nil {
		return nil, marshalErr
	}

	req.GetTemplate().RawUserData = userDataJson

	templateData, getTemplateErr := s.provider.InsertTemplate(ctx, req.GetProcessId(), req.GetTemplate())
	if getTemplateErr != nil {
		return nil, getTemplateErr
	}
	logger.WithField(s.logger.GrpcMetadataKey, templateData.String()).Info("[SaveTemplate] Insert template success")

	return &pb.GeneralBodyResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Save template data success",
	}, nil
}

func (s *Server) UpdateTemplate(ctx context.Context, req *pb.UpdateTemplateRequest) (*pb.GetTemplateDetailResponse, error) {
	span, _ := apm.StartSpan(ctx, "UpdateTemplate", "process")
	span.Action = "execute"
	defer span.End()

	logger := s.logger.WithTaskID(req.GetProcessId())
	logger.Debug("[UpdateTemplate] Start processing request")
	logger.WithField(s.logger.GrpcMetadataKey, req.String()).Info("[UpdateTemplate] Receive new request")

	updatedData, rowsAffected, updateTemplateErr := s.provider.UpdateTemplate(ctx, req.GetTemplate())
	if updateTemplateErr != nil {
		return nil, updateTemplateErr
	}
	logger.WithField(s.logger.GrpcMetadataKey, fmt.Sprintf("updated_data: %v, rows_affected: %d", updatedData, rowsAffected)).Info("[UpdateTemplate] Update template success")

	return &pb.GetTemplateDetailResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Update template data success, total_affected: %d", rowsAffected),
		Data:    updatedData,
	}, nil
}

func (s *Server) DeleteTemplate(ctx context.Context, req *pb.DeleteTemplateRequest) (*pb.GeneralBodyResponse, error) {
	span, _ := apm.StartSpan(ctx, "DeleteTemplate", "process")
	span.Action = "execute"
	defer span.End()

	logger := s.logger.WithTaskID(req.GetProcessId())
	logger.Debug("[DeleteTemplate] Start processing request")
	logger.WithField(s.logger.GrpcMetadataKey, req.String()).Info("[DeleteTemplate] Receive new request")

	rowsAffected, updateTemplateErr := s.provider.DeleteTemplate(ctx, req.GetTemplateId())
	if updateTemplateErr != nil {
		return nil, updateTemplateErr
	}
	logger.WithField(s.logger.GrpcMetadataKey, fmt.Sprintf("rows_affected: %d", rowsAffected)).Info("[DeleteTemplate] Delete template success")

	return &pb.GeneralBodyResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Delete template data success, total_affected: %d", rowsAffected),
	}, nil
}
