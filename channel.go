package main

import (
	"fmt"
	"time"
)

func goRoutineA(a <-chan int) {
	val := <-a
	fmt.Println("goRoutine A received the data", val)
}

func goRoutineB(b <-chan int) {
	val := <-b
	fmt.Println("goRoutine B received the data", val)
}

func main() {
	ch := make(chan int)
	go goRoutineA(ch)
	go goRoutineB(ch)
	ch <- 3
	time.Sleep(time.Second * 1)
}
