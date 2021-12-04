package util

import "sync/atomic"

func InterlockIncrement(value *int32) int32 {
	return atomic.AddInt32(value, int32(1))
}

func InterlockDecrement(value *int32) int32 {
	return atomic.AddInt32(value, int32(-1))
}

func InterlockIncrement64(value *int64) int64 {
	return atomic.AddInt64(value, int64(1))
}

func InterlockDecrement64(value *int64) int64 {
	return atomic.AddInt64(value, int64(-1))
}

func InterlockedCompareExchange(value *int32, exchange int32, compare int32) bool {
	return atomic.CompareAndSwapInt32(value, compare, exchange)
}

func InterlockedCompareExchange64(value *int64, exchange int64, compare int64) bool {
	return atomic.CompareAndSwapInt64(value, compare, exchange)
}