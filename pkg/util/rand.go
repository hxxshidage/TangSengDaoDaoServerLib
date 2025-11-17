package util

import (
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	words = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	wLen  = len(words)
)

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
	defaultStrip  = 8
)

var (
	lockMask    = uint64(defaultStrip - 1)
	srcLocks    []*srcLock
	srcCounter  atomic.Uint64
	randCounter atomic.Uint64
)

type srcLock struct {
	mu  sync.Mutex
	src rand.Source
}

func init() {
	sli := make([]*srcLock, defaultStrip)
	for i := 0; i < defaultStrip; i++ {
		sli[i] = &srcLock{
			src: rand.NewSource(time.Now().UnixNano() + int64(srcCounter.Add(1))),
		}
	}

	srcLocks = sli
}

func RandStr(len uint16) string {
	appender := strings.Builder{}
	rLen := int(len)
	appender.Grow(rLen)

	sl := srcLocks[randCounter.Add(1)&lockMask]

	sl.mu.Lock()
	defer sl.mu.Unlock()

	src := sl.src

	for i, cache, remain := rLen-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < wLen {
			appender.WriteByte(words[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return appender.String()
}

func RandInt(start, end int) int {
	idx := randCounter.Add(1) & lockMask
	sl := srcLocks[idx]

	sl.mu.Lock()
	defer sl.mu.Unlock()

	n := end - start + 1
	if n&(n-1) == 0 {
		return start + int(sl.src.Int63()&int64(n-1))
	}

	_max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	var r int64
	for {
		r = sl.src.Int63()
		if r <= _max {
			break
		}
	}
	return start + int(r%int64(n))
}
