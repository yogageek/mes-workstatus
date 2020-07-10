package logic

import (
	"mes-workstatus/model"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

/*
	// usages := map[string]float64{}
	usages := map[string]string{}
	for _, result := range results {
		fmt.Println("result:", result) //同
		fmt.Println()
		for _, series := range result.Series {
			fmt.Println("series:", series)
			fmt.Println()
			//查到幾個欄位就跑幾次
			fmt.Println(len(series.Columns))
			//i=0是欄位名 i=1是值
			for i := 1; i < len(series.Columns); i++ { // only columns[0],columns[1] =>[xxIdTag],["20200404"]
				// if columnName := series.Columns[i]; columnName == "" {

				// }
				for row := 0; row < len(series.Values); row++ {
					fmt.Println(series.Values[row][i])
				}
			}
		}
	}
	fmt.Println(usages)
*/

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

//process time == time per pcs of this wo
func takeProcessTimeFromWo(m model.WorkOrders, stepid int) int {
	stepidstr := strconv.Itoa(stepid)
	// if stepid == "" {
	// 	return 0
	// }
	// trimmedStepid := strings.Split(stepid, "-")[1]
	for _, w := range m.Product.Route {
		// fmt.Println(" id: ", w.Id)
		if w.Id == stepidstr {
			for _, t := range w.Lines {
				// fmt.Println("id: ", w.Id, " processtime: ", t.ProcessTime)
				t, err := strconv.Atoi(t.ProcessTime)
				if err != nil {
					glog.Error(err)
				}
				return t
			}
		}
	}
	glog.Error("time per pcs of this workorder not found ")
	return 0
}

func takeTextFromWo(m model.WorkOrders, stepid int) string {
	stepidstr := strconv.Itoa(stepid)
	// if stepid == "" {
	// 	return ""
	// }
	// trimmedStepid := strings.Split(stepid, "-")[1]
	for _, w := range m.Product.Route {
		// fmt.Println(" id: ", w.Id)
		if w.Id == stepidstr {
			return w.Text
		}
	}
	glog.Error("text not found ")
	return ""
}

//depre
func takeOIDFromWo(stepid string) string {
	if stepid == "" {
		return ""
	}
	oid := strings.Split(stepid, "-")[0]
	return oid
}

//depre
func takeWoIDFromWo(stepid string) string {
	if stepid == "" {
		return ""
	}
	woid := strings.Split(stepid, "-")[1]
	return woid
}
