package logwriter

import "sync/atomic"

type inc int32

func (i *inc) Add() {
	atomic.AddInt32((*int32)(i), 1)
}

func (i *inc) Get() int32 {
	return atomic.LoadInt32((*int32)(i))
}

func (i *inc) Reset() {
	atomic.StoreInt32((*int32)(i), 0)
}
