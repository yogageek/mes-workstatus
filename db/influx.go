package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang/glog"
	influxdb "github.com/influxdata/influxdb1-client/v2"
)

const (
	Influx_addr     = "http://59.124.112.31:31593"
	Influx_username = "admin"
	Influx_password = "admin12345"
	Influx_database = "mydb"
)

var Influx_topic = os.Getenv("INFLUX_TOPIC")

func newInfluxConn() {
	fmt.Println("Start")

	//Influx Connection
	conn, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     Influx_addr,
		Username: Influx_username,
		Password: Influx_password,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success_Conn: ", conn)
}

func influxConn() influxdb.Client {
	fmt.Println("Start...")
	fmt.Println("ENV INFLUX_TOPIC=", Influx_topic)

	//Influx Connection
	conn, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     Influx_addr,
		Username: Influx_username,
		Password: Influx_password,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success_Conn: ", conn)

	return conn
}

//Influx Insert
func setBP() influxdb.BatchPoints {
	bp, _ := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database: Influx_database,
	})
	return bp
}

//EXTERNAL------------------------------------

//Influx conn instance
type Influx struct {
	Conn influxdb.Client
}

// NewInflux new influx conn instance
func NewInflux() *Influx {
	// o:= new(influx)
	// o.conn = getInfluxConn()
	return &Influx{
		Conn: influxConn(),
	}
}

//InfluxInsert Influx Insert
func (i *Influx) InfluxInsert(fields map[string]interface{}, tags map[string]string, time time.Time) {
	bp := setBP()
	pt, _ := influxdb.NewPoint(Influx_topic, tags, fields, time)
	bp.AddPoint(pt)
	err := i.Conn.Write(bp)
	if err != nil {
		glog.Error(err)
	} else {
		//# bug
		// fmt.Println("Success_Insert")
	}
}

//InfluxQuery Influx query
func (i *Influx) InfluxQuery(sql string) *influxdb.Response {
	q := influxdb.Query{
		Command:  sql,
		Database: Influx_database,
	}
	res, err := i.Conn.Query(q)
	if err != nil {
		glog.Error(err)
	}
	return res //may nil
}

//InfluxQuery Influx query
func (i *Influx) InfluxDrop() *influxdb.Response {
	sql := "DROP MEASUREMENT \"" + Influx_topic + "\""
	q := influxdb.Query{
		Command:  sql,
		Database: Influx_database,
	}
	res, err := i.Conn.Query(q)
	if err != nil {
		glog.Error(err)
	}
	fmt.Println("drop measurement...")
	return res //may nil
}
