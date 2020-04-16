package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"net"
	"log"
	"time"

	"golang.org/x/net/proxy"
)

// UserAgent 全局 UA
var UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36"

// S5Proxy 设置 s5 代理
func S5Proxy(proxyURL string) (transport *http.Transport) {
	dialer, err := proxy.SOCKS5("tcp", proxyURL,
		nil,
		&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		},
	)
	if err != nil {
		log.Panic(" [S5 Proxy Error]: ", err)
	}
	transport = &http.Transport{
		Proxy:               nil,
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return
}

// HTTPClient 设置 http 请求
func HTTPClient(proxy string) (client http.Client) {
	client = http.Client{Timeout: 30 * time.Second}
	if proxy != "" {
		transport := S5Proxy(proxy)
		client = http.Client{Timeout: 30 * time.Second, Transport: transport}
	}
	return
}

// YahooGet HTTP for Yahoo
func YahooGet(url string, reader io.Reader) (resURL string, resBody []byte) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	postReq, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Panic(" [YahooPost - Request Error]: ", err)
	}
	postReq.Header.Set("User-Agent", UserAgent)
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postRes, err := httpClient.Do(postReq)
	if err != nil {
		log.Panic(" [YahooPost - Response Error]: ", err)
	}
	location, _ := postRes.Location()
	cookies := postRes.Cookies()

	realURL := location.String()
	getReq, err := http.NewRequest("GET", realURL, nil)
	if err != nil {
		log.Panic(" [YahooGet - Request Error]: ", err)
	}
	getReq.Header.Set("User-Agent", UserAgent)
	for _, cookie := range cookies {
		getReq.AddCookie(cookie)
	}
	getRes, err := httpClient.Do(getReq)
	if err != nil {
		log.Panic(" [YahooGet - Response Error]: ", err)
	}
	body, err := ioutil.ReadAll(getRes.Body)
	resURL = realURL
	resBody = body
	return
}
