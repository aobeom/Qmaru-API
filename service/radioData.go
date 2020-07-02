package service

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"qmaru-api/config"
	"qmaru-api/utils"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

// Common Header 公共头
const (
	XRadikoDevice     = "pc"
	XRadikoUser       = "dummy_user"
	XRadikoApp        = "pc_html5"
	XRadikoAppVersion = "0.0.1"
)

var miniReq = utils.Minireq

// DLServer 下载控制器
type DLServer struct {
	WG    sync.WaitGroup
	Gonum chan string
}

// RadioData 请求参数
type RadioData struct {
	stationID string
	startAt   string
	endAt     string
	ft        string
	to        string
	l         string
	rtype     string
}

type radioJSON struct {
	Name string
	URL  string
}

// encodeKey 根据偏移长度生成 KEY
func encodeKey(authkey string, offset int64, length int64) (partialkey string) {
	reader := strings.NewReader(authkey)
	buff := make([]byte, length)
	_, err := reader.ReadAt(buff, offset)
	if err != nil {
		log.Panic(err)
	}
	partialkey = base64.StdEncoding.EncodeToString(buff)
	return
}

// radikoJSKey 提取 JS 的密钥
func radikoJSKey(client http.Client) (authkey string) {
	playerURL := "http://radiko.jp/apps/js/playerCommon.js"
	req, _ := http.NewRequest("GET", playerURL, nil)
	req.Header.Add("User-Agent", utils.UserAgent)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	regKeyRule := regexp.MustCompile(`[0-9a-z]{40}`)
	authkeyMap := regKeyRule.FindAllString(string(body), -1)
	authkey = authkeyMap[0]
	return
}

// radikoChunklist 提取播放地址
func radikoChunklist(playlist string) (url string) {
	regURLRule := regexp.MustCompile(`https://.*?\.m3u8`)
	urlList := regURLRule.FindAllString(playlist, -1)
	url = urlList[0]
	return
}

// radikoAAC 提取 AAC 文件的地址
func radikoAAC(m3u8 string) (urls []string) {
	regURLRule := regexp.MustCompile(`https://.*?\.aac`)
	urls = regURLRule.FindAllString(m3u8, -1)
	return
}

// radikoAuth1 获取 token / offset / length
func radikoAuth1() (token string, offset int64, length int64) {
	auth1URL := "https://radiko.jp/v2/api/auth1"
	headers := make(http.Header)
	headers.Add("User-Agent", utils.UserAgent)
	headers.Add("x-radiko-device", XRadikoDevice)
	headers.Add("x-radiko-user", XRadikoUser)
	headers.Add("x-radiko-app", XRadikoApp)
	headers.Add("x-radiko-app-version", XRadikoAppVersion)

	res := miniReq.GetRes(auth1URL, headers, nil)
	resHeader := res.Header
	token = resHeader.Get("X-Radiko-AuthToken")
	Keyoffset := resHeader.Get("X-Radiko-Keyoffset")
	Keylength := resHeader.Get("X-Radiko-Keylength")
	offset, _ = strconv.ParseInt(Keyoffset, 10, 64)
	length, _ = strconv.ParseInt(Keylength, 10, 64)
	return
}

// radikoAuth2 获取地区代码
func radikoAuth2(token string, partialkey string) (areaid string) {
	auth2URL := "https://radiko.jp/v2/api/auth2"

	headers := make(http.Header)
	headers.Add("User-Agent", utils.UserAgent)
	headers.Add("x-radiko-device", XRadikoDevice)
	headers.Add("x-radiko-user", XRadikoUser)
	headers.Add("x-radiko-authtoken", token)
	headers.Add("x-radiko-partialkey", partialkey)

	body := miniReq.GetBody(auth2URL, headers, nil)
	areaSplit := strings.Split(string(body), ",")
	areaid = areaSplit[0]
	return
}

// radikoHLS 获取 AAC 下载地址
func radikoHLS(token string, areaid string, radioData *RadioData) (aacURLs []string) {
	playlistURL := "https://radiko.jp/v2/api/ts/playlist.m3u8"

	headers := make(http.Header)
	headers.Add("User-Agent", utils.UserAgent)
	headers.Add("X-Radiko-AuthToken", token)
	headers.Add("X-Radiko-AreaId", areaid)

	params := make(map[string]string)
	params["station_id"] = radioData.stationID
	params["start_at"] = radioData.startAt
	params["ft"] = radioData.ft
	params["end_at"] = radioData.endAt
	params["to"] = radioData.to
	params["l"] = radioData.l
	params["type"] = radioData.rtype

	chunklistBody := miniReq.GetBody(playlistURL, headers, params)
	chunklistURL := radikoChunklist(string(chunklistBody))

	m3u8ReqHeaders := make(http.Header)
	m3u8ReqHeaders.Add("User-Agent", utils.UserAgent)

	m3u8Body := miniReq.GetBody(chunklistURL, m3u8ReqHeaders, nil)
	aacURLs = radikoAAC(string(m3u8Body))
	if len(aacURLs) == 0 {
		aacURLs = []string{}
	}
	return
}

// rThread 并发下载
func rThread(url string, ch chan []byte, dl *DLServer) {
	headers := make(http.Header)
	headers.Add("User-Agent", utils.UserAgent)
	body := miniReq.GetBody(url, headers, nil)
	dl.WG.Done()
	<-dl.Gonum
	ch <- body
}

// rEngine 下载器
func rEngine(urls []string, savePath string) {
	aacFile, _ := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	total := float64(len(urls))
	var thread int
	dl := new(DLServer)
	if total < 16 {
		thread = int(total)
	} else {
		thread = 16
	}
	var num int

	ch := make([]chan []byte, 1024)
	dl.Gonum = make(chan string, thread)
	dl.WG.Add(len(urls))

	for i, url := range urls {
		dl.Gonum <- url
		ch[i] = make(chan []byte)
		go rThread(url, ch[i], dl)
	}
	for _, d := range ch {
		if num == int(total) {
			break
		}
		tmp, _ := <-d
		offset, _ := aacFile.Seek(0, 2)
		_, _ = aacFile.WriteAt(tmp, offset)
		num++
	}

	dl.WG.Wait()
	defer aacFile.Close()
}

// RadioFromDB 读取 Radio 的数据
func RadioFromDB(name string) (data map[string]interface{}) {
	radikoColl := DataBase.Collection("radiko")
	fData := bson.D{
		{Key: "name", Value: name},
	}
	radikoData := MFind(radikoColl, 0, 0, fData)
	if len(radikoData) != 0 {
		data = radikoData[0]
	} else {
		data = map[string]interface{}{}
	}
	return
}

// Radio2DB 保存 Radio 的数据
func Radio2DB(name, url string) {
	radikoColl := DataBase.Collection("radiko")
	rdata := radioJSON{
		Name: name,
		URL:  url,
	}
	MInsertOne(radikoColl, rdata)
}

// RadioGet 获取 Radio 的数据
func RadioGet(fileName, station, startAt, endAt string) (dlurl string) {
	// authkey := radikoJSKey(client)
	authkey := "bcd151073c03b352e1ef2fd66c32209da9ca0afa"

	radioData := new(RadioData)
	radioData.stationID = station
	radioData.startAt = startAt
	radioData.endAt = endAt
	radioData.ft = startAt
	radioData.to = endAt
	radioData.l = "15"
	radioData.rtype = "b"

	token, offset, length := radikoAuth1()
	partialkey := encodeKey(authkey, offset, length)
	area := radikoAuth2(token, partialkey)
	aacURLs := radikoHLS(token, area, radioData)

	if len(aacURLs) != 0 {
		mediaInfo := config.MediaCfg()
		mediaPath := mediaInfo["media_path"].(string)
		wwwPath := mediaInfo["www_path"].(string)

		savePath := filepath.Join(mediaPath, fileName)
		rEngine(aacURLs, savePath)
		dlurl = filepath.Join(wwwPath, fileName)
	} else {
		dlurl = ""
	}
	return
}
