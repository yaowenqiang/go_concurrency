package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go printEven(i, &wg)
	}

	// wg.Wait()
	// for i := 0; i < 10; i++ {
	// 	go fmt.Printf("Goroutine numbeer: %d\n", i)
	// }
	// fmt.Println("loop finished")
	// time.Sleep(1 * time.Second)

	val := 0
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		// go increment(&val, &wg)
		increment(&val, &wg)
	}

	wg.Wait()

	fmt.Printf("Final value was %d\n", val)
}

func printEven(x int, wg *sync.WaitGroup) {
	if x%2 == 0 {
		fmt.Printf("%d is even\n", x)
	}

	wg.Done()
}
func increment(ptr *int, wg *sync.WaitGroup) {
	i := *ptr
	fmt.Printf("i is %d\n", i)
	*ptr = i + 1
	wg.Done()
}
