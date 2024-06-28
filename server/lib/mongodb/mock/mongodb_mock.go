package mongodbmock

import (
	"context"
	"errors"
	"testing"

	mongodbwrapper "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/mongodb/wrapper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbMock struct {
	mongodbwrapper.MongoDbInterface
}

func NewMock(t *testing.T) *mtest.T {
	return mtest.New(t, mtest.NewOptions().CreateClient(false).ClientType(mtest.Mock))
}

func (mdm *MongoDbMock) Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	if opts[0].GetURI() == "" {
		return nil, errors.New("failed connect to mongodb server")
	}

	return &mongo.Client{}, nil
}

func (mdm *MongoDbMock) Disconnect(ctx context.Context) error {
	return nil
}
