package main

import (
	"fmt"
	"io"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	fmt.Println("Hello World!")

	// Loops
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)

	// Optional Init and Post Statements
	sum = 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Println(sum)

	// Infinite loop
	for {
		break
	}

	// If
	fmt.Println(sqrt(2), sqrt(-4))

	// If with short statement
	fmt.Println(
		pow(3, 2, 10),
		pow(3, 3, 20),
	)
	// If-else scope
	fmt.Println(
		pow2(3, 2, 10),
		pow2(3, 3, 20),
	)

	// Switch
	fmt.Println("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux")
	default:
		// freebsd, opnbsd
		// plan9, windows...
		fmt.Printf("%s. \n", cases.Title(language.AmericanEnglish).String(os))
	}

	// Interesting - Switch evaluation order does not drop through cases
	// Interesting -- switch without condition is switch true, allowing clean long if-then-else chains
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}

	// Defer until completion of the surrounding function
	// Arguments are evaluated immediately, not lazily
	deferFunc()

	// Stacking defers is LIFO
	deferStack()

	// Go has pointers
	fmt.Println("Pointers")
	i, j := 42, 2701
	p := &i
	fmt.Println(*p)
	*p = 21
	fmt.Println(i)

	p = &j
	*p = *p / 37
	fmt.Println(j)

	// There is no pointer arithmetic

	// Structs
	fmt.Println(Vertex{1, 2})
	// Accessed through .
	v := Vertex{1, 2}
	v.X = 4
	fmt.Println(v.X)

	// No explicit dereference needed
	pv := &v
	// Explicit
	fmt.Println((*pv).X)
	// Implicit
	fmt.Println(pv.X)

	// Struct literals
	var (
		v1  = Vertex{1, 2}  // Has type Vertex
		v2  = Vertex{X: 1}  // Y:0 is implicit
		v3  = Vertex{}      // X:0 and Y:0 is implicit
		pv4 = &Vertex{1, 2} // has type *Vertex
	)
	fmt.Println(v1, v2, v3, pv4)

	// Arrays (size, type)
	var a [10]int
	a[0] = 1
	a[1] = 2
	fmt.Println(a)

	var b [2]string
	b[0] = "Hello"
	b[1] = "World"
	fmt.Println(b[0], b[1])
	fmt.Println(b)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)

	// Arrays are fixed, slices are dynamically sized
	// Follows syntax a[low : high] inclusive, exclusive
	var s []int = primes[1:4]
	fmt.Println(s)

	// Slices are references to arrays, like Rust
	names := [4]string{
		"John",
		"Paul",
		"George",
		"Ringo",
	}
	fmt.Println(names)

	names1 := names[0:2]
	names2 := names[1:3]
	fmt.Println(names1, names2)

	names2[0] = "XXX"
	fmt.Println(names1, names2)
	fmt.Println(names)

	// Slice literals are arrays without length, but creates slice referencing it (like Rust)
	q := []int{2, 3, 5, 7, 11, 13}
	fmt.Println(q)

	r := []bool{true, false, true, true, false, true}
	fmt.Println(r)

	sa := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
		{7, true},
		{11, false},
		{13, true},
	}
	fmt.Println(sa)

	// Slice defaults
	sb := []int{2, 3, 5, 7, 11, 13}
	sb = sb[1:4]
	fmt.Println(sb)

	sb = sb[:2]
	fmt.Println(s)

	sb = sb[1:]
	fmt.Println(sb)

	// For this array
	// 		var a [10]int
	// These are equivalent
	// 		a[0:10]
	// 		a[:10]
	// 		a[0:]
	// 		a[:]

	// Slice length and capacity
	sc := []int{2, 3, 5, 7, 11, 13}

	// Zero length
	sc = sc[:0]
	printSlice(sc)

	// Extend
	sc = sc[:4]
	printSlice(sc)

	// Drop first two values
	sc = sc[2:]
	printSlice(sc)

	// Extend beyond capacity
	// 		sc = sc[:8] // panic: runtime error: slice bounds out of range [:8] with capacity 4

	// Zero value of slice is nil (len 0 cap 0 and no underlying array)
	var sd []int
	fmt.Println(sd, len(sd), cap(sd))
	if sd == nil {
		fmt.Println("nil!")
	}

	// Built-in make function is used to create dynamically sliced collections, like arrays
	ds := make([]int, 5)
	printSlice2("ds", ds)

	// Capacity is third argument
	dsa := make([]int, 0, 5)
	printSlice2("dsa", dsa)

	dsab := dsa[:2]
	printSlice2("dsab", dsab)

	dsabc := dsab[2:5]
	printSlice2("dsabc", dsabc)

	// You can have slices of slices
	// Create a tic-tac-toe board
	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}
	// The players take turns.
	board[0][0] = "X"
	board[2][2] = "O"
	board[1][2] = "X"
	board[1][0] = "O"
	board[0][2] = "X"

	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}

	// You can append to slices
	var appending []int
	printSlice(appending)

	// Append works on nil slices
	appending = append(appending, 0)
	printSlice(appending)

	// The slice grows as needed
	appending = append(appending, 1)
	printSlice(appending)

	// We can add multiple elements at a time
	appending = append(appending, 2, 3, 4)
	printSlice(appending)

	// We have ranges!!
	var powaa = []int{1, 2, 4, 8, 16, 32, 64, 128}
	for i, v := range powaa {
		fmt.Printf("2**%d = %d\n", i, v)
	}

	// You can skip stuff using underscores (_)
	powb := make([]int, 10)
	for i := range powb {
		powb[i] = 1 << uint(i) // == 2**i
	}
	for _, value := range powb {
		fmt.Printf("%d\n", value)
	}

	// We got maps (map[key]Value)
	m := make(map[string]Coord)
	m["Bell Labs"] = Coord{40.68433, -74.39967}
	fmt.Println(m["Bell Labs"])

	// Map literals
	// 		They are struct literals, but keys are required
	n := map[string]Coord{
		"Bell Labs": Coord{ // Type is implicit, no need to define Coord
			40.68433, -74.39967,
		},
		"Google": {
			37.42202, -122.08408,
		},
	}
	fmt.Println(n)

	// We can mutate maps
	// Insert or update (n[key] = element)
	// Retrieve (element = n[key])
	// Delete (delete(n, key))
	// Test key is present with two-value assignment
	// element, ok = n[key]
	element, ok := n["Dell"]
	fmt.Println(element, ok)

	// Functions are values too!
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}
	fmt.Println(hypot(5, 12))

	fmt.Println(compute(hypot))
	fmt.Println(compute(math.Pow))

	// Interesting - Function closures
	fmt.Println("Closures")
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(
			pos(i),
			neg(-2*i),
		)
	}

	fmt.Println("Testing my closure understanding")
	// let's test the stateful stuff
	myAdder := adder()
	fmt.Println(myAdder(5))  // Expecting 5
	fmt.Println(myAdder(5))  // Expecting 10?
	fmt.Println(myAdder(-5)) // Expecting 5?
	// I am right. This is useful for counters etc.

	// Let's make a fibonacci closure
	fmt.Println("Fibonacci closure")
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}

	// Go does not have classes
	// But types can have methods
	// Methods are functions with a special receiver argument
	// Interesting, kind of like traits in rust?
	myV := MyV{3, 4}
	fmt.Println(myV.Abs())

	// Methods are functions. It just has a receiver argument
	fmt.Println((myV))

	// Methods can be defined on non-struct types too
	// You can not do this on types not defined in the same package as the method (such as on built-in types)
	fmt.Println("Method on non struct types")
	myf := MyFloat(-math.Sqrt2)
	fmt.Println(myf.Abs())

	// Pointer receivers
	// Everything is implicitly passed by value for Go?
	myV2 := MyV{3, 4}

	// Pointer indirection. We don't need to pass &v because it is implicit. this is why myV.Scale(5) works even though myV is a value.
	// It interprets it as (&myV).Scale(5) since the method has a pointer receiver
	myV2.Scale(10)

	fmt.Println(myV2.Abs())

	// Pointer indirection also happens in the reverse direction.
	// Functions that take a value argument must take a value of that specific type
	// var v Vertex
	// fmt.Println(AbsFunc(v)) // OK
	// fmt.Println(AbsFunc(&v)) // Compile error!

	// While methods with value receivers can take either a value or a pointer as the receiver when they are called, implicit dereference
	// var v Vertex
	// fmt.Println(v.Abs()) // OK
	// p := &v
	// fmt.Println(p.Abs()) // OK
	//
	// p.Abs() is interpreted as (*p).Abs()

	// Don't mix value and pointer receivers on a given type... Makes working with interfaces difficult? And more -- let's look at it
	var abser Abser
	ffff := MyFloat(-math.Sqrt2)
	vvvv := MyV{3, 4}

	abser = ffff
	abser = vvvv

	fmt.Println(abser.Abs())

	// Interface valeus with nil. If the concrete value inside the interface is nil, the method will be called with a nil received. No NPE!
	fmt.Println("Nil interface impl")
	var nilAbser Abser
	var myNilT *T
	nilAbser = myNilT
	fmt.Println(nilAbser.Abs())

	// Nil interface value is run-time error, no type inside tuple to indicate which concrete method to call
	//		var nilAbser2 Abser
	//		nilAbser2.Abs() // panic: runtime error: invalid memory address or nil pointer dereference

	// The empty interface
	// An empty interface may hold values of any type, because every type implements at least zero methods
	// They are used by code that handles values of unknown type. For example, fmt.Print takes any number of arguments of type interface{}
	// No type erasure!
	fmt.Println("Empty interface")
	var empty interface{}
	describe(empty)

	empty = 42
	describe(empty)

	empty = "hello"
	describe(empty)

	// Type assertion provides access to an interface value's underlying concrete type
	// We can test it with thee two-value return
	fmt.Println("Type assertion")
	var myI interface{} = "hello"
	myIS := myI.(string)
	fmt.Println(myIS)

	// Test
	myIS2, ok := myI.(string)
	fmt.Println(myIS2, ok)

	myIF, ok := myI.(float64)
	fmt.Println(myIF, ok)

	// myIF = myI.(float64) // panic
	// fmt.Println(myIF)

	// Type switches
	// It has to be named "type"
	switch ixxx := myI.(type) {
	case string:
		fmt.Println("It's a string", ixxx)
	case int:
		fmt.Println("It's an integer", ixxx)
	default:
		fmt.Println("No match", ixxx)
	}

	// Stringer interface is defined by the fmt package
	// type Stringer interface {
	// 	String() string
	// }

	// I've implemented Stringer on MyV
	stringerExample := MyV{10, 20}
	fmt.Println(stringerExample)

	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}

	// Readers
	reader := strings.NewReader("Hello, Reader!")
	byt := make([]byte, 8)
	for {
		n, err := reader.Read(byt)
		fmt.Printf("n = %v, err = %v b = %v\n", n, err, byt)
		fmt.Printf("b[:n] = %q\n", byt[:n])
		if err == io.EOF {
			break
		}
	}

	//reader.Validate(MyInfinityReader{})

	// Type parameters
	si := []int{10, 20, 15, -10}
	fmt.Println(Index(si, 15))

	ss := []string{"foo", "bar", "baz"}
	fmt.Println(Index(ss, "hello"))

	root := List[string]{}
	fmt.Println("Linked list")
	root.Add("5").Add("6").Add("7")
	fmt.Println(root)

	// Goroutines
	fmt.Println("Goroutines")

	go say("world")
	say("hello")

	// Channels
	fmt.Println("Channels")

	sssss := []int{7, 2, 8, -9, 4, 0}
	ch := make(chan int)

	go sumch(sssss[:len(sssss)/2], ch)
	go sumch(sssss[len(sssss)/2:], ch)
	x, y := <-ch, <-ch

	fmt.Println(x, y, x+y)

	fmt.Println("Buffered Channels")

	bufch := make(chan int, 2)
	bufch <- 1
	bufch <- 2
	fmt.Println(<-bufch)
	fmt.Println(<-bufch)

	// Range and close
	fmt.Println("Range and close")
	rch := make(chan int, 10)
	go fibonaccich(cap(rch), rch)
	for i := range rch {
		fmt.Println(i)
	}

	// Mutex
	fmt.Println("Mutex")
	c := SafeCounter{v: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go c.Inc("somekey")
	}
	time.Sleep(time.Second)
	fmt.Println(c.Value("somekey"))
}

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	c.v[key]++
	c.mu.Unlock()
}

func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v[key]
}

func fibonaccich(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func sumch(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

// Type parameters appear between square brackets, defining constraints which are interfaces
func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}

// Generic type. any is alias for interface{}
type List[T any] struct {
	next *List[T]
	val  T
}

func (l *List[T]) Add(element T) *List[T] {
	l.next = &List[T]{val: element}
	return l.next
}

func (l List[T]) String() string {
	s := make([]string, 0)
	for {
		s = append(s, fmt.Sprint(l.val))
		if l.next == nil {
			break
		}
		l = *l.next
	}
	return strings.Join(s, " ")
}

type IPAddr [4]byte

type MyInfinityReader struct{}

func (r MyInfinityReader) Read(b []byte) (int, error) {
	for i := range b {
		b[i] = 'A'
	}
	return len(b), nil
}

func (v IPAddr) String() string {
	s := make([]string, 0, len(v))
	for _, i := range v {
		s = append(s, strconv.Itoa(int(i)))
	}
	return strings.Join(s, ".")
}

// Interfaces are implicitly implemented on any supported type
// Under the hood, interface values can be thought of as tuples of a value and a concrete type
// (value, type)
// It holds a value of a specific underlying concrete type
// Calling a method on an interface value executes the method of the same signature on its underlying type -- Very interesting
type Abser interface {
	Abs() float64
}

type T struct {
	S string
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func (t *T) Abs() float64 {
	if t == nil {
		fmt.Println("<nil>")
		return 0
	}
	return 1
}

func (v *MyV) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
	// v.X, v.Y = v.X * f, v.Y * f // This is equivalent
}

type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	// Implicit conversion to float64 -- cool
	return float64(f)
}

type MyV struct {
	X, Y float64
}

func (v MyV) String() string {
	return fmt.Sprintf("X = %v Y = %v\n", v.X, v.Y)
}

// func, receiver (received type), the function name, return type
// Abs has a receiver of type MyV named v
func (v MyV) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Very strange :) see above
func Abs(v MyV) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func fibonacci() func() int {
	f2, f1 := 0, 1
	return func() int {
		f := f2
		f2, f1 = f1, f+f1
		return f
	}
}

func alsoFibonacci() func() int {
	f2 := 0
	f1 := 1
	return func() int {
		f := f2
		f2 = f1
		f1 = f + f1
		return f
	}
}

func adder() func(int) int {
	// The function is bound to the variable
	// Very funky, not used to this.. kind of like creating a class with a method, that then acts on internal state
	// Stateful function?
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap%d %v\n", len(s), cap(s), s)
}

// No method overloading in Go
func printSlice2(s string, x []int) {
	fmt.Printf("%s len=%d cap=%d %v\n", s, len(x), cap(x), x)
}

func deferFunc() {
	defer fmt.Println("world")
	fmt.Println("Hello, ")
}

func deferStack() {
	fmt.Println("Counting")
	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}
	fmt.Println("Done")
}

// if statements
func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}

func pow(x, n, lim float64) float64 {
	// Interesting
	if v := math.Pow(x, n); v < lim {
		return v
	}
	return lim
}

// If and else
func pow2(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	} else {
		fmt.Printf("%g >= %g\n", v, lim)
	}
	// Can't use v here
	return lim
}

type Vertex struct {
	X int
	Y int
}

type Coord struct {
	Lat, Long float64
}

// package tree

// import (
// 	"fmt"

// 	"golang.org/x/tour/tree"
// )

// // Walk walks the tree t sending all values
// // from the tree to the channel ch.
// func Walk(t *tree.Tree, ch chan int) {
// 	defer close(ch) // <- closes the channel when this function returns
// 	var walk func(t *tree.Tree)
// 	walk = func(t *tree.Tree) {
// 		if t == nil {
// 			return
// 		}
// 		walk(t.Left)
// 		ch <- t.Value
// 		walk(t.Right)
// 	}
// 	walk(t)
// }

// // Same determines whether the trees
// // t1 and t2 contain the same values.
// func Same(t1, t2 *tree.Tree) bool {
// 	ch1 := make(chan int)
// 	ch2 := make(chan int)

// 	go Walk(t1, ch1)
// 	go Walk(t2, ch2)

// 	for k := range ch1 {
// 		select {
// 		case g := <-ch2:
// 			if k != g {
// 				return false
// 			}
// 		default:
// 			break
// 		}
// 	}
// 	return true

// }

// func aaa() {
// 	fmt.Println(Same(tree.New(1), tree.New(1)))
// 	fmt.Println(Same(tree.New(1), tree.New(2)))
// }
