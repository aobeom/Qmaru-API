package config

var mode int = 1

// Deployment 部署模式 0: debug 1: release
func Deployment() bool {
	if mode == 0 {
		return true
	}
	return false
}
