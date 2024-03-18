package db

import (
	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/database"
)

type Provider struct {
	dbSql *database.DbSql
}

func NewProvider(dbSql *database.DbSql) *Provider {
	return &Provider{
		dbSql: dbSql,
	}
}

func (p *Provider) GetDbSql() *database.DbSql {
	return p.dbSql
}
