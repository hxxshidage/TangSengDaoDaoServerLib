package guc

import (
	"hash/crc32"
	"math"
	"sync"
)

type Action[T any] func() (T, error)

type DcAction[T any] func() (T, bool, error)

// DcSegLocker dc mean double check
type DcSegLocker[T any] interface {
	DoWithDc(string, DcAction[T], Action[T]) (T, error)
}

type muLocker[T any] struct {
	locks []*sync.Mutex
	mask  uint32
}

func (l *muLocker[T]) DoWithDc(key string, dc DcAction[T], a Action[T]) (T, error) {
	lock := l.getLock(key)
	lock.Lock()
	defer lock.Unlock()

	if v, ok, e := dc(); e != nil {
		var zero T
		return zero, e
	} else if ok {
		return v, nil
	}

	return a()
}

func NewSegMuLocker[T any](strip int) DcSegLocker[T] {
	strip = CorrectStrip(strip)

	locks := make([]*sync.Mutex, strip)
	for i := 0; i < strip; i++ {
		locks[i] = new(sync.Mutex)
	}

	return &muLocker[T]{
		locks: locks,
		mask:  uint32(strip - 1),
	}
}

func (l *muLocker[T]) getLock(key string) *sync.Mutex {
	idx := crc32.ChecksumIEEE([]byte(key)) & l.mask
	return l.locks[idx]
}

const (
	defaultStrip = 4
)

func CorrectStrip(strip int) int {
	if strip == 0 {
		return defaultStrip
	}

	return 1 << uint(math.Log2(float64(strip-1))+1)
}
