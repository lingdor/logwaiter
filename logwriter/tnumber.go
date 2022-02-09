package logwriter

import (
	"sync/atomic"
	"time"
)

// TNumber for second number, thread safe
type TNumber struct {
	second int64
	num    int64
}

func (t *TNumber) Incr() {
	second := time.Now().Unix()
	if atomic.LoadInt64(&t.second) == second {
		atomic.AddInt64(&t.num, 1)
	} else {
		atomic.SwapInt64(&t.num, 1)
		atomic.SwapInt64(&t.second, second)
	}
}

func (t *TNumber) Get() int64 {

	second := time.Now().Unix()
	if atomic.LoadInt64(&t.second) == second {
		return atomic.LoadInt64(&t.num)
	}
	return 0
}
