package utils

import (
	"log"
	"strconv"
	"time"
)

// TimeCtlBasic 时间转换
type TimeCtlBasic struct{}

// TimeCtl 初始化
var TimeCtl *TimeCtlBasic

func init() {
	TimeCtl = NewTimeCtl()
}

// NewTimeCtl 初始化 TimeCtl
func NewTimeCtl() (n *TimeCtlBasic) {
	n = new(TimeCtlBasic)
	return
}

// AnyFormat 显示转换 2006/01/02 → 2006.01.02
func (tc *TimeCtlBasic) AnyFormat(oldLayout, newLayout string, t string) string {
	ns, err := time.Parse(oldLayout, t)
	if err != nil {
		log.Panic("Time Parse Error", err)
	}
	fTime := ns.Format(newLayout)
	log.Println(fTime)
	return fTime
}


// String2Unix 字符串转 unix 时间戳 layout → 1136142245
func (tc *TimeCtlBasic) String2Unix(layout string, t string) int64 {
	localTime, err := time.LoadLocation("Local")
	if err != nil {
		log.Panic("Timezone Error", err)
	}
	tString, err := time.ParseInLocation(layout, t, localTime)
	if err != nil {
		log.Panic("LocalTime Error", err)
	}
	uTime := tString.Unix()
	return uTime
}

// UnixInt2String unix int64 时间戳转字符串
func (tc *TimeCtlBasic) UnixInt2String(layout string, t int64) string {
	return time.Unix(t, 0).Format(layout)
}

// Unix2String unix 时间戳转字符串 1136142245 → layout
func (tc *TimeCtlBasic) Unix2String(layout string, t string) string {
	if len(t) == 13 {
		t = t[0 : len(t)-3]
	}
	timeInt, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		log.Panic("String to Int64 Error", err)
	}
	sTime := time.Unix(timeInt, 0).Format(layout)
	return sTime
}

// Utc2string UTC 时间转字符串 2006-01-02T03:04:50Z0700  → layout
func (tc *TimeCtlBasic) Utc2string(layout string, t string) string {
	rawTime, _ := time.Parse("2006-01-02T15:04:05Z0700", t)
	uTime := rawTime.Format(layout)
	return uTime
}
