package globals

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Db          *mongo.Client
	UserDb      *mongo.Collection
	MongoCtx    context.Context
	MongoDBUrl  string
	RabbitMQUrl string
)
