package common

import (
	"fmt"
)

func ZigZag(n int64) int64 {
	return (n << 1) ^ (n >> 63)
}

func Decimal2Binary(n int64) string {
	var s string
	for {
		s = fmt.Sprintf("%d", n%2) + s

		n = n / 2
		if n == 0 {
			break
		}
	}
	return s
}
