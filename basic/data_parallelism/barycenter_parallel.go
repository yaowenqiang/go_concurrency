package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// MassPoint struct
type MassPoint struct {
	x, y, z, mass int
}

func stringToPointAsync(s string, c chan<- MassPoint, wg *sync.WaitGroup) {
	defer wg.Done()

	var newMassPoint MassPoint

	_, err := fmt.Sscanf(s, "%d:%d:%d:%d", &newMassPoint.x, &newMassPoint.y, &newMassPoint.z, &newMassPoint.mass)

	if err != nil {
		panic(err)
	}

	c <- newMassPoint

}

func avgMassPointsWeightedAsync(a, b MassPoint, c chan<- MassPoint) {
	c <- avgMassPointsWeighted(a, b)
}

func addMassPoints(a, b MassPoint) MassPoint {
	return MassPoint{
		a.x + b.x,
		a.y + b.y,
		a.z + b.z,
		a.mass + b.mass,
	}
}

func avgMassPoints(a, b MassPoint) MassPoint {
	sum := addMassPoints(a, b)
	return MassPoint{
		sum.x / 2,
		sum.y / 2,
		sum.z / 2,
		sum.mass / 2,
	}
}

func toWeihtedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		a.x * a.mass,
		a.y * a.mass,
		a.y * a.mass,
		a.mass,
	}
}

func fromWeihtedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		a.x / a.mass,
		a.y / a.mass,
		a.y / a.mass,
		a.mass,
	}
}

func avgMassPointsWeighted(a, b MassPoint) MassPoint {
	aWeighted := toWeihtedSubspace(a)
	bWeighted := toWeihtedSubspace(b)

	return fromWeihtedSubspace(avgMassPoints(aWeighted, bWeighted))
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func closeFile(f *os.File) {
	err := f.Close()
	handle(err)
}
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Incorrect number of arguments.")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	handle(err)

	defer closeFile(file)

	var massPoints []MassPoint
	startLoading := time.Now()

	r := bufio.NewReader(file)
	massPointsChan := make(chan MassPoint, 128)

	var wg sync.WaitGroup

	for {
		str, err := r.ReadString('\n')
		if len(str) == 0 || err != nil {
			break
		}
		wg.Add(1)
		go stringToPointAsync(str, massPointsChan, &wg)
	}

	synChan := make(chan bool)

	go func() { wg.Wait(); synChan <- false }()

	run := true

	for run || len(massPointsChan) > 0 {
		select {
		case value := <-massPointsChan:
			massPoints = append(massPoints, value)
		case _ = <-synChan:
			run = false
		}
	}

	// for {
	// 	var newMassPoint MassPoint
	// 	_, err := fmt.Fscanf(file, "%d:%d:%d:%d", &newMassPoint.x, &newMassPoint.y, &newMassPoint.z, &newMassPoint.mass)
	// 	if err == io.EOF {
	// 		break
	// 	} else if err != nil {
	// 		handle(err)
	// 		// continue
	// 	}

	// 	massPoints = append(massPoints, newMassPoint)
	// }

	fmt.Printf("Loaded %d values from file in %s.\n", len(massPoints), time.Since(startLoading))

	if len(massPoints) <= 1 {
		handle(errors.New("insufficient values"))
	}

	c := make(chan MassPoint, len(massPoints)/2)

	startCalculation := time.Now()

	for len(massPoints) > 1 {
		var newMasPoints []MassPoint

		goroutines := 0

		for i := 0; i < len(massPoints)-1; i += 2 {
			go avgMassPointsWeightedAsync(massPoints[i], massPoints[i+1], c)
			goroutines++
		}

		for i := 0; i < goroutines; i++ {
			newMasPoints = append(newMasPoints, <-c)
		}

		if len(massPoints)%2 != 0 {
			newMasPoints = append(newMasPoints, massPoints[len(massPoints)-1])
		}

		massPoints = newMasPoints
	}

	systemAverage := massPoints[0]

	fmt.Printf("system barycenter is at (%d, %d, %d) and the system's mass is %d.\n",
		systemAverage.x,
		systemAverage.y,
		systemAverage.z,
		systemAverage.mass,
	)

	fmt.Printf("Calculation took %s.\n", time.Since(startCalculation))

	// if len(os.Args) < 2 {
	// 	os.Exit(1)
	// }
	// nbodies, err := strconv.Atoi(os.Args[1])
	// if err != nil {
	// 	os.Exit(1)
	// }
	// rand.Seed(time.Now().Unix())

	// posMax := 500
	// massMax := 5

	// for i := 0; i < nbodies; i++ {
	// 	posX := rand.Intn(posMax*2) - posMax
	// 	posy := rand.Intn(posMax*2) - posMax
	// 	posz := rand.Intn(posMax*2) - posMax

	// 	mass := rand.Intn(massMax-1) + 1

	// 	fmt.Printf("%d:%d:%d:%d\n", posX, posy, posz, mass)
	// }

}
