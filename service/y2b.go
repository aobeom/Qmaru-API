package service

import (
	"log"
	"os/exec"
	"qmaru-api/config"
	"strings"
)

var extapiCfg = config.AwsCfg()
var y2bcfg = extapiCfg["youtube"].(map[string]interface{})
var y2bexec = y2bcfg["exec"].(string)
var y2bfile = y2bcfg["file"].(string)

// Y2BDownload 下载 Y2B 仅最佳匹配
func Y2BDownload(url string) string {
	cmd := exec.Command(y2bexec, y2bfile, url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("cmd.Run() failed", err)
	}
	results := string(out)
	return strings.TrimSpace(results)
}
