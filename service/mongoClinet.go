package service

import (
	"qmaru-api/config"
)

var (
	dbCfg  = config.DBCfg()
	dbHost = dbCfg["dbhost"].(string)
	dbPort = dbCfg["dbport"].(string)
	dbName = dbCfg["dbname"].(string)
)

// DBTest 数据库测试
func DBTest() {
	MConnect(dbHost, dbPort)
}

// Client 数据库连接
var Client = MClient(dbHost, dbPort)

// DataBase 主数据库
var DataBase = Client.Database(dbName)
