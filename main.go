package main

import (
	"qmaru-api/apis"
	"qmaru-api/config"
)

func main() {
	apis.Run(config.Deployment())
}
