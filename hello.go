package main

import (
	"fmt"
	"math"

	"github.com/axmty/go-playground/morestrings"
	"github.com/google/go-cmp/cmp"
)

func powIf(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	}
	return lim
}

func powMin(x, n, lim float64) float64 {
	return math.Min(lim, math.Pow(x, n))
}

func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}

func main() {
	fmt.Println(morestrings.ReverseRunes("Hello world!"))
	fmt.Println(cmp.Diff("Hello world", "Hello go"))

	// For loop
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)

	// "While" loop
	sum = 0
	for sum < 1000 {
		sum++
	}
	fmt.Println(sum)

	// Infinite loop
	// for {
	// }
}
