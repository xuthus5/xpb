package driver

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"pastebin/config"
	"time"
)

var (
	client *mongo.Client
	confer = config.GetConfig()
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var err error

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(confer.Mongo.URI).
		SetMaxPoolSize(16).SetMaxConnIdleTime(5*time.Second))
	if err != nil {
		panic(err)
	}
}
