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
	curBucket   *list.Element
	bucketIndex uint64
	head        *list.Element
	c           chan *Event
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
		curBucket:   buckets.Front(),
		head:        buckets.Front(),
		c:           make(chan *Event),
	}
}

func (t *TimeWheel) Start() {
	ticker := time.NewTicker(time.Duration(t.tick) * time.Millisecond)
	defer ticker.Stop()
	log.Infof("ticker has started, tick is %dms", t.tick)
	for {
		select {
		case <-ticker.C:
			event, ok := t.Lookup()
			if ok {
				t.c <- event
			}
		}
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

func (t *TimeWheel) Lookup() (*Event, bool) {
	t.bucketIndex += 1
	t.curBucket = t.curBucket.Next()
	// circle queue
	if t.curBucket == nil {
		t.curBucket = t.head
	}
	bucket := (t.curBucket.Value).(*bucket)
	event, err := bucket.Lookup()
	if err != nil {
		log.Error(err)
		return nil, false
	}
	if event == nil {
		return nil, false
	}
	return event, true
}
