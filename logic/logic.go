package logic

import (
	"fmt"
	"log"
	"mes-workstatus/db"
	"mes-workstatus/model"
)

func selectWorkorder() ([]model.WorkOrders, error) {
	fmt.Println("exec selectWorkorder...")
	var workorders []model.WorkOrders

	db := db.NewPostgres().SqlDB
	defer db.Close()
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
	fmt.Println("exec selectOrderJoinWorkorder...")
	var orderJoinWorkorders []model.OrderJoinWorkorder
	// orderProcessQty := make(map[string]int)

	db := db.NewPostgres().SqlDB
	defer db.Close()
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
		err = rows.Scan(&m.Time, &m.OrderId, &m.TotalQty, &m.Qty, &m.AccGood, &m.AccNg, &m.RequiredTime, &m.WorkedTime, &m.ProductId, &m.ProductName)
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
	fmt.Println("exec selectOrderJoinManorder...")
	var orderJoinManorders []model.OrderJoinManorder

	db := db.NewPostgres().SqlDB
	defer db.Close()

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
		err = rows.Scan(&m.Time, &m.DueDate, &m.ManorderId, &m.OrderId, &m.TotalQty, &m.Qty, &m.RequiredTime, &m.WorkedTime, &m.AccGood, &m.AccNg)
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
