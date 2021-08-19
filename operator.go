package main

import "fmt"

type T struct {
	x int
	y *int
}

var TestArray2 = [...]float32{.03, .02}

func main() {
	var t T
	p := &t.x
	fmt.Printf("%T\n", p)
	*p++
	*p--

	t.y = p
	a := *t.y
	fmt.Printf("%T\n", a)

	const M = 2
	fmt.Printf("%T\n", M)
	var _ = 1.0 << M
	var N = 2.0
	//var S = 1.0 << N
	fmt.Printf("%T\n", N)

}
