package util

import "math/rand/v2"

func RandInt(start, end int) int {
	return rand.IntN(end-start+1) + start
}
