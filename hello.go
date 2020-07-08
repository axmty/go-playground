package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/axmty/go-playground/morestrings"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/tour/reader"
	"golang.org/x/tour/wc"
)

// Vertex is a simple vector struct
type Vertex struct {
	X float64
	Y float64
}

type person struct {
	Name string
	Age  int
}

type myFloat float64

type abser interface {
	abs() float64
}

type ipAddr [4]byte

func (addr *ipAddr) String() string {
	return fmt.Sprintf(
		"%v.%v.%v.%v", addr[0], addr[1], addr[2], addr[3])
}

func (p person) String() string {
	return fmt.Sprintf("%v (%v years)", p.Name, p.Age)
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

type errNegativeSqrt float64

func (e errNegativeSqrt) Error() string {
	// If we just use e instead of float64(e), infinite loop.
	// Because it will print using the Error() method on errNegativeSqrt,
	// since e implements the interface error.
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func sqrtError(x float64) (float64, error) {
	if x < 0 {
		return 0, errNegativeSqrt(x)
	}
	return math.Sqrt(x), nil
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
	// "slice bounds out of range" error because extending is impossible.
	// a = a[:3]

	// OK: slicing slice [2] with [:2] and with the underlying array [1 2 3]
	// will extend with the element 3.
	a = a[:2]
	fmt.Printf("len=%d cap=%d %v\n", len(a), cap(a), a) // 2 2 [2 3]

	// ==> Capacity is a guarantee that we will not extend beyond
	// the underlying array elements.
	var nilSlice []int = nil
	fmt.Println(nilSlice, len(nilSlice), cap(nilSlice)) // [] 0 0
	if nilSlice == nil {
		fmt.Println("nil slice!")
	}
}

func makeSlice() {
	a := make([]int, 5)            // len(a) = 5, cap(a) = 5
	fmt.Println(a, len(a), cap(a)) // [0 0 0 0 0] 5 5

	b := make([]int, 0, 5)         // len(b) = 0, cap(b) = 5
	fmt.Println(b, len(b), cap(b)) // [] 0 5
	b = b[:cap(b)]                 // len(b) = 5, cap(b) = 5
	fmt.Println(b, len(b), cap(b)) // [0 0 0 0 0] 5 5
	b = b[:cap(b)-1]               // len(b) = 4, cap(b) = 5
	fmt.Println(b, len(b), cap(b)) // [0 0 0 0] 4 5
	b = b[1:]                      // len(b) = 3, cap(b) = 4
	fmt.Println(b, len(b), cap(b)) // [0 0 0] 3 4

	// Compile error: length cannot be larger than capacity
	// c := make([]int, 5, 3)
}

func appendToSlice() {
	arr := [3]int{1, 2, 3}
	s := arr[:2]
	fmt.Printf("arr: %v, s: %v\n", arr, s) // arr: [1 2 3], s: [1 2]

	fmt.Println("appending 4 and 5 to slice s")
	// Building a new slice with append() does not modify the original slice.
	v := append(s, 4, 5)
	fmt.Printf("arr: %v, s: %v, v: %v\n", arr, s, v) // arr: [1 2 3], s: [1 2 4 5]
}

func ranges() {
	// i, v = index, element copy
	for i, v := range []int{1, 2, 4, 8, 16} {
		fmt.Printf("2**%d = %d\n", i, v)
	}

	// Omitting index
	for _, v := range []int{1, 2, 4} {
		fmt.Printf("Value %d\n", v)
	}

	// Omitting value
	for i := range []int{1, 2, 3} {
		fmt.Printf("At index %d\n", i)
	}
}

func playMaps() {
	m := make(map[string]Vertex)
	m["ABC"] = Vertex{
		4, -4,
	}
	fmt.Println(m["AB"])
	fmt.Println(m["ABC"])

	// Using map litteral.
	ages := map[string]int{
		"Alex": 28,
		"Nico": 31,
	}
	fmt.Println(ages) // map[Alex:28 Nico:31]

	// Omitting the type name.
	m2 := map[string]Vertex{
		"ABC": {4, -4},
	}
	fmt.Println(m2) // map[ABC:{4 -4}]

	delete(m, "ABC")
	fmt.Println(m) // map[]

	elem, ok := m["ABC"]
	fmt.Println(ok, elem) // false {0 0}

	elem, ok = m2["ABC"]
	fmt.Println(ok, elem) // true {4 -4}
}

func wordCount(s string) map[string]int {
	m := make(map[string]int)
	words := strings.Fields(s)
	for _, w := range words {
		_, ok := m[w]
		if ok {
			m[w]++
		} else {
			m[w] = 1
		}
	}
	return m
}

func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func fibonacci() func() int {
	now, next := 0, 1
	return func() int {
		next = now + next
		now = next - now
		return next - now
	}
}

func playFuncs() {
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}
	fmt.Println(hypot(5, 12))

	fmt.Println(compute(hypot))
	fmt.Println(compute(math.Pow))

	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(
			pos(i),
			neg(-2*i),
		)
	}

	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

func (v *Vertex) abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (f *myFloat) abs() float64 {
	if f == nil {
		return 0
	} else if *f < 0 {
		return float64(-*f)
	}
	return float64(*f)
}

func (v *Vertex) scale(f float64) {
	v.X *= f
	v.Y *= f
}

func scale(v *Vertex, f float64) {
	v.X *= f
	v.Y *= f
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func playInterfaces() {
	var a abser
	f := myFloat(34)
	a = &f
	a = &Vertex{1, 2}

	fmt.Println(a.abs())

	var nilF *myFloat
	a = nilF
	fmt.Println(a) // <nil>

	// Error, nilF pointer is nil, and no nil pointer handling inside
	// *myFloat.abs method (we still know the concrete type thanks to
	// the statement 'var nilF *myFloat', see next).
	// This can be fixed with adding nil handling.
	// fmt.Println(a.abs())

	// nilF pointer is nil and the underlying value of
	// the a interface value is nil, but there is a concrete type.
	fmt.Printf("(%v, %T)\n", a, a) // (<nil>, *main.myFloat)

	a = nil
	fmt.Printf("(%v, %T)\n", a, a) // (<nil>, <nil>)
	// Error, the a interface value is nil, and there is no underlying type.
	// How could we know which method to call then?
	// fmt.Println(a.abs())

	var i interface{}
	describe(i) // (<nil>, <nil>)

	var x *Vertex
	i = x
	describe(i) // (<nil>, *main.Vertex)

	i = 42
	describe(i) // (42, int)

	i = "hello"
	describe(i) // (hello, string)
}

func playTypeAssertions() {
	var i interface{} = "hello"

	s := i.(string)
	fmt.Println(s) // hello

	s, ok := i.(string)
	fmt.Println(s, ok) // hello, true

	// panic: interface conversion: interface {} is string, not float64.
	// f := i.(float64)

	// No panic.
	f, ok := i.(float64)
	fmt.Println(f, ok) // 0 false (0 is float64 zero value)

	fn := func(i interface{}) {
		switch v := i.(type) {
		case int:
			fmt.Printf("Twice %v is %v\n", v, v*2)
		case string:
			fmt.Printf("%q is %v bytes long\n", v, len(v))
		default:
			fmt.Printf("I don't know about type %T\n", v)
		}
	}

	fn(34)
	fn("Hello world!")
	fn(true)
}

func playStringer() {
	a := person{"Arthur Dent", 42}
	z := person{"Zaphod Beeblebrox", 9001}
	fmt.Println(a, z)

	addr := ipAddr{127, 0, 0, 1}
	fmt.Println(&addr)
}

type myError struct {
	When time.Time
	What string
}

func (e *myError) Error() string {
	return fmt.Sprintf("at %v, %s", e.When, e.What)
}

func runError() error {
	return &myError{
		time.Now(),
		"it didn't work",
	}
}

func playErrors() {
	if err := runError(); err != nil {
		fmt.Println(err)
	}
}

type myReader struct{}

func (r myReader) Read(b []byte) (int, error) {
	for i := 0; i < len(b); i++ {
		b[i] = 'A'
	}
	return len(b), nil
}

func playReaders() {
	r := strings.NewReader("漢字!")
	b := make([]byte, 3)
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}

	fmt.Printf("%q\n", [3]byte{230, 188, 162}) // 漢 (3 bytes for this char)

	reader.Validate(myReader{})
}

type rot13Reader struct {
	r io.Reader
}

func (rr rot13Reader) Read(b []byte) (int, error) {
	n, err := rr.r.Read(b)
	for i := 0; i < len(b); i++ {
		c := b[i]
		if c >= 'A' && c <= 'Z' {
			b[i] = 'A' + (c-'A'+13)%26
		} else if c >= 'a' && c <= 'z' {
			b[i] = 'a' + (c-'a'+13)%26
		}
	}
	return n, err
}

func playRot13Reader() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
	fmt.Println()
}

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func fibo(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		fmt.Printf("Fibo step %v\n", i+1)
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func fibo2Chans(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func defaultSelect() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func playGoroutines() {
	go say("world")
	say("hello")

	s := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)

	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c

	fmt.Println(x, y, x+y)

	c = make(chan int)
	go fibo(10, c)
	for i := range c {
		fmt.Println(i)
	}

	c = make(chan int)
	quit := make(chan int)

	fmt.Println("start fibo with 2 channels")
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibo2Chans(c, quit)

	defaultSelect()

	fmt.Println("end method playGoroutines")
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
	makeSlice()
	appendToSlice()
	ranges()
	playMaps()
	wc.Test(wordCount)
	playFuncs()

	v := Vertex{3, 4}
	scale(&v, 10)
	v.scale(10)
	fmt.Println(v.abs())

	playInterfaces()
	playTypeAssertions()
	playStringer()
	playErrors()

	sqrtError(2)
	_, err := sqrtError(-2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	playReaders()
	playRot13Reader()
	playGoroutines()
}
