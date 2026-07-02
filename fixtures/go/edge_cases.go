package sample

import "fmt"

// Simple function
func simple() int {
	return 42
}

// Function with error handling
func process(data []int) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data")
	}

	sum := 0
	for _, v := range data {
		if v < 0 {
			continue
		}
		if v > 100 {
			break
		}
		sum += v
	}
	return sum, nil
}

// Complex function with switch
func classify(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

// Function with defer
func cleanup() {
	defer fmt.Println("cleanup done")
	fmt.Println("working")
}

// Method with receiver
type Counter struct {
	count int
}

func (c *Counter) Increment() {
	c.count++
}

func (c *Counter) Get() int {
	return c.count
}

// Anonymous function
var transform = func(x int) int {
	return x * 2
}

// Function with multiple returns
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}
