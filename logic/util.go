package logic

import (
	"math"
	"reflect"
	"strconv"
)

func transStatusStrToInt(status string) int {
	switch status {
	case "未進站":
		return 0
	case "生產中":
		return 1
	case "已出站":
		return 2
	default:
		return 3
	}
}

//計算完成率
func calCompletionPct(goodqty, ngqty, totalqty int) int {
	//完成率只能用良品去計算
	ngqty = 0

	// 返回整數百分比
	var percent float64
	done := float64(goodqty + ngqty)
	total := float64(totalqty)
	percent = (done / total) * 100
	// fmt.Println(done, total, percent)
	if math.IsNaN(percent) {
		// log.Println("percent is Nan")
		return 0
	}
	return int(percent)
}

//計算完成率
func calCompletionPctStr(goodqty, ngqty, totalqty int) string {
	//完成率只能用良品去計算
	ngqty = 0

	// 返回字串 ex:50/100
	a := strconv.Itoa(goodqty + ngqty)
	b := strconv.Itoa(totalqty)
	c := a + "/" + b
	return c
}

//計算標準完成時間
func calStardardCompletionTime(totalqty, timeperPcs int) int {
	return totalqty * timeperPcs
}

func struct2MapString(obj interface{}) map[string]string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).String()
	}
	return data
}

// func struct2Map(obj interface{}) map[string]interface{} {
// 	t := reflect.TypeOf(obj)
// 	v := reflect.ValueOf(obj)

// 	var data = make(map[string]interface{})
// 	for i := 0; i < t.NumField(); i++ {
// 		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
// 	}
// 	return data
// }

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
