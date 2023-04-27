package server

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"tron/global"
)

type MongoLog struct {
	Tokens   string `json:"tokens" bson:"tokens"`
	BlockNum uint64 `json:"block_num" bson:"block_num"`
	IsError  bool   `json:"is_error" bson:"is_error"`
	Time     string `json:"time" bson:"time"`
	Err      string `json:"err" bson:"err"`
}

// WriteLog 录入块高日志
func WriteLog(blockNum uint64, err error) {
	collection := global.TRON_MongoDB.Database(global.TRON_CONFIG.MongoDB.Db).Collection("block_log")
	var log MongoLog
	log.Tokens = global.TRON_CONFIG.Tron.Tokens
	log.BlockNum = blockNum
	log.IsError = err != nil
	log.Time = time.Now().Format("2006-01-02 15:04:05")
	if err != nil {
		log.Err = err.Error()
	}
	filter := bson.M{"block_num": blockNum}
	update := bson.M{"$set": log}
	_, err = collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		global.TRON_LOG.Error(fmt.Sprint(err))
	}
}

// ReadEndBlockNum
func ReadEndBlockNum() uint64 {
	collection := global.TRON_MongoDB.Database(global.TRON_CONFIG.MongoDB.Db).Collection("block_log")
	option := options.FindOne()
	filter := bson.M{}
	option.SetSort(bson.D{{"block_num", -1}}) // Sort by block_num in descending order
	var result MongoLog
	err := collection.FindOne(context.Background(), filter, option).Decode(&result)
	if err != nil {
		if err.Error() == fmt.Sprintf("mongo: no documents in result") {
			return 0
		} else {
			panic(err)
		}
	}
	return result.BlockNum
}

func ReadStartBlockNum() uint64 {
	collection := global.TRON_MongoDB.Database(global.TRON_CONFIG.MongoDB.Db).Collection("block_log")
	option := options.FindOne()
	filter := bson.M{}
	option.SetSort(bson.D{{"block_num", 1}}) // Sort by block_num in descending order
	var result MongoLog
	err := collection.FindOne(context.Background(), filter, option).Decode(&result)
	if err != nil {
		if err.Error() == fmt.Sprintf("mongo: no documents in result") {
			return 0
		} else {
			panic(err)
		}
	}

	return result.BlockNum
}

// CheckMissingBlocks 检查缺失块高并且IsError为true的区块
func CheckMissingBlocks() []map[uint64]string {
	missBlockNum := make([]map[uint64]string, 0)
	endBlockNum := ReadEndBlockNum()
	startBlockNum := ReadStartBlockNum()
	for blockNum := startBlockNum; blockNum <= endBlockNum; blockNum++ {
		collection := global.TRON_MongoDB.Database(global.TRON_CONFIG.MongoDB.Db).Collection("block_log")
		option := options.FindOne()
		filter := bson.M{"block_num": blockNum}
		var result MongoLog
		err := collection.FindOne(context.Background(), filter, option).Decode(&result)
		if err != nil {
			if err.Error() == fmt.Sprintf("mongo: no documents in result") {
				missBlockNum = append(missBlockNum, map[uint64]string{blockNum: "missing block"})
			} else {
				panic(err)
			}
		}
		if result.IsError == true {
			missBlockNum = append(missBlockNum, map[uint64]string{blockNum: result.Err})
		}
	}
	return missBlockNum
}
