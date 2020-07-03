package logic

import (
	"mes-workstatus/db"
	"mes-workstatus/model"
	"time"

	"github.com/fatih/structs"
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
		wofield.Status = TransStatusStrToInt(o.State.Status)
		wofield.ProcessTimePerPcs = o.TakeProcessTimeFromWo()
		wofield.StardardCompletionTime = CalStardardCompletionTime(wofield.Qty, wofield.ProcessTimePerPcs)
		wofield.CurrentStatus = CalCompletionPctStr(wofield.GoodQty, wofield.NgQty, wofield.Qty)
		// fmt.Printf("\n%+v\n", wofield)

		tags := Struct2MapString(model.WoTag{
			Type:           "processStatus",
			WorkorderIdTag: o.WorkorderId,
			ManorderIdTag:  o.OrderInfo.ManorderId,
			OrderIdTag:     o.OrderInfo.OrderId,
		})

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
			OrderOEE: CalCompletionPct(o.WorkedTime, o.AccNg, o.RequiredTime),
			DueDate:  o.DueDate.Format("2006-01-02 15:04:05"), //一定要放這個日期才能轉 2006-01-02 15:04:05
		}
		// fmt.Println(o.DueDate)
		// fmt.Println(o.DueDate.Round(time.Minute))    //not work
		// fmt.Println(o.DueDate.Truncate(time.Minute)) //not work

		manordertag := Struct2MapString(model.ManorderTag{
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
			OrderOEE: CalCompletionPct(o.WorkedTime, o.AccNg, o.RequiredTime),
		}
		ordertag := Struct2MapString(model.OrderTag{
			Type:       "orderStatus",
			OrderIdTag: o.OrderId,
		})
		ifx.InfluxInsert(structs.Map(orderfield), ordertag, timeUnix)
	}
}

// 不用sql join 而用程式group寫法
// func test(fs model.WoFields) {
// 	wofields := fs.WoFields
// 	// bigset := make(map[string]map[string]struct{}) //key: string value: map[string]struct{}

// 	distinctOrders := make([]string, 0)
// 	for _, f := range wofields { //all work orders
// 		oid := f.OrderId
// 		distinctOrders = append(distinctOrders, oid)
// 	}
// 	distinctOrders = RemoveRepeatedElement(distinctOrders)

// 	// set := make(map[string]interface{})
// 	var dary map[string]map[string][]model.WoField //key: string value: map[string]struct{} //{78:{78-1,78-2}}, {79:{79-1,79-2}}
// 	for _, dId := range distinctOrders {
// 		dset := make(map[string][]model.WoField) //放78的array 78:{78-1,78-2}
// 		for _, f := range wofields {
// 			if dId == f.OrderId {
// 				dset[dId] = append(dset[dId], f)
// 			}
// 		}
// 		fmt.Printf("\n%+v\n", dset)
// 		dary[dId] = dset
// 	}
// 	fmt.Printf("\n%+v\n", dary)
// 判断key是否存在的问题
// if v, ok := m["a"]; ok {
// 	fmt.Println(v)
// } else {
// 	fmt.Println("Key Not Found")
// }

// orderIds := RemoveRepeatedElement(fs.OrderIds)
// for _, id := range orderIds {
// 	set := make(map[string]struct{})
// 	// set[id] =
// }

// distinctOrderIDs := RemoveRepeatedElement(oid)
// for _, id := range distinctOrderIDs {
// 	fmt.Println(id)
// }
// }
