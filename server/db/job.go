package db

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/pb"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/utils"
	"go.elastic.co/apm"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (p *Provider) InsertJobTransactionPending(ctx context.Context, processId string, jobData *pb.JobTransactionStatusPending) (*pb.JobTransactionStatusPending, error) {
	span, _ := apm.StartSpan(ctx, "InsertJobTransactionPending", "db")
	span.Subtype = "postgresql"
	span.Action = "query"
	defer span.End()

	sqlQuery := `INSERT INTO transaction_pending.job_transaction_pending (task_id, status, created_at, updated_at, type)  
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	sqlStmt, stmtErr := p.GetDbSql().SqlDb.Prepare(sqlQuery)
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer sqlStmt.Close()

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, p.GetDbSql().GetTimeout())
	defer ctxCancel()

	var lastInsertedId uint64

	execErr := sqlStmt.QueryRowContext(ctxTimeout,
		jobData.GetTaskId(),
		jobData.GetStatus().String(),
		"NOW()",
		"NOW()",
		jobData.GetType(),
	).Scan(&lastInsertedId)
	if execErr != nil {
		return nil, execErr
	}

	jobData.Id = lastInsertedId

	return jobData, nil
}

func (p *Provider) UpdateJobTransactionPending(ctx context.Context, jobData *pb.JobTransactionStatusPending) (*pb.JobTransactionStatusPending, error) {
	span, _ := apm.StartSpan(ctx, "UpdateJobTransactionPending", "db")
	span.Subtype = "postgresql"
	span.Action = "query"
	defer span.End()

	sqlQuery := `UPDATE transaction_pending.job_transaction_pending 
                 SET task_id = $1, status = $2, updated_at = NOW() 
                 WHERE id = $3`

	sqlStmt, stmtErr := p.GetDbSql().SqlDb.Prepare(sqlQuery)
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer sqlStmt.Close()

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, p.GetDbSql().GetTimeout())
	defer ctxCancel()

	_, execErr := sqlStmt.ExecContext(ctxTimeout,
		jobData.TaskId,
		jobData.Status.String(),
		jobData.Id,
	)
	if execErr != nil {
		return nil, execErr
	}

	return jobData, nil
}

func (p *Provider) GetJobTransactionPendingById(ctx context.Context, id uint64) (*pb.JobTransactionStatusPending, error) {
	span, _ := apm.StartSpan(ctx, "GetJobTransactionPendingById", "db")
	span.Subtype = "postgresql"
	span.Action = "query"
	defer span.End()

	sqlQuery := `SELECT id, task_id, type,  status, created_at, updated_at 
                 FROM transaction_pending.job_transaction_pending 
                 WHERE id = $1`

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, p.GetDbSql().GetTimeout())
	defer ctxCancel()

	row := p.GetDbSql().SqlDb.QueryRowContext(ctxTimeout, sqlQuery, id)

	var jobData pb.JobTransactionStatusPending
	var status string

	var createdAt, updatedAt sql.NullTime

	err := row.Scan(
		&jobData.Id,
		&jobData.TaskId,
		&jobData.Type,
		&status,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	jobData.CreatedAt = timestamppb.New(createdAt.Time)
	jobData.UpdatedAt = timestamppb.New(updatedAt.Time)

	jobData.Status = pb.StatusInquiryJob(pb.StatusInquiryJob_value[status])

	return &jobData, nil
}

func (p *Provider) GetAllJobTransactionsPending(ctx context.Context, pagination *pb.Pagination) ([]*pb.JobTransactionStatusPending, error) {
	span, _ := apm.StartSpan(ctx, "GetAllJobTransactionsPending", "db")
	span.Subtype = "postgresql"
	span.Action = "query"
	defer span.End()

	sqlQuery := `SELECT id, task_id, type, status, created_at, updated_at 
                 FROM transaction_pending.job_transaction_pending`

	allowedOrderBy := []string{
		"id",
		"task_id",
		"type",
		"status",
		"updated_at",
		"created_at",
	}

	if pagination.GetFilter() != "" {
		sqlQuery += " WHERE " + pagination.GetFilter()
	}

	sort := `id`
	if len(pagination.GetSort()) > 0 {
		if slices.Contains(allowedOrderBy, pagination.GetSort()) {
			sort = pagination.GetSort()
		}
	}

	sqlQuery = utils.GeneratePagination(sqlQuery, sort,
		pagination.GetDir(), pagination.GetLimit(), pagination.GetPage())

	sqlStmt, stmtErr := p.GetDbSql().SqlDb.Prepare(sqlQuery)
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer sqlStmt.Close()

	var jobTransactions []*pb.JobTransactionStatusPending

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, p.GetDbSql().GetTimeout())
	defer ctxCancel()

	rows, execErr := sqlStmt.QueryContext(ctxTimeout)
	if execErr != nil {
		return nil, execErr
	}

	for rows.Next() {
		var jobData pb.JobTransactionStatusPending
		var status string

		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&jobData.Id,
			&jobData.TaskId,
			&jobData.Type,
			&status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		jobData.CreatedAt = timestamppb.New(createdAt.Time)
		jobData.UpdatedAt = timestamppb.New(updatedAt.Time)

		jobData.Status = pb.StatusInquiryJob(pb.StatusInquiryJob_value[status])
		jobTransactions = append(jobTransactions, &jobData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobTransactions, nil
}
