package globals

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Db             *mongo.Client
	UserDb         *mongo.Collection
	MongoCtx       context.Context
	mongoUsername  = "user-service"
	mongoPwd       = "vLxxhmS0eJFwmteF"
	ConnURI        = "mongodb+srv://" + mongoUsername + ":" + mongoPwd + "@cluster0.fpedw5d.mongodb.net/"
	DbName         = "UserService"
	CollectionName = "Users"
)
