package main_peter

import (
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

const (
	Influx_addr     = "http://59.124.112.31:31593"
	Influx_username = "admin"
	Influx_password = "admin12345"
	Influx_database = "mydb"
	Influx_topic    = "workstatus"
)

func main() {

	fmt.Println("Start")

	//Influx Connection
	conn, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     Influx_addr,
		Username: Influx_username,
		Password: Influx_password,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success_Conn: ", conn)

	//Influx Insert
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database: Influx_database,
	})

	//固定
	tags := map[string]string{
		"workorderID":            "78-1",
		"stardardCompletionTime": "100s",
		"totalQty":               "90", //總數
	}
	//浮動
	fields := map[string]interface{}{
		"goodQty":              80, //良品
		"ngQty":                10, //不良品
		"completionPercentage": 50, //完成率(即時)

	}

	pt, _ := client.NewPoint(Influx_topic, tags, fields, time.Now())

	bp.AddPoint(pt)

	if err := client.Client.Write(conn, bp); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success_Insert")

	// endTime := time.Now()
	// endUnixTime := endTime.Unix()
	// startUnixTime := endUnixTime - 100000
	// endInfluxTs := endUnixTime * 1000000000
	// startInfluxTs := startUnixTime * 1000000000

	//Influx Query
	q := client.Query{
		Command:  "show tag keys from \"III/+/ITO\"", //query直接帶tablename
		Database: Influx_database,
	}
	fmt.Println("Query: ", q)

	res, _ := conn.Query(q)
	fmt.Println("Response: ", res.Results[0].Series[0].Values)

	// sT := strings.Replace("2020-04-14T09:00:19.886Z", "T", " ", 1)
	// var year, mon, day, hh, mm, ss int
	// fmt.Sscanf(sT, "%d-%d-%d %d:%d:%d ", &year, &mon, &day, &hh, &mm, &ss)
	// time_string_to_parse := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d+00:00", year, mon, day, hh, mm, ss)
	// sT_url_create_time, _ := time.Parse(time.RFC3339, time_string_to_parse)

	// eT := strings.Replace("2020-04-14T09:00:24", "T", " ", 1)
	// fmt.Sscanf(eT, "%d-%d-%d %d:%d:%d ", &year, &mon, &day, &hh, &mm, &ss)
	// time_string_to_parse = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d+00:00", year, mon, day, hh, mm, ss)
	// eT_url_create_time, _ := time.Parse(time.RFC3339, time_string_to_parse)

	// fmt.Println(sT_url_create_time)
	// fmt.Println(eT_url_create_time)

	// difference := eT_url_create_time.Sub(sT_url_create_time)

	// fmt.Println(difference)
	// result := difference.Seconds() / 288
	// fmt.Println(result)

}
