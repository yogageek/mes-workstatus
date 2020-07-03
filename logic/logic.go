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
	//# bug
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
	sqlStatement := `select orders.order_datetime,(orders.order_info->>'order_id')::text as order_id,T.qty as total_qty,orders.qty,COALESCE(T.acc_good,0)as acc_good,COALESCE(T.acc_ng,0)as acc_ng,(orders.product->>'id') as productID,(orders.product->>'name') as productName from mes.orders left join 
	(select order_id,sum(acc_good) as acc_good , sum(acc_ng) as acc_ng, sum(qty) as qty from
	(select (work_orders.order_info->>'order_id')::text as order_id, (work_orders.order_info->>'qty')::int as qty, workorder_id,(work_orders.state->>'acc_good')::int as acc_good,(work_orders.state->>'acc_ng')::int as acc_ng from mes.work_orders) as rTable
	group by order_id) T
	on (orders.order_info->>'order_id')::text = T.order_id
	order by T.order_id
	`
	rows, err := db.Query(sqlStatement)
	//defer不能放這 如果有err會錯
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return nil, err //需要return 不然接著rows.next會錯
	}
	defer rows.Close() //if row=nil, will cause this error
	for rows.Next() {  // iterate over the rows
		var m model.OrderJoinWorkorder
		err = rows.Scan(&m.Time, &m.OrderId, &m.TotalQty, &m.Qty, &m.AccGood, &m.AccNg, &m.ProductId, &m.ProductName)
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
