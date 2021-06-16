package main

import "fmt"

func main() {
	ch := make(chan int, 5)
	chs := make(chan string, 5)

	ch <- 123

	select {
	case msg := <-ch:
		fmt.Println("redeived message", msg)
	case msgs := <-chs:
		fmt.Println("redeived message", msgs)
	default:
		fmt.Println("no message received")
	}

}
