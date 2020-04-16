package api

import (
	"qmaru-api/service"
	"testing"
)

func Test_Picdown(t *testing.T) {
	testSite := map[string]string{
		"mdpr.jp":           "https://mdpr.jp/news/detail/1763404",
		"ameblo.jp":         "https://ameblo.jp/sayaka-kanda/entry-12372153694.html",
		"thetv.jp":          "https://thetv.jp/news/detail/145669/",
		"tokyopopline.com":  "https://tokyopopline.com/archives/100688",
		"hustlepress.co.jp": "https://hustlepress.co.jp/shiraishi_20190105_interview/",
		"lineblog.me":       "https://lineblog.me/mamoru_miyano/archives/1744026.html",
		"instagram.com":     "https://www.instagram.com/p/Bwsi9mvBgBx/",
	}
	for site, url := range testSite {
		imgs := service.PicData(site, url)
		if len(imgs) != 0 {
			t.Log(site + " ok")
		} else {
			t.Error(site + " 规则改变")
		}
	}
}
