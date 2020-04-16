package utils

import (
	"encoding/json"
	"fmt"
)

// DataConvBasic 数据类型转换
type DataConvBasic struct{}

// DataConvert 数据转换
var DataConvert *DataConvBasic

func init() {
	DataConvert = NewDataConv()
}

// NewDataConv 初始化 DataConv
func NewDataConv() (n *DataConvBasic) {
	n = new(DataConvBasic)
	return
}

// Map2String map类型保存 JSON
func (dc *DataConvBasic) Map2String(data interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}

// String2Maps 读取 JSON 转换为 []map
func (dc *DataConvBasic) String2Maps(data []byte) (m []map[string]string) {
	json.Unmarshal(data, &m)
	return
}

// String2Map 读取 JSON 转换为 map
func (dc *DataConvBasic) String2Map(data []byte) (m map[string]interface{}) {
	json.Unmarshal(data, &m)
	return
}

// String2Array 读取 JSON 转换为数组
func (dc *DataConvBasic) String2Array(data []byte) (m interface{}) {
	json.Unmarshal(data, &m)
	return
}

// Print 打印输出
func (dc *DataConvBasic) Print(i interface{}) {
	fmt.Println(i)
}
