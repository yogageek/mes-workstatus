package logic

import (
	"encoding/json"
	"fmt"
	"mes-workstatus/model"

	"github.com/wesovilabs/koazee"
)

// #deprecated
// func CheckIfWoHasNewData() (bool, int, int) {
// 	pg := db.NewPostgres().SqlDB
// 	defer pg.Close()
// 	var pg_woId int
// 	sqlStatement := `SELECT wo_id FROM mes.work_orders order by wo_id desc limit 1`
// 	err := pg.QueryRow(sqlStatement).Scan(&pg_woId)
// 	if err != nil {
// 		log.Printf("Unable to execute the query. %v", err)
// 	}

// 	var ifx_woId int
// 	r := db.NewInflux().InfluxQuery(`SELECT Id FROM "workstatus" order by time desc limit 1`)
// 	if r.Results[0].Series != nil {
// 		res := r.Results[0].Series[0].Values //會拿到所有欄位值
// 		for _, row := range res {
// 			// row[0] == time
// 			// fmt.Println("index:", i, " Id:", row[1])
// 			// fmt.Println(reflect.TypeOf(row[1]))
// 			id, err := row[1].(json.Number).Int64()
// 			if err != nil {
// 				glog.Error(err)
// 			}
// 			ifx_woId = int(id)
// 			// t, err := time.Parse(time.RFC3339, row[0].(string))
// 			// if err != nil {
// 			// 	log.Fatal(err)
// 			// }
// 			// fmt.Println(reflect.TypeOf(row[1]))
// 		}
// 		if pg_woId > int(ifx_woId) {
// 			return true, pg_woId, ifx_woId
// 		}
// 		return false, pg_woId, ifx_woId
// 	}
// 	glog.Error("influx no data")
// 	return true, 0, 0
// }
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func ppt(b interface{}) {
	s, _ := json.MarshalIndent(b, "", "\t")
	fmt.Print(string(s))
}

func SumGoodSpendedTimeBeta(wos []model.WorkOrders) int {
	Mwos := make(map[string][]model.WorkOrders)
	for _, wo := range wos {
		mid := wo.OrderInfo.ManorderId
		Mwos[mid] = append(Mwos[mid], wo)
	}

	fmt.Println("制令單總數", len(Mwos))
	for k, v := range Mwos {

		fmt.Println("制令單", k, "有", len(v), "張工單")
		fmt.Println("內容為", v, "\n")
		// spew.Dump(v)
		for i, v := range Mwos[k] {
			fmt.Println("工單順序", i, "工單內容為", v)
		}
	}
	return 0

}

//#1.0.1
func SumGoodSpendedTime(wos []model.WorkOrders) int {
	// fmt.Println(wos)
	var list = []string{}
	stream := koazee.StreamOf(wos)

	// output := stream.ForEach(func(wos model.WorkOrders) string { return wos.OrderInfo.ManorderId })
	stream.ForEach(func(m model.WorkOrders) {
		fmt.Println(m.OrderInfo.ManorderId)
		mid := m.OrderInfo.ManorderId
		list = append(list, mid)
	}).Do().RemoveDuplicates()
	fmt.Println(len(list))
	stream = koazee.StreamOf(list).RemoveDuplicates().Do()
	distinctManorders := stream.RemoveDuplicates().Out().Val()
	fmt.Println(distinctManorders)

	// output.RemoveDuplicates().Do()
	// fmt.Printf("stream: %v\n", stream.Out().Val())
	return 0
	/*
		//宣告空map
		mapStruct := make(map[string][]model.WorkOrders)
		listStruct := map[string][]model.WorkOrders{}
		//append if unique
		for _, wo := range wos {
			mid := wo.OrderInfo.ManorderId
			if _, value := mapStruct[mid]; !value { //如果mapStruct沒有這筆key
				mapStruct[mid] = wo //將這筆加入map
				listStruct = append(listStruct, wo)
			}
		}
		fmt.Println(mapStruct)
		fmt.Println(listStruct)
	*/
	// you can use the ,ok idiom to check for existing keys
	// if _, ok := set[1]; ok {
	// 	fmt.Println("element found")
	// } else {
	// 	fmt.Println("element not found")
	// }

	// stream := koazee.StreamOf(wos)
	// fmt.Print("stream.RemoveDuplicates(): ")
	// fmt.Println(stream.RemoveDuplicates().Out().Val())

	//無法把reflect value轉回struct
	// out, _ := stream.GroupBy(func(wos model.WorkOrders) string { return wos.OrderInfo.ManorderId })

	// fmt.Println("SumGoodSpendedTime:", out)

	// theType := out.Type()
	// fmt.Println(theType)

	// thevalue := reflect.ValueOf(out)
	// i := thevalue.Interface()

	// result := i.(map[string][]model.WorkOrders)
	// PrettyPrint(result)
	// stream2 := koazee.StreamOf(wos).ForEach(func(wos model.WorkOrders) {
	// 	fmt.Println(wos.State.AccGood, wos.TakeProcessTimeFromWo)
	// })
	// fmt.Println("Operations are not evaluated until we perform stream.Do()\n")
	// stream2.Do()

	// fmt.Println("SumGoodSpendedTime:", out)
	// datatype := out.Type()
	// fmt.Println(datatype)
	// var data map[string][]model.WorkOrders //依照manorderid分類
	// for k := range out {
	// 	keys = append(keys, k)
	// }
	// a := out.MapKeys
	// fmt.Print("\n%+v\n", a)

	// fmt.Printf("\n%+v\n", out)
	// PrettyPrint(out)

	return 0
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
func structToMap2() {

	// s := &Server{
	// 	Name:    "gopher",
	// 	ID:      123456,
	// 	Enabled: true,
	// }

	// => {"Name":"gopher", "ID":123456, "Enabled":true}
	// m := structs.Map(s)
}

// func structToMap(item interface{}) map[string]interface{} {

// 	res := map[string]interface{}{}
// 	if item == nil {
// 		return res
// 	}
// 	v := reflect.TypeOf(item)
// 	reflectValue := reflect.ValueOf(item)
// 	reflectValue = reflect.Indirect(reflectValue)

// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	for i := 0; i < v.NumField(); i++ {
// 		tag := v.Field(i).Tag.Get("json")
// 		field := reflectValue.Field(i).Interface()
// 		if tag != "" && tag != "-" {
// 			if v.Field(i).Type.Kind() == reflect.Struct {
// 				res[tag] = structToMap(field)
// 			} else {
// 				res[tag] = field
// 			}
// 		}
// 	}
// 	return res
// }
