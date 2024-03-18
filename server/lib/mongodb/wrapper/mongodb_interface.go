package mongodbwrapper

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbInterface interface {
	Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error)
	Disconnect(ctx context.Context) error
}
