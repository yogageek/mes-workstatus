func main() {
    input := make(chan interface{})
    //producer - produce the messages
    go func() {
    for i := 0; i < 5; i   {
    input <- i
    }
    input <- "hello, world"
    }()
    t1 := time.NewTimer(time.Second * 5)
    t2 := time.NewTimer(time.Second * 10)
    for {
    select {
    //consumer - consume the messages
    case msg := <-input:
    fmt.Println(msg)
    case <-t1.C:
    println("5s timer")
    t1.Reset(time.Second * 5)
    case <-t2.C:
    println("10s timer")
    t2.Reset(time.Second * 10)
    }
    }
}


// func timetest() {
// 	c := cron.New()

// 	c.AddFunc("*/1 * * * * *", func() {
// 		fmt.Println("every 1 seconds executing")
// 		ifnew, woId := logic.CheckIfWoHasNewData()
// 		fmt.Println(ifnew, woId)
// 	})

// 	go c.Start()
// 	defer c.Stop()

// 	select {
// 	case <-time.After(time.Second * 100):
// 		return
// 	}
// }

	// var t time.Time
	// fmt.Println(t)

	// sT := strings.Replace("2020-06-30T00:00:00.000Z", "T", " ", 1)
	// var year, mon, day, hh, mm, ss int
	// fmt.Sscanf(sT, "%d-%d-%d %d:%d:%d ", &year, &mon, &day, &hh, &mm, &ss)
	// timeStringToParse := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d+00:00", year, mon, day, hh, mm, ss)
	// sTUrlCreateTime, _ := time.Parse(time.RFC3339, timeStringToParse)

	//// layoutISO := "2006-01-02"
	// layoutISO := "2006-01-02 00:00:00 "
	// date := "2020-6-30"
	// t, _ := time.Parse(layoutISO, date)

	// t := sTUrlCreateTime
	// unixt := t.Unix()
	// unixt++
	// unixTimeUTC := time.Unix(1405544146, 0)
	// fmt.Println(unixt) // 1999-12-31 00:00:00 +0000 UTC

	// dayTime := now.Truncate(24 * time.Hour)
	// t = t.Add(time.Second * 10)

	// var orderQty int
	// var orderCompletionPct int

	// var orderid []string

    		// // day := time.Now().Day()
		// const (
		// 	layoutISO = "2006-01-02"
		// 	layoutUS  = "January 2, 2006"
		// )
		// date := "1999-12-31"
		// t, _ := time.Parse(layoutISO, date)