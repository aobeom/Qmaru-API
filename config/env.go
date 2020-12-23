package config

var mode int = 0

// Deployment 部署模式 0: debug 1: release
func Deployment() bool {
	if mode == 0 {
		return true
	}
	return false
}
