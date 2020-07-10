package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Print("Hello")
	usages := []string{"a", "b", "c", "d", "e"}

	var wg sync.WaitGroup

	fmt.Println("Running for loop…")
	retArray := []string{}

	for _, usage := range usages {
		wg.Add(1)
		go func(uu string) {
			defer wg.Done()

			fmt.Println("Insert " + uu)

			retArray = append(retArray, uu)
		}(usage)

	}

	wg.Wait()
	fmt.Println("Finished for loop")
	fmt.Println(retArray)
}

func main() {
	fmt.Print("Hello")
	usages := []string{"a", "b", "c", "d", "e"}

	var wg sync.WaitGroup

	fmt.Println("Running for loop…")
	retArray := []string{}
	//總共有多少筆，其實靠lens(usages)就知道，把wg.Add(1)移到for外面可以節省點資源
	wg.Add(len(usages))
	for _, usage := range usages {

		go func(uu string) {
			defer wg.Done()

			fmt.Println("Insert " + uu)

			retArray = append(retArray, uu)
		}(usage)

	}

	wg.Wait()
	fmt.Println("Finished for loop")
	fmt.Println(retArray)
}
