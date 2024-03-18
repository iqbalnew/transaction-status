package mongodb

import (
	"context"
	"time"

	mongodbwrapper "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/mongodb/wrapper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Uri            string
	DatabaseName   string
	CollectionName string
	Timeout        int64
}

type MongoDb struct {
	ctx         context.Context
	MongoClient *mongo.Client
	Config      Config
	mdi         mongodbwrapper.MongoDbInterface
}

func New(ctx context.Context, config Config, mdi mongodbwrapper.MongoDbInterface) *MongoDb {
	return &MongoDb{
		ctx:    ctx,
		Config: config,
		mdi:    mdi,
	}
}

func (md *MongoDb) GetMongoDatabaseName() string {
	return md.Config.DatabaseName
}

func (md *MongoDb) GetMongoCollectionName() string {
	return md.Config.CollectionName
}

func (md *MongoDb) Connect(auth options.Credential) error {
	timeout := time.Duration(md.Config.Timeout) * time.Second
	opts := options.Client().ApplyURI(md.Config.Uri).SetAuth(auth).SetTimeout(timeout)

	var mongoErr error
	md.MongoClient, mongoErr = md.mdi.Connect(md.ctx, opts)
	if mongoErr != nil {
		return mongoErr
	}

	return nil
}

func (md *MongoDb) Close() error {
	return md.mdi.Disconnect(md.ctx)
}
