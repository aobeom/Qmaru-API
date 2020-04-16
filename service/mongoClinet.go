package service

import (
	"qmaru-api/config"
)

var dbCfg map[string]interface{} = config.DBCfg()
var dbHost string = dbCfg["dbhost"].(string)
var dbPort string = dbCfg["dbport"].(string)
var dbName string = dbCfg["dbname"].(string)

// DBTest 数据库测试
func DBTest() {
	MConnect(dbHost, dbPort)
}

// Client 数据库连接
var Client = MClient(dbHost, dbPort)

// DataBase 主数据库
var DataBase = Client.Database(dbName)
