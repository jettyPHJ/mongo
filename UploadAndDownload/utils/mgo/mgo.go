package mgo

import (
	cfg "UploadAndDownload/utils/config"
	"UploadAndDownload/utils/log"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Database *mongo.Database

func init() {
	clientOptions := options.Client().ApplyURI(cfg.Param.Mgo.Uri)
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	Database = client.Database("base1")
}
