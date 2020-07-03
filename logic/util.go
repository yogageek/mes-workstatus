package logic

import (
	"math"
	"strconv"
)

func TransStatusStringToInt(status string) int {
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
//#1.0.1 需要加乘process time比重
/*
1.0.0
良品數/制訂單總數(qty*工單數)
1.0.1 完成時間比例
各工單(良品數*標準完成時間)/制令單完成時間(duration欄位)
*/

func CalCompletionPct(goodqty, ngqty, totalqty int) int {

	//#1.0.1 完成率只能用良品去計算
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
//#1.0.1 需要加乘process time比重
func CalCompletionPctStr(goodqty, ngqty, totalqty int) string {
	//#1.0.1 完成率只能用良品去計算
	ngqty = 0

	// 返回字串 50/100
	a := strconv.Itoa(goodqty + ngqty)
	b := strconv.Itoa(totalqty)
	c := a + "/" + b
	return c
}

//計算標準完成時間
func CalStardardCompletionTime(totalqty, timeperPcs int) int {
	return totalqty * timeperPcs
}
