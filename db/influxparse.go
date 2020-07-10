package db

import (
	"fmt"
	// . "mes-workstatus/model"

	influxdb "github.com/influxdata/influxdb1-client/v2"
)

//查詢出來的多筆結果轉map
func ParseQueryResult(response *influxdb.Response) []map[string]interface{} {
	//response.Results
	for _, result := range response.Results {
		//series is group of result (like group by consumerId, each conusmerId in one serie)
		for _, series := range result.Series {
			// rowsdata := RowsData{} //多筆資料
			// consumerId := series.Tags["consumerId"]
			// usage.ConsumerId = consumerId

			//values[row][col]
			if len(series.Columns) > 0 {
				var multiRows []map[string]interface{}
				for row := 0; row < len(series.Values); row++ {
					// rowdata := RowData{}
					// var cvs []map[string]interface{}
					// cvs := []map[string]interface
					// cvs := []ColumnValue{}

					//columns
					for i := 1; i < len(series.Values[row]); i++ {
						oneRow := make(map[string]interface{})
						// cv := make(map[string]interface{})
						// cv := map[string]interface{}
						// cv := ColumnValue{}.CV

						colname := series.Columns[i] // cv[col] = ""
						if series.Values[row][i] != nil {
							// dataUsage.Quantity, _ = series.Values[row][i].(json.Number).Float64()
							// fmt.Println(series.Values[row][i])
							oneRow[colname] = series.Values[row][i]
							// //int
							// if value, ok := series.Values[row][i].(json.Number); ok {
							// 	cv[colname] = value.Float64
							// }
							// //string
							// if value, ok := series.Values[row][i].(json.); ok {
							// 	cv[colname] = value.string
							// }
						}
						multiRows = append(multiRows, oneRow)
					}
					// rowdata.Row = cvs //一筆資料
					// rowsdata.Rows = append(rowsdata.Rows, rowdata)
				}
				// fmt.Println(len(multiRows)) //not work
				return multiRows
			}
		}
	}
	var none []map[string]interface{}
	fmt.Println("parse result no data")
	return none
}
