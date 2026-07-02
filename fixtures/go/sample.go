package sample

import "fmt"

// Simple function
func Simple() int {
	return 42
}

// Complex function with decision points
func ComplexFunc(x, y int) int {
	if x > 0 {
		for i := 0; i < x; i++ {
			if i%2 == 0 {
				y += i
			} else if i%3 == 0 {
				y -= i
			} else {
				y *= 2
			}
		}
	} else if x < 0 {
		for y > 0 {
			y--
		}
	} else {
		if y != 0 {
			y = x / y
		} else {
			y = 0
		}
	}
	return y
}

// Function with switch
func Classify(x int) string {
	switch {
	case x > 100:
		return "large"
	case x > 10:
		return "medium"
	case x > 0:
		return "small"
	default:
		return "negative"
	}
}

// Calculator
type Calculator struct{}

func (c Calculator) Compute(op string, a, b int) (int, error) {
	switch op {
	case "add":
		return a + b, nil
	case "sub":
		return a - b, nil
	case "mul":
		return a * b, nil
	case "div":
		if b != 0 {
			return a / b, nil
		}
		return 0, fmt.Errorf("division by zero")
	default:
		return 0, fmt.Errorf("unknown operation: %s", op)
	}
}
