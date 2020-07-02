package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Example:
// filter := bson.D{{Key: "name", Value: "Ash"}}
// argStage := mongo.Pipeline{
//     bson.D{
//         {Key: "$sort", Value: bson.D{
//             {Key: "city", Value: 1},
//         }},
//     },
//     bson.D{
//         {Key: "$limit", Value: 1},
//     },
// }

// MConnect 数据库连接
func MConnect(host string, port string) {
	if host == "" || port == "" {
		log.Fatal("Database Address Error")
	}
	address := host + ":" + port
	log.Println("Database: " + address)

	client := MClient(host, port)

	err := client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Database connection failed")
	}

	log.Println("Database is OK")
}

// MClient 创建连接
func MClient(host string, port string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	address := host + ":" + port

	clientOptions := options.Client().ApplyURI("mongodb://" + address)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panic(err)
	}
	return client
}

// MInsertOne 插入一条数据
func MInsertOne(collection *mongo.Collection, data interface{}) interface{} {
	insertRes, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Panic(err)
	}
	return insertRes.InsertedID
}

// MInserMany 批量插入数据
func MInserMany(collection *mongo.Collection, data []interface{}) []interface{} {
	insertManyResult, err := collection.InsertMany(context.TODO(), data)
	if err != nil {
		log.Panic(err)
	}
	return insertManyResult.InsertedIDs
}

// MUpdate 更新数据
func MUpdate(collection *mongo.Collection, filter bson.D, update bson.D) (matched int64, modified int64) {
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Panic(err)
	}
	return updateResult.MatchedCount, updateResult.ModifiedCount
}

// MFind 查找数据
func MFind(collection *mongo.Collection, limit int64, skip int64, filter bson.D) (results []map[string]interface{}) {
	findOptions := options.Find()
	findOptions.SetLimit(skip)
	findOptions.SetLimit(limit)
	findOptions.SetProjection(bson.D{{Key: "_id", Value: 0}})
	if filter == nil {
		filter = bson.D{{}}
	}
	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Panic(err)
	}

	for cur.Next(context.TODO()) {
		var bsonData bson.D
		var jsonData map[string]interface{}
		var tmpB []byte
		_ = cur.Decode(&bsonData)
		tmpB, _ = bson.MarshalExtJSON(bsonData, true, true)
		_ = json.Unmarshal(tmpB, &jsonData)
		results = append(results, jsonData)
	}

	if err := cur.Err(); err != nil {
		log.Panic(err)
	}
	_ = cur.Close(context.TODO())
	return
}

// MDelete 删除数据
func MDelete(collection *mongo.Collection, filter bson.D) int64 {
	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Panic(err)
	}
	return deleteResult.DeletedCount
}

// MAggregate 聚合查询
func MAggregate(collection *mongo.Collection, argStage []bson.D) []map[string]interface{} {
	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline(argStage))
	if err != nil {
		log.Panic(err)
	}
	var results []map[string]interface{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Panic(err)
	}
	return results
}
