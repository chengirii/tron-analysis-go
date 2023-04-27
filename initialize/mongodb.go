package initialize

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"tron/global"
)

func MongoDB() *mongo.Client {
	mongodbCfg := global.TRON_CONFIG.MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongodbCfg.Dns()))
	if err != nil {
		panic(fmt.Errorf("connet fail check mongodb config : %s\n", err))
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(fmt.Errorf("mongodb connect ping failed : %s\n", err))
	}
	checkDatabaseExists(client, mongodbCfg.Db)
	return client

}

// checkDatabaseExists 检查数据是否存在
func checkDatabaseExists(client *mongo.Client, dbName string) {
	dbNames, err := client.ListDatabaseNames(context.Background(), bson.M{})
	if err != nil {
		return
	}
	for _, name := range dbNames {
		if name == dbName {
			return
		}
	}
	createDatabase(client, dbName)
}

// createDatabase 创建数据库
func createDatabase(client *mongo.Client, dbName string) {
	database := client.Database(dbName)
	err := database.CreateCollection(context.Background(), "block_log", nil)
	if err != nil {
		panic(fmt.Errorf("Failed to create :%s\n", err))
	}
	err = createIndex(client, "block_log", "block_num", true)
	if err != nil {
		panic(fmt.Errorf("Failed to create index :%s\n", err))
	}
}

// createIndex 创建索引
func createIndex(client *mongo.Client, collectionName string, fieldName string, isUnique bool) error {
	collection := client.Database(global.TRON_CONFIG.MongoDB.Db).Collection(collectionName)
	index := bsonx.Doc{{Key: fieldName, Value: bsonx.Int32(1)}}
	options := options.Index().SetUnique(isUnique)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: index, Options: options})
	return err
}
