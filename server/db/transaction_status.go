package db

import (
	"context"
	"database/sql"
	"log"
	"slices"

	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/pb"
	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/utils"
	"go.elastic.co/apm"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (p *Provider) GetAllTransactionPending(ctx context.Context, pagination *pb.Pagination) ([]*pb.TransactionPending, *pb.Pagination, error) {
	span, _ := apm.StartSpan(ctx, "GetAllTransactionPending", "db")
	span.Subtype = "postgresql"
	span.Action = "query"
	defer span.End()

	sqlQuery := `SELECT t.type, t.updated_at, t.id, t.task_id, t.status
			FROM public.transaction_today_pending_no_job t`

	allowedOrderBy := []string{
		"id",
		"task_id",
		"type",
		"status",
		"updated_at",
	}

	sort := `id`
	if len(pagination.GetSort()) > 0 {
		if slices.Contains(allowedOrderBy, pagination.GetSort()) {
			sort = pagination.GetSort()
		}
	}

	sqlQuery = utils.GeneratePagination(sqlQuery, sort,
		pagination.GetDir(), pagination.GetLimit(), pagination.GetPage())

	log.Println(sqlQuery)

	sqlStmt, stmtErr := p.GetDbSql().SqlDb.Prepare(sqlQuery)
	if stmtErr != nil {
		return nil, nil, stmtErr
	}
	defer sqlStmt.Close()

	result := make([]*pb.TransactionPending, 0)
	var totalData uint64 = 0

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, p.GetDbSql().GetTimeout())
	defer ctxCancel()

	dataRow, execErr := sqlStmt.QueryContext(ctxTimeout)
	if execErr != nil {
		return nil, nil, execErr
	}
	for dataRow.Next() {
		data := &pb.TransactionPending{}
		var updatedAt sql.NullTime
		errScan := dataRow.Scan(
			&data.Type,
			&updatedAt,
			&data.Id,
			&data.TaskId,
			&data.Status,
		)
		if errScan != nil {
			return nil, nil, errScan
		}

		data.UpdatedAt = timestamppb.New(updatedAt.Time)

		result = append(result, data)
	}

	pagination.TotalPages = utils.CalculateTotalPage(uint64(pagination.GetLimit()), totalData)
	pagination.TotalRows = totalData

	return result, pagination, nil
}
