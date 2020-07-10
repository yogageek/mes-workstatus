package logic

import (
	"fmt"
	"mes-workstatus/db"
	"mes-workstatus/model"
	"time"

	"github.com/fatih/structs"
)

const (
	t1 = "OrderIdTag"
	t2 = "ManorderIdTag"
	t3 = "WorkorderIdTag"
)

const (
	theTime int64 = 1593475200
)

var timeUnix time.Time = time.Unix(theTime, 0) //時間轉unix格式->2020-06-30 08:00:00 +0800 CST

//定時 pg into influx
func SelectPgIntoInflux() {

	ifx := db.NewInflux()

	//工單
	wofields := model.WoFields{}
	workorders, _ := selectWorkorder()
	for _, o := range workorders {
		wofield := model.WoField{}
		wofield.WorkorderId = o.WorkorderId
		wofield.Seq = o.OrderInfo.StepId
		wofield.Process = o.TakeTextSeq()
		wofield.Qty = o.OrderInfo.Qty
		wofield.GoodQty = o.State.AccGood
		wofield.NgQty = o.State.AccNg
		wofield.Status = transStatusStrToInt(o.State.Status)
		wofield.ProcessTimePerPcs = o.TakeProcessTimeFromWo()
		wofield.StardardCompletionTime = calStardardCompletionTime(wofield.Qty, wofield.ProcessTimePerPcs)
		wofield.CurrentStatus = calCompletionPctStr(wofield.GoodQty, wofield.NgQty, wofield.Qty)
		// fmt.Printf("\n%+v\n", wofield)

		tags := struct2MapString(model.WoTag{
			Type:           "processStatus",
			WorkorderIdTag: o.WorkorderId,
			ManorderIdTag:  o.OrderInfo.ManorderId,
			OrderIdTag:     o.OrderInfo.OrderId,
		})

		// 		n := int64(32)
		// str := strconv.FormatInt(n, 10)
		// fmt.Println(str)  // Prints "32"
		// s := strconv.FormatInt(time.Now().Unix(), 10) // s == "61" (hexadecimal)
		// fmt.Println(s)
		// fmt.Println(time.Now().Unix())

		//#requirement peter讓工序能排序
		timeSeq := theTime + int64(wofield.Seq)
		timeSeqUnix := time.Unix(timeSeq, 0)

		ifx.InfluxInsert(structs.Map(wofield), tags, timeSeqUnix)

		//加入
		wofields.Add(wofield)
	}

	//制令單
	manorders, _ := selectOrderJoinManorder()
	for _, o := range manorders {
		manorderfield := &model.ManorderField{
			PlanProduction: o.Qty,
			//#1.0.1 計算增加包含 Manorder層級 時間.pcs 加權比重
			OrderOEE: calCompletionPct(o.WorkedTime, o.AccNg, o.RequiredTime),
			DueDate:  o.DueDate.Format("2006-01-02 15:04:05"), //一定要放這個日期才能轉 2006-01-02 15:04:05
		}
		// fmt.Println(o.DueDate)
		// fmt.Println(o.DueDate.Round(time.Minute))    //not work
		// fmt.Println(o.DueDate.Truncate(time.Minute)) //not work

		manordertag := struct2MapString(model.ManorderTag{
			Type:          "manorderStatus",
			ManorderIdTag: o.ManorderId,
			OrderIdTag:    o.OrderId,
		})
		ifx.InfluxInsert(structs.Map(manorderfield), manordertag, timeUnix)
	}

	//定單
	orders, _ := selectOrderJoinWorkorder()
	for _, o := range orders {
		orderfield := &model.OrderField{
			ProductId:      o.ProductId,
			ProductName:    o.ProductName,
			PlanProduction: o.Qty,
			//#1.0.1 計算增加包含 Manorder層級 時間.pcs 加權比重
			OrderOEE: calCompletionPct(o.WorkedTime, o.AccNg, o.RequiredTime),
		}
		ordertag := struct2MapString(model.OrderTag{
			Type:       "orderStatus",
			OrderIdTag: o.OrderId,
		})
		ifx.InfluxInsert(structs.Map(orderfield), ordertag, timeUnix)
	}

}

//cch新方法:插入influx時field帶指定值, 查詢全部出來後用go判斷field值是否有變 拿到field的tag後去刪
func SyncPgAndInfluxForDelete() {
	fmt.Println("SyncPgAndInfluxForDelete...")
	ifx := db.NewInflux()
	//查influx
	//iql := fmt.Sprintf("SHOW TAG VALUES FROM %s with KEY= \"%s\"", ifx.Topic, "ManorderIdTag")
	iql := fmt.Sprintf("SELECT * FROM %s ", ifx.Topic)
	data := ifx.InfluxQuery(iql)
	iorders := getDiscinctIds(t1, data)
	imanorders := getDiscinctIds(t2, data)
	iworkorders := getDiscinctIds(t3, data)
	influxids := append(append(iorders, imanorders...), iworkorders...)
	fmt.Println("len of influx idTag rows:", len(influxids))

	//查pg
	porders, _ := selectDistc(t1)
	pmanorders, _ := selectDistc(t2)
	pworkorders, _ := selectDistc(t3)
	pgids := append(append(porders, pmanorders...), pworkorders...)
	fmt.Println("len of pg idTag rows:", len(pgids))

	//刪除
	del := func(delId string) {
		var dql string
		dql = fmt.Sprintf("DELETE FROM %s WHERE OrderIdTag = '%s'", ifx.Topic, delId)
		res1 := ifx.InfluxDelete(dql)
		dql = fmt.Sprintf("DELETE FROM %s WHERE ManorderIdTag = '%s'", ifx.Topic, delId)
		res2 := ifx.InfluxDelete(dql)
		dql = fmt.Sprintf("DELETE FROM %s WHERE WorkorderIdTag = '%s'", ifx.Topic, delId)
		res3 := ifx.InfluxDelete(dql)
		fmt.Println("delete influx idTag=", delId, "result:", res1, res2, res3)
	}

	//如果influxid不在pgid裡面
	for _, influxid := range influxids {
		if !stringInSlice(influxid, pgids) {
			del(influxid)
		}
	}
}
