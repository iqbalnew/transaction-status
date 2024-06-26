package db

import (
	"context"

	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/pb"
	"go.elastic.co/apm"
)

func (p *Provider) InsertJobTransactionPending(ctx context.Context, processId string, jobData *pb.JobTransactionStatusPending) (*pb.JobTransactionStatusPending, error) {
	span, _ := apm.StartSpan(ctx, "InsertJobTransactionPending", "db")
	span.Subtype = "postgresql"
	span.Action = "query"
	defer span.End()

	sqlQuery := `INSERT INTO transaction_pending.job_transaction_pending (task_id, status, created_at, updated_at)  
	VALUES ($1, $2, $3, $4) RETURNING id`

	sqlStmt, stmtErr := p.GetDbSql().SqlDb.Prepare(sqlQuery)
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer sqlStmt.Close()

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, p.GetDbSql().GetTimeout())
	defer ctxCancel()

	var lastInsertedId uint64

	execErr := sqlStmt.QueryRowContext(ctxTimeout,
		jobData.TaskId,
		jobData.Status.String(),
		"NOW()",
		"NOW()",
	).Scan(&lastInsertedId)
	if execErr != nil {
		return nil, execErr
	}

	jobData.Id = lastInsertedId

	return jobData, nil
}
