package main

import (
	"fmt"
	"math/rand"
	"time"
)

func cakeMaker(kind string, number int, to chan<- string) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < number; i++ {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		to <- kind
	}
	close(to)

}

func main() {
	chocolateChan := make(chan string, 8)
	redVelvetChan := make(chan string, 8)

	go cakeMaker("chocolateChan", 4, chocolateChan)
	go cakeMaker("redVelvetChan", 3, redVelvetChan)

	moreChocolate := true
	moreRedVelvet := true

	var cake string

	for moreChocolate || moreRedVelvet {
		select {
		case cake, moreChocolate = <-chocolateChan:
			if moreChocolate {
				fmt.Printf("Got a cake from the first factory : %s\n", cake)
			}
		case cake, moreRedVelvet = <-redVelvetChan:
			if moreRedVelvet {
				fmt.Printf("Got a cake from the second factory : %s\n", cake)
			}
		case <-time.After(250 * time.Millisecond):
			fmt.Printf("Timed out")
			return
		}
	}
}
