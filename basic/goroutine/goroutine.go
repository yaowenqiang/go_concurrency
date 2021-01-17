package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		go fmt.Printf("Goroutine numbeer: %d\n", i)
	}
	fmt.Println("loop finished")
	time.Sleep(1 * time.Second)
}
