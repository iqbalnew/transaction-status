package mongodb

import (
	"context"
	"errors"
	"testing"

	mongodbmock "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/mongodb/mock"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *MongoDbTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func TestInitMongoDb(t *testing.T) {
	suite.Run(t, new(MongoDbTestSuite))
}

func (s *MongoDbTestSuite) TestMongoDb_New() {
	type expectation struct {
		out *MongoDb
	}

	tests := map[string]struct {
		config   Config
		expected expectation
	}{
		"SuccessCreateNewMongoDb": {
			config: Config{},
			expected: expectation{
				out: &MongoDb{},
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			out := New(s.ctx, tt.config, &mongodbmock.MongoDbMock{})

			if out == nil {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, out)
			}
		})
	}
}

func (s *MongoDbTestSuite) TestMongoDb_GetMongoDatabaseName() {
	type expectation struct {
		out string
	}

	tests := map[string]struct {
		mongoDb  *MongoDb
		expected expectation
	}{
		"SuccessGetMongoDatabaseName": {
			mongoDb: &MongoDb{
				Config: Config{
					DatabaseName: "unit-test-mongo",
				},
			},
			expected: expectation{
				out: "unit-test-mongo",
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			mongoDatabase := tt.mongoDb.GetMongoDatabaseName()

			if mongoDatabase != tt.expected.out {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, mongoDatabase)
			}
		})
	}
}

func (s *MongoDbTestSuite) TestMongoDb_GetMongoCollectionName() {
	type expectation struct {
		out string
	}

	tests := map[string]struct {
		mongoDb  *MongoDb
		expected expectation
	}{
		"SuccessGetMongoCollectionName": {
			mongoDb: &MongoDb{
				Config: Config{
					CollectionName: "unit-test-mongo",
				},
			},
			expected: expectation{
				out: "unit-test-mongo",
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			mongoCollection := tt.mongoDb.GetMongoCollectionName()

			if mongoCollection != tt.expected.out {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, mongoCollection)
			}
		})
	}
}

func (s *MongoDbTestSuite) TestMongoDb_Connect() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		mongoDb  *MongoDb
		expected expectation
	}{
		"Success": {
			mongoDb: &MongoDb{
				Config: Config{
					Uri: "mongodb://test:test@localhost:27012",
				},
				mdi: &mongodbmock.MongoDbMock{},
			},
			expected: expectation{
				err: nil,
			},
		},
		"Failed": {
			mongoDb: &MongoDb{
				mdi: &mongodbmock.MongoDbMock{},
			},
			expected: expectation{
				err: errors.New("failed connect to mongodb server"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			credential := options.Credential{
				Username: "",
				Password: "",
			}
			err := tt.mongoDb.Connect(credential)

			if err != nil {
				if err.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *MongoDbTestSuite) TestMongoDb_Disconnect() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		mongoDb  *MongoDb
		expected expectation
	}{
		"Success": {
			mongoDb: &MongoDb{
				ctx: s.ctx,
				mdi: &mongodbmock.MongoDbMock{},
			},
			expected: expectation{
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			err := tt.mongoDb.Close()

			if err != nil {
				if err.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}
