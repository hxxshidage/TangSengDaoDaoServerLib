package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkRandInt(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			RandInt(1, 999)
		}
	})
}

func BenchmarkRandStr(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			RandStr(32)
		}
	})
}

func TestRandStr(t *testing.T) {
	num := 100_000_000
	m := make(map[string]struct{}, num)
	for i := 0; i < num; i++ {
		m[RandStr(12)] = struct{}{}
	}

	assert.True(t, len(m) == num)
}
