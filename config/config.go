package config

import (
	"path/filepath"
	"qmaru-api/utils"
)

var cfgRoot = "config"
var fileCtl = utils.NewFileCtl()
var dataConv = utils.NewDataConv()

func readCfg(name string) (d map[string]interface{}) {
	cfgPath := filepath.Join(fileCtl.LocalPath(Deployment()), cfgRoot, name)
	d = dataConv.String2Map(fileCtl.Read(cfgPath))
	return
}

// DBCfg 数据库连接配置
func DBCfg() (d map[string]interface{}) {
	d = readCfg("database.json")
	return
}

// MediaCfg 静态文件配置
func MediaCfg() (d map[string]interface{}) {
	d = readCfg("media.json")
	return
}

// TweetCfg 推认证
func TweetCfg() (d map[string]interface{}) {
	d = readCfg("tweet.json")
	return
}

// AwsCfg aws 网关
func AwsCfg() (d map[string]interface{}) {
	d = readCfg("extapi.json")
	return
}
