package config

var mode string = "release"

// Deployment 部署模式
func Deployment() string {
	return mode
}
