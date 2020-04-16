package service

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// getDateRange 截取当前月份
func getDateRange() (mint, maxt string) {
	year := time.Now().Year()
	month := time.Now().Month()

	var newMonth string
	if month < 10 {
		newMonth = fmt.Sprintf("%0*d", 2, month)
	} else {
		newMonth = fmt.Sprintf("%d", month)
	}
	mint = fmt.Sprintf("%d-%s-00", year, newMonth)
	maxt = fmt.Sprintf("%d-%s-31", year, newMonth)
	return
}

// DaramaData 读取 Drama 的数据
func DaramaData(dtype string) (data []map[string]interface{}) {
	dramaColl := DataBase.Collection("drama_info")

	start, end := getDateRange()
	// fixsub 没有更新日期 直接倒叙返回前 15 条数据
	if dtype == "fixsub" {
		var argStage = []bson.D{
			{
				{Key: "$match", Value: bson.D{
					{Key: "type", Value: dtype},
				}},
			},
			{
				{Key: "$sort", Value: bson.D{
					{Key: "_id", Value: -1},
				}},
			},
			{
				{Key: "$limit", Value: 15},
			},
			{
				{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "createdat", Value: 0},
					{Key: "updatedat", Value: 0},
				}},
			},
		}
		data = MAggregate(dramaColl, argStage)
		// 其他返回当前月份的数据
	} else {
		var argStage = []bson.D{
			{
				{Key: "$match", Value: bson.D{
					{Key: "type", Value: dtype},
					{Key: "date", Value: bson.D{
						{Key: "$lt", Value: end},
						{Key: "$gt", Value: start},
					}},
				}},
			},
			{
				{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "createdat", Value: 0},
					{Key: "updatedat", Value: 0},
				}},
			},
			{
				{Key: "$sort", Value: bson.D{
					{Key: "date", Value: -1},
				}},
			},
		}
		data = MAggregate(dramaColl, argStage)
	}
	return
}
