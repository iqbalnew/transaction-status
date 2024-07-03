package apigrpc

import (
	"context"

	pb "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/pb"
	"github.com/sirupsen/logrus"
)

func (s *Server) RegisterJobTransactionPending(ctx context.Context, req *pb.RegisterJobTransactionPendingeRequest) (*pb.RegisterJobTransactionPendingResponse, error) {

	result, _, err := s.provider.GetAllTransactionPending(
		ctx,
		&pb.Pagination{
			Limit: 100,
			Page:  1,
			Sort:  "updated_at",
			Dir:   pb.Direction_ASC,
		},
	)
	if err != nil {
		logrus.Errorln("failedGet Data")
		return nil, err
	}

	for i, v := range result {

		s.logger.Infoln(i, " loop")

		job := &pb.JobTransactionStatusPending{
			TaskId: v.TaskId,
			Status: pb.StatusInquiryJob_NEW,
			Type:   v.GetType(),
		}

		_, err := s.provider.InsertJobTransactionPending(ctx, "0", job)

		if err == nil {
			s.logger.Println("Success insert data")
		} else {
			s.logger.Infoln(i, " error: ", err)
		}

	}

	// just insert the data
	return &pb.RegisterJobTransactionPendingResponse{
		Error:   false,
		Code:    200,
		Message: "Success Register New Job Pending",
	}, nil
}

func (s *Server) UpdatedJobTransactionStatus(ctx context.Context, req *pb.UpdatedJobTransactionStatusRequest) (*pb.UpdatedJobTransactionStatusResponse, error) {

	data := &pb.JobTransactionStatusPending{
		Id:     req.GetIdJob(),
		Status: req.GetStatus(),
	}

	_, err := s.provider.UpdateJobTransactionPending(context.Background(), data)
	if err != nil {
		s.logger.Errorln("failed updated UpdateJobTransactionPending error :", err)
		return nil, err
	}

	return &pb.UpdatedJobTransactionStatusResponse{
		Error:   false,
		Code:    200,
		Message: "Success Register New Job Pending",
	}, nil
}
