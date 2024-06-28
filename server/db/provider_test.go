package db

import (
	"context"
	"database/sql"
	"testing"

	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type ProviderTestSuite struct {
	suite.Suite
	ctx     context.Context
	mockDb  *sql.DB
	sqlMock sqlmock.Sqlmock
}

func (s *ProviderTestSuite) SetupTest() {
	s.ctx = context.Background()
	var err error

	s.mockDb, s.sqlMock, err = sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
}

func TestInitProvider(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}

func (s *ProviderTestSuite) TestProvider_New() {
	type expectation struct {
		out *Provider
	}

	tests := map[string]struct {
		expected expectation
	}{
		"Success": {
			expected: expectation{
				out: &Provider{},
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			provider := NewProvider(nil)

			if *provider != *tt.expected.out {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, provider)
			}
		})
	}
}

func (s *ProviderTestSuite) TestProvider_GetDbSql() {
	type expectation struct {
		out *database.DbSql
	}

	tests := map[string]struct {
		expected expectation
	}{
		"Success": {
			expected: expectation{
				out: &database.DbSql{},
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			provider := &Provider{
				dbSql: &database.DbSql{},
			}
			dbSql := provider.GetDbSql()

			if *dbSql != *tt.expected.out {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, dbSql)
			}
		})
	}
}
