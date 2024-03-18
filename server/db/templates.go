package db

import (
	"context"
	"database/sql"
	"slices"

	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/pb"
	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/utils"
	"go.elastic.co/apm"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (p *Provider) GetAllTemplates(ctx context.Context, pagination *pb.Pagination) ([]*pb.Templates, *pb.Pagination, error) {
	span, _ := apm.StartSpan(ctx, "GetAllTemplates", "db")
	span.Subtype = "postgresql"
	span.Action = "query"
	defer span.End()

	sqlQuery := `SELECT t.template_id, t.process_id, t.company_id, t.company_name, t.status_id, t.status_description, 
	t.raw_user_data, t.created_by, t.created_dt, t.updated_by, t.updated_dt, COUNT(1) OVER() AS total_rows 
	FROM "transaction".templates t`

	allowedOrderBy := []string{
		"template_id",
		"company_id",
		"company_name",
		"status_id",
	}

	sort := `template_id`
	if len(pagination.GetSort()) > 0 {
		if slices.Contains(allowedOrderBy, pagination.GetSort()) {
			sort = pagination.GetSort()
		}
	}

	sqlQuery = utils.GeneratePagination(sqlQuery, sort,
		pagination.GetDir(), pagination.GetLimit(), pagination.GetPage())

	sqlStmt, stmtErr := p.GetDbSql().SqlDb.Prepare(sqlQuery)
	if stmtErr != nil {
		return nil, nil, stmtErr
	}
	defer sqlStmt.Close()

	result := make([]*pb.Templates, 0)
	var totalData uint64 = 0

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, p.GetDbSql().GetTimeout())
	defer ctxCancel()

	dataRow, execErr := sqlStmt.QueryContext(ctxTimeout)
	if execErr != nil {
		return nil, nil, execErr
	}
	for dataRow.Next() {
		data := &pb.Templates{}
		var createdBy, updatedBy sql.NullString
		var createdDt, updatedDt sql.NullTime
		errScan := dataRow.Scan(
			&data.TemplateId,
			&data.ProcessId,
			&data.CompanyId,
			&data.CompanyName,
			&data.StatusId,
			&data.StatusDescription,
			&data.RawUserData,
			&createdBy,
			&createdDt,
			&updatedBy,
			&updatedDt,
			&totalData,
		)
		if errScan != nil {
			return nil, nil, errScan
		}

		data.CreatedBy = createdBy.String
		data.UpdatedBy = updatedBy.String
		data.CreatedDt = timestamppb.New(createdDt.Time)
		data.UpdatedDt = timestamppb.New(updatedDt.Time)

		result = append(result, data)
	}

	pagination.TotalPages = utils.CalculateTotalPage(uint64(pagination.GetLimit()), totalData)
	pagination.TotalRows = totalData

	return result, pagination, nil
}

func (p *Provider) GetTemplateById(ctx context.Context, templateId uint64) (*pb.Templates, error) {
	span, _ := apm.StartSpan(ctx, "GetTemplateById", "db")
	span.Subtype = "postgresql"
	span.Action = "query"
	defer span.End()

	sqlQuery := `SELECT t.template_id, t.process_id, t.company_id, t.company_name, t.status_id, 
	t.status_description, t.raw_user_data, t.created_by, t.created_dt, t.updated_by, t.updated_dt
	FROM "transaction".templates t
	WHERE t.template_id = $1`

	sqlStmt, stmtErr := p.GetDbSql().SqlDb.Prepare(sqlQuery)
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer sqlStmt.Close()

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, p.GetDbSql().GetTimeout())
	defer ctxCancel()

	result := &pb.Templates{}
	var createdBy, updatedBy sql.NullString
	var createdDt, updatedDt sql.NullTime
	errScan := sqlStmt.QueryRowContext(ctxTimeout, templateId).Scan(
		&result.TemplateId,
		&result.ProcessId,
		&result.CompanyId,
		&result.CompanyName,
		&result.StatusId,
		&result.StatusDescription,
		&result.RawUserData,
		&createdBy,
		&createdDt,
		&updatedBy,
		&updatedDt)
	if errScan != nil {
		if errScan == sql.ErrNoRows {
			return result, nil
		}
		return nil, errScan
	}

	result.CreatedBy = createdBy.String
	result.UpdatedBy = updatedBy.String
	result.CreatedDt = timestamppb.New(createdDt.Time)
	result.UpdatedDt = timestamppb.New(updatedDt.Time)

	return result, nil
}

func (p *Provider) InsertTemplate(ctx context.Context, processId string, templateData *pb.Templates) (*pb.Templates, error) {
	span, _ := apm.StartSpan(ctx, "InsertTemplate", "db")
	span.Subtype = "postgresql"
	span.Action = "query"
	defer span.End()

	sqlQuery := `INSERT INTO "transaction".templates (process_id, company_id, company_name, status_id, status_description, raw_user_data, created_by, created_dt) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING template_id`

	sqlStmt, stmtErr := p.GetDbSql().SqlDb.Prepare(sqlQuery)
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer sqlStmt.Close()

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, p.GetDbSql().GetTimeout())
	defer ctxCancel()

	var lastInsertedId uint64
	tempNilData := sql.NullString{
		Valid:  false,
		String: "",
	}
	var rawUserData interface{}
	if len(templateData.GetRawUserData()) <= 0 {
		rawUserData = tempNilData
	} else {
		rawUserData = templateData.GetRawUserData()
	}

	execErr := sqlStmt.QueryRowContext(ctxTimeout,
		processId,
		templateData.CompanyId,
		templateData.CompanyName,
		templateData.StatusId,
		templateData.StatusDescription,
		rawUserData,
		templateData.CreatedBy,
		"NOW()",
	).Scan(&lastInsertedId)
	if execErr != nil {
		return nil, execErr
	}

	templateData.TemplateId = lastInsertedId

	return templateData, nil
}

func (p *Provider) UpdateTemplate(ctx context.Context, templateData *pb.Templates) (*pb.Templates, int64, error) {
	span, _ := apm.StartSpan(ctx, "UpdateTemplate", "db")
	span.Subtype = "postgresql"
	span.Action = "query"
	defer span.End()

	sqlQuery := `UPDATE "transaction".templates SET status_id = $1, status_description = $2, updated_by = $3, updated_dt = $4 
	WHERE template_id = $5`

	sqlStmt, stmtErr := p.GetDbSql().SqlDb.Prepare(sqlQuery)
	if stmtErr != nil {
		return nil, 0, stmtErr
	}
	defer sqlStmt.Close()

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, p.GetDbSql().GetTimeout())
	defer ctxCancel()

	result, execErr := sqlStmt.ExecContext(ctxTimeout,
		templateData.StatusId,
		templateData.StatusDescription,
		templateData.UpdatedBy,
		"NOW()",
		templateData.GetTemplateId(),
	)
	if execErr != nil {
		return nil, 0, execErr
	}

	rowAffected, _ := result.RowsAffected()

	return templateData, rowAffected, nil
}

func (p *Provider) DeleteTemplate(ctx context.Context, templateId uint64) (int64, error) {
	span, _ := apm.StartSpan(ctx, "DeleteTemplate", "db")
	span.Subtype = "postgresql"
	span.Action = "query"
	defer span.End()

	sqlQuery := `DELETE FROM "transaction".templates t WHERE t.template_id = $1`

	sqlStmt, stmtErr := p.GetDbSql().SqlDb.Prepare(sqlQuery)
	if stmtErr != nil {
		return 0, stmtErr
	}
	defer sqlStmt.Close()

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, p.GetDbSql().GetTimeout())
	defer ctxCancel()

	result, execErr := sqlStmt.ExecContext(ctxTimeout,
		templateId,
	)
	if execErr != nil {
		return 0, execErr
	}

	rowAffected, _ := result.RowsAffected()

	return rowAffected, nil
}
