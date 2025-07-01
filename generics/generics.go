package main

import "fmt"

func generics_circunference[r int | float64](radius r) {
	c := 2 * 3 * radius
	fmt.Println("The circunferece is: ", c)
}

// Parametrized Types
type Radius interface {
	int64 | int8 | float64
}

func generics_circunference2 [R Radius](radius R) {
	var c R
	c = 2 * 3 * radius
	fmt.Println("The circunferece is: ", c)
}

// Constraints GO

type Number interface {
	int | float64 | complex64
}

func Add[T Number](a, b T) T {
	return a + b
}

func main() {
	var radius int = 8
	var radius2 float64 = 9.5

	generics_circunference(radius)
	generics_circunference2(radius2)
}