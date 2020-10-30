package service

import (
	"github.com/antchfx/htmlquery"
	"qmaru-api/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type programJSON struct {
	Keyword   string              `bson:"keyword"`
	AreaCode  string              `bson:"area_code"`
	YahooURL  string              `bson:"yahoo_url"`
	ProgInfo  []map[string]string `bson:"prog_info"`
	CreatedAt time.Time           `bson:"created_at"`
}

// ProgramFromDB 读取 Program 的数据
func ProgramFromDB(kw, ac string) (data map[string]interface{}) {
	programColl := DataBase.Collection("program_info")
	fData := bson.D{
		{Key: "keyword", Value: kw},
		{Key: "areacode", Value: ac},
	}
	programData := MFind(programColl, 0, 0, fData)
	if len(programData) != 0 {
		data = programData[0]
	} else {
		data = map[string]interface{}{}
	}
	return
}

// Program2DB 保存 Program 的数据
func Program2DB(kw, ac, tvurl string, tvinfo []map[string]string) {
	programColl := DataBase.Collection("program_info")
	var pdata programJSON
	pdata.Keyword = kw
	pdata.AreaCode = ac
	pdata.YahooURL = tvurl
	pdata.ProgInfo = tvinfo
	pdata.CreatedAt = time.Now()
	MInsertOne(programColl, pdata)
}

// YahooTV 获取 YahooTV 数据
func YahooTV(kw, code string) (tvurl string, tvinfo []map[string]string) {
	yahooSite := "https://tv.yahoo.co.jp"
	apiURL := yahooSite + "/search/category/"

	headers := utils.MiniHeaders{
		"User-Agent": utils.UserAgent,
	}

	data := utils.MiniFormData{
		"q":   kw,
		"a":   code,
		"oa":  "1",
		"tv":  "1",
		"bsd": "1",
	}

	utils.Minireq.NoRedirect()
	tRes := utils.Minireq.Post(apiURL, headers, data)
	tLocation, _ := tRes.RawRes.Location()
	tRealURL := tLocation.String()
	res := utils.Minireq.Get(tRealURL)
	tvurl = tRealURL

	doc, _ := htmlquery.Parse(strings.NewReader(string(res.RawData())))
	nodes := htmlquery.Find(doc, `//ul[@class="programlist"]/li`)
	for _, node := range nodes {
		d := make(map[string]string)
		emTag := htmlquery.Find(node, `./div[@class="leftarea"]/p/em`)
		udate := htmlquery.InnerText(emTag[0])
		utime := htmlquery.InnerText(emTag[1])

		aTag := htmlquery.FindOne(node, `./div[@class="rightarea"]/p[1]/a`)
		title := htmlquery.InnerText(aTag)
		yurl := htmlquery.SelectAttr(aTag, "href")

		spanTag := htmlquery.FindOne(node, `./div[@class="rightarea"]/p[2]/span[1]`)
		station := htmlquery.InnerText(spanTag)

		d["date"] = udate
		d["time"] = utime
		d["url"] = yahooSite + yurl
		d["title"] = title
		d["station"] = station
		tvinfo = append(tvinfo, d)
	}
	return
}
