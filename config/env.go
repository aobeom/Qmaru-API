package config

var mode int = 0

// Deployment 部署模式 0: debug 1: release
func Deployment() int {
	return mode
}
