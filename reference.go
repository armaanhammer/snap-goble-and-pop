package main

import (
	"fmt"
)

func main() {
	var x int = 10
	var y int = 10

	var z int = x + y

	sum := x + y

	fmt.Println(z)
	fmt.Println(sum)

	if x < y {
		fmt.Println(x, "is less than", y)
	} else if x > y {
		fmt.Println(x, "is greater than", y)
	} else {
		fmt.Println(x, "is equal to", y)
	}


	// Arrays
	var a [5]int
	a[2]=100 // arrays are zero-indexed

	fmt.Println(a)

	bottles := [5]int{100,99,98,97,96}

	fmt.Println(bottles)

	// Slices https://tour.golang.org/moretypes/7
	on_the_wall := []int{100,99,98,97,96}
	on_the_wall = append(on_the_wall,95) // append returns a new slice

	fmt.Println(on_the_wall)


	// Maps
	// map[key]value
	verticies := make(map[string]int)

	verticies["line"] = 0
	verticies["triangle"] = 3
	verticies["square"] = 4
	verticies["cube"] = 8

	delete(verticies, "line")

	fmt.Println(verticies)
	fmt.Println(verticies["triangle"])


	// Loops
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// while loop
	i := -4
	for i < 5 {
		fmt.Println(i)
		i++
	}
}
