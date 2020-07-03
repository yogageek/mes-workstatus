package main

import (
	"fmt"
	"mes-workstatus/logic"
	"time"
)

func init() {
	// 啟動清空influxdb
	// ifx.InfluxDrop()
}

func main() {
	//make a chan first
	input := make(chan interface{})
	//producer - produce the messages
	go func() {
		for i := 0; i < 5; i++ {
			input <- i //put data into chan
		}
		input <- "v1.0.1"
	}()

	t1 := time.NewTimer(time.Second * 1) //這裡設幾秒 就會等幾秒
	// t2 := time.NewTimer(time.Second * 10)
	fmt.Println("print this line first and then...")

	for {
		select {
		//consumer - consume the messages
		case msg := <-input: //take data from chan
			fmt.Println(msg) //will print helle world
		case <-t1.C: //t1.C拿出channel
			println("t1s timer")
			logic.SelectPgIntoInflux()
			t1.Reset(time.Second * 5) //使t1重新開始計時
			// case <-t2.C:
			// 	println("10s timer")
			// 	t2.Reset(time.Second * 10)
		}
	}
}
