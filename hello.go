package main

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/axmty/go-playground/morestrings"
	"github.com/google/go-cmp/cmp"
)

// Vertex is a simple vector struct
type Vertex struct {
	X int
	Y int
}

// Sqrt computes the square root of x with a precision of epsilon
func Sqrt(x, epsilon float64) float64 {
	z := 1.0
	for math.Abs(z*z-x) > epsilon {
		z -= (z*z - x) / (2 * z)
		fmt.Println(z, z*z)
	}
	return z
}

func powIf(x, n, lim float64) float64 {
	v := math.Pow(x, n)
	if v < lim {
		return v
	}
	fmt.Printf("%g >= %g\n", v, lim)
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

func deferPrint(s string) {
	defer fmt.Println(s)
	fmt.Println("Deferred string:")
}

func deferPrints(arr []string) {
	fmt.Println("Deferred strings:")
	for _, s := range arr {
		// Each defer statement push the target statement onto a stack
		// (Last In First Out)
		defer fmt.Println(s)
	}
}

func playPointers() {
	i := 1
	p := &i
	i = 2
	fmt.Println(i, p, *p)
}

func playStructs() {
	v := Vertex{1, 2}
	v.X = 4
	p := &v
	fmt.Println(v.X)
	v.X = 5
	fmt.Println((*p).X)
	fmt.Println(p.X) // No need to write explicit dereference (*p).X

	v1 := Vertex{1, 2}  // X = 1, Y = 2
	v2 := Vertex{X: 1}  // Y = 0
	v3 := Vertex{}      // X = 0, Y = 0
	p1 := &Vertex{1, 2} // Type *Vertex

	fmt.Println(v1, v2, v3, p1)
}

func playArrays() {
	var arr1 [2]string // Length 2 is part of the type
	arr2 := [6]int{1, 2, 3}
	fmt.Println(arr1)
	fmt.Println(arr2)

	var slice1 []int = arr2[1:4] // Slice, dynamically-sized view of arr2
	var slice2 []int = arr2[1:4]
	slice3 := arr2[2:]
	slice4 := arr2[:]
	slice5 := arr2[:3]
	fmt.Println(slice1, slice2, slice3, slice4, slice5)
	slice1[0] = 100

	// Modifying slice1 element modifies the underlying array,
	// and other slices that share the same underlying array
	// will see those changes (they are just views on arrays)
	fmt.Println(slice1, slice2, arr2)

	// Slice literal, creates the underlying array [3]int
	// and builds the slice that references it
	q := []int{2, 3, 4}

	s := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
	}

	fmt.Println(q, s)

	// Length is 3, capacity is 3 too, because the slice reference
	// all the elements, from the first one.
	a := []int{1, 2, 3}
	fmt.Printf("len=%d cap=%d %v\n", len(a), cap(a), a) // 3 3 [1 2 3]

	// Lenth is 1, but capacity is 2: 2 elements in the underlying array
	// (2 and 3), counting from the first element (2).
	a = a[1:2]
	fmt.Printf("len=%d cap=%d %v\n", len(a), cap(a), a) // 1 2 [2]

	// Error: capacity of a is now 2, so cannot extend the length to 3.
	// It would have no sense, because a is [2] from the previous assignment,
	// and the underlying arr is [1 2 3].
	// Slicing slice [2] with [:3] would generate a
	// "slice bounds out of range" error.
	// a = a[:3]

	// OK: slicing slice [2] with [:2] and with the underlying array [1 2 3]
	// will extend with the element 3.
	a = a[:2]
	fmt.Printf("len=%d cap=%d %v\n", len(a), cap(a), a) // 2 2 [2 3]

	// ==> Capacity is a guarantee that we will not extend beyond
	// the underlying array elements.
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

	fmt.Println(Sqrt(3, 0.001))

	// Simple switch conditions
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		fmt.Printf("%s.\n", os)
	}

	// More complex switch conditions
	switch today := time.Now().Weekday(); time.Saturday {
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}

	// Switch with no condition
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}

	deferPrint("Hello!!!")
	deferPrints([]string{"s1", "s2", "s3"})

	playPointers()
	playStructs()
	playArrays()
}
