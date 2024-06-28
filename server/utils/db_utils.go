package utils

import (
	"fmt"

	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/pb"
)

func GeneratePagination(sqlQuery string, sortBy string, direction pb.Direction, limit int32, page int32) string {
	switch direction {
	case pb.Direction_DESC:
		sqlQuery = sqlQuery + " ORDER BY " + sortBy + " DESC"
	default:
		sqlQuery = sqlQuery + "  ORDER BY " + sortBy + " ASC"
	}

	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	sqlQuery = sqlQuery + fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	return sqlQuery
}
