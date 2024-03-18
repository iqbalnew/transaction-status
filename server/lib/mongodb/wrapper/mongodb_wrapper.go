package mongodbwrapper

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbWrapper struct {
	mongoDbClient *mongo.Client
	MongoDbInterface
}

func (mdw *MongoDbWrapper) Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	return mongo.Connect(ctx, opts...)
}

func (mdw *MongoDbWrapper) Disconnect(ctx context.Context) error {
	return mdw.mongoDbClient.Disconnect(ctx)
}
