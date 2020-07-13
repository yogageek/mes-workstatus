package logic

import (
	"fmt"
	"log"
	"mes-workstatus/db"
	"mes-workstatus/model"
)

func selectWorkorder() ([]model.WorkOrders, error) {
	fmt.Print("exec selectWorkorder...")
	var workorders []model.WorkOrders

	db := db.NewPostgres().SqlDB
	//# fix 取消db.close 修正初始化一次連線後被關掉導致無法重複利用的問題
	//defer db.Close()

	// 暫時過濾不符規定的測試資料
	// sqlStatement := `SELECT timestamp, wo_id, workorder_id, state, manorder FROM mes.work_orders order by wo_id`
	sqlStatement := `SELECT timestamp, workorder_id, state, order_info, product FROM mes.work_orders where timestamp is not null order by workorder_id`
	rows, err := db.Query(sqlStatement)

	//defer不能放這 如果有err會錯
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return nil, err //需要return 不然接著rows.next會錯
	}
	defer rows.Close() //if row=nil, will cause this error

	for rows.Next() { // iterate over the rows
		var m model.WorkOrders
		err = rows.Scan(&m.Time, &m.WorkorderId, &m.State, &m.OrderInfo, &m.Product)
		if err != nil {
			log.Printf("Unable to scan the row. %v", err)
		}
		// get any error encountered during iteration
		err = rows.Err()
		if err != nil {
			log.Printf("err during iteration. %v", err)
		}
		workorders = append(workorders, m)

		// fmt.Println(" time: ", m.Timestamp, " serial id: ", m.WoId, " order id: ", takeOIDFromWo(m.WorkorderId), " wo id: ", takeWoIDFromWo(m.WorkorderId), " status/good/ng: ", m.State, " qty: ", m.Manorder.Qty, " text: ", takeTextFromWo(m, m.WorkorderId), " processtime: ", takeProcessTimeFromWo(m, m.WorkorderId))

	}
	return workorders, err
}

func selectOrderJoinWorkorder() ([]model.OrderJoinWorkorder, error) {
	fmt.Print("exec selectOrderJoinWorkorder...")
	var orderJoinWorkorders []model.OrderJoinWorkorder
	// orderProcessQty := make(map[string]int)

	db := db.NewPostgres().SqlDB
	//# fix 取消db.close 修正初始化一次連線後被關掉導致無法重複利用的問題
	//defer db.Close()
	sqlStatement := selectOrderJoinWorkorderGroupByOrderIncludeCalTime
	rows, err := db.Query(sqlStatement)
	//defer不能放這 如果有err會錯
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return nil, err //需要return 不然接著rows.next會錯
	}
	defer rows.Close() //if row=nil, will cause this error
	for rows.Next() {  // iterate over the rows
		var m model.OrderJoinWorkorder
		err = rows.Scan(&m.Time, &m.OrderId, &m.TotalQty, &m.Qty, &m.AccGood, &m.AccNg, &m.RequiredTime, &m.WorkedTime)
		if err != nil {
			log.Printf("Unable to scan the row. %v", err)
		}
		// get any error encountered during iteration
		err = rows.Err()
		if err != nil {
			log.Printf("err during iteration. %v", err)
		}
		orderJoinWorkorders = append(orderJoinWorkorders, m)

		// fmt.Printf("orderProcessQty: %+v\n ", m)
	}
	return orderJoinWorkorders, err
}

//#1.0.1 增加制令單預計完成日
func selectOrderJoinManorder() ([]model.OrderJoinManorder, error) {
	fmt.Print("exec selectOrderJoinManorder...")
	var orderJoinManorders []model.OrderJoinManorder

	db := db.NewPostgres().SqlDB
	//# fix 取消db.close 修正初始化一次連線後被關掉導致無法重複利用的問題
	//defer db.Close()

	sqlStatement := selectOrderJoinWorkorderGroupByManorderIncludeCalTime
	rows, err := db.Query(sqlStatement)
	//defer不能放這 如果有err會錯
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return nil, err //需要return 不然接著rows.next會錯
	}
	defer rows.Close() //if row=nil, will cause this error

	for rows.Next() { // iterate over the rows
		var m model.OrderJoinManorder
		err = rows.Scan(&m.Time, &m.DueDate, &m.ManorderId, &m.OrderId, &m.TotalQty, &m.Qty, &m.RequiredTime, &m.WorkedTime, &m.AccGood, &m.AccNg, &m.ProductId, &m.ProductName)
		if err != nil {
			log.Printf("Unable to scan the row. %v", err)
		}
		// get any error encountered during iteration
		err = rows.Err()
		if err != nil {
			log.Printf("err during iteration. %v", err)
		}
		orderJoinManorders = append(orderJoinManorders, m)

		// fmt.Println(" time: ", m.Timestamp, " serial id: ", m.WoId, " order id: ", takeOIDFromWo(m.WorkorderId), " wo id: ", takeWoIDFromWo(m.WorkorderId), " status/good/ng: ", m.State, " qty: ", m.Manorder.Qty, " text: ", takeTextFromWo(m, m.WorkorderId), " processtime: ", takeProcessTimeFromWo(m, m.WorkorderId))
	}
	return orderJoinManorders, err
}

func selectDistc(theIdTag string) ([]string, error) {
	fmt.Printf("exec selectDistc(%v)...", theIdTag)

	db := db.NewPostgres().SqlDB
	//# fix 取消db.close 修正初始化一次連線後被關掉導致無法重複利用的問題
	//defer db.Close()

	var sqlStatement string
	switch theIdTag {
	case t1:
		// --查詢訂單
		sqlStatement = `select distinct ((work_orders.order_info->>'order_id')::text) T from mes.work_orders order by T`
	case t2:
		// --查詢制令單
		sqlStatement = `select distinct ((work_orders.order_info->>'manorder_id')::text) T  from mes.work_orders order by T`
	case t3:
		// 	--查詢工單
		sqlStatement = `select distinct workorder_id from mes.work_orders order by workorder_id`
	}

	rows, err := db.Query(sqlStatement)
	//defer不能放這 如果有err會錯
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return nil, err //需要return 不然接著rows.next會錯
	}
	defer rows.Close() //if row=nil, will cause this error

	var ids []string
	for rows.Next() { // iterate over the rows
		var id string
		err = rows.Scan(&id)
		if err != nil {
			log.Printf("Unable to scan the row. %v", err)
		}
		// get any error encountered during iteration
		err = rows.Err()
		if err != nil {
			log.Printf("err during iteration. %v", err)
		}
		ids = append(ids, id)
	}
	return ids, err
}

func getDiscinctIds(theIdtag string, data []map[string]interface{}) []string {
	var discinctIdList []string
	for _, d := range data {
		for k, v := range d { //一筆資料
			if k == theIdtag {
				if !stringInSlice(v.(string), discinctIdList) {
					discinctIdList = append(discinctIdList, v.(string))
				}
			}
		}
	}
	fmt.Printf("%v len=%v \n", theIdtag, len(discinctIdList))
	return discinctIdList
}
