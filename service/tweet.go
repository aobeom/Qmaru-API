package service

import (
	"net/http"
	"net/url"
	"strings"

	"qmaru-api/config"
	"qmaru-api/utils"
)

func tweetStatusID(url string) (i string) {
	firstSplit := strings.Split(url, "/")
	idRaw := firstSplit[len(firstSplit)-1]
	lastSplit := strings.Split(idRaw, "?")
	i = lastSplit[0]
	return
}

func tweetToken(k, s string) (t string) {
	// https://developer.twitter.com/en/docs/basics/authentication/api-reference/token
	oauthAPI := "https://api.twitter.com/oauth2/token"

	headers := make(http.Header)
	headers.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	headers.Add("username", k)
	headers.Add("password", s)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	reader := strings.NewReader(data.Encode())

	body := utils.Minireq.PostBody(oauthAPI, headers, reader)
	bodyJSON := utils.DataConvert.String2Map(body)

	t = bodyJSON["access_token"].(string)
	return
}

func tweetData(statusID, token string) (vurl string) {
	showAPI := "https://api.twitter.com/1.1/statuses/show.json"
	headers := make(http.Header)
	headers.Add("Authorization", "Bearer "+token)

	params := make(map[string]string)
	params["id"] = statusID
	params["tweet_mode"] = "extended"

	res := utils.Minireq.GetBody(showAPI, headers, params)

	resJSON := utils.DataConvert.String2Map(res)

	if _, ok := resJSON["extended_entities"]; ok {
		tweetExtendedEntities := resJSON["extended_entities"].(map[string]interface{})
		tweetMediaList := tweetExtendedEntities["media"].([]interface{})
		tweetMedia := tweetMediaList[0].(map[string]interface{})
		tweetVideoInfo := tweetMedia["video_info"].(map[string]interface{})
		tweetVariants := tweetVideoInfo["variants"].([]interface{})
		tweetBitrate := 0.0
		for _, v := range tweetVariants {
			tweetValue := v.(map[string]interface{})
			if _, ok := tweetValue["bitrate"]; ok {
				tweetBitrateM := tweetValue["bitrate"].(float64)
				if tweetBitrateM > tweetBitrate {
					tweetBitrate = tweetBitrateM
					vurl = tweetValue["url"].(string)
				}
			}
		}
	} else {
		vurl = ""
	}
	return
}

// TweetVideo 获取 Tweet 视频
func TweetVideo(url string) (vurl string) {
	cfg := config.TweetCfg()
	token := cfg["token"].(string)
	if token == "" {
		key := cfg["twitter_key"].(string)
		secret := cfg["twitter_secret"].(string)
		token = tweetToken(key, secret)
	}

	statusID := tweetStatusID(url)
	vurl = tweetData(statusID, token)

	return
}
