package main

import (
	"fmt"
	"mes-workstatus/db"
	"mes-workstatus/logic"
	"time"
)

func init() {
	// 啟動清空influxdb
	ifx := db.NewInflux()
	ifx.InfluxDrop()
}

const (
	timer1sec = 5
	timer2sec = 10
)

func main() {
	//make a chan first
	input := make(chan interface{})
	//producer - produce the messages
	go func() {
		for i := 0; i < 5; i++ {
			input <- i //put data into chan
		}
		input <- "v1.0.3"
	}()

	timer1 := time.NewTimer(time.Second * 3) //這裡設幾秒 就會等幾秒
	timer2 := time.NewTimer(time.Second * 10)
	// fmt.Println("print this line first and then...")

	for {
		select {
		case msg := <-input: //take data from chan
			fmt.Println(msg) //will print helle world
		case <-timer1.C: //t1.C拿出channel
			println("timer1 sec...", timer1sec)
			logic.SelectPgIntoInflux()
			timer1.Reset(time.Second * timer1sec) //使t1重新開始計時
		case <-timer2.C:
			println("timer2 sec...", timer2sec)
			logic.SyncPgAndInfluxForDelete()
			timer2.Reset(time.Second * timer2sec)
		}
	}
}
