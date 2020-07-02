package service

import (
	"go.mongodb.org/mongo-driver/bson"
)

// STData 读取 STchannel 的最后 15 条数据
func STData() (data []map[string]interface{}) {
	stColl := DataBase.Collection("st_info")
	var argStage = []bson.D{
		{
			{Key: "$sort", Value: bson.D{
				{Key: "date", Value: -1},
			}},
		},
		{
			{Key: "$limit", Value: 15},
		},
		{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
			}},
		},
	}
	data = MAggregate(stColl, argStage)
	return
}