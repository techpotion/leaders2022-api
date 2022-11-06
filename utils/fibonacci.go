package utils

const (
	firstElemOffset  = 1
	secondElemOffset = 2
)

func FibonacciRecursion(n int) int {
	if n <= 1 {
		return n
	}

	return FibonacciRecursion(n-firstElemOffset) + FibonacciRecursion(n-secondElemOffset)
}
