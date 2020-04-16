package service

import (
	"go.mongodb.org/mongo-driver/bson"
)

// CronTime 获取定时任务执行的时间
func CronTime(cronType string) (cTime string) {
	cronColl := DataBase.Collection("crond_time")
	fData := bson.D{
		{Key: "type", Value: cronType},
	}
	cronResult := MFind(cronColl, 1, 0, fData)
	if len(cronResult) != 0 {
		cronData := cronResult[0]
		cTime = cronData["time"].(string)
	} else {
		cTime = ""
	}
	return
}
