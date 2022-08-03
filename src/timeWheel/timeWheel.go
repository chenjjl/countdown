package timeWheel

import (
	"container/list"
	"countdown/src/logger"
	"errors"
	"sync"
	"time"
)

var log = logger.GetLogger("timeWheel")

type TimeWheel struct {
	tick        uint64
	wheelSize   uint64
	buckets     *list.List
	bucketIndex uint64
	head        *list.Element
	mu          sync.Mutex
}

func NewTimeWheel(tick time.Duration, wheelSize uint64) *TimeWheel {
	_tick := uint64(tick / time.Millisecond)
	if _tick <= 0 {
		panic(errors.New("tick must be greater than 0ms"))
	}
	if wheelSize <= 0 {
		panic(errors.New("size of timeWheel must be greater than 0"))
	}
	buckets := list.New()
	for i := uint64(0); i < wheelSize; i++ {
		bucket := NewBucket()
		buckets.PushBack(bucket)
	}
	return &TimeWheel{
		tick:        _tick,
		wheelSize:   wheelSize,
		buckets:     buckets,
		bucketIndex: 1,
		head:        buckets.Front(),
	}
}

func (t *TimeWheel) Add(event *Event) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	e := (event.expiration + t.bucketIndex*t.tick) / t.tick
	event.round = e / t.wheelSize
	index := e % t.wheelSize
	_bucket := t.head
	for i := uint64(0); i < index-1; i++ {
		_bucket = _bucket.Next()
	}
	bucket := (_bucket.Value).(*bucket)
	err := bucket.Add(event)
	if err != nil {
		return err
	}
	return nil
}
