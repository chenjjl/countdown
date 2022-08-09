package timeWheel

import (
	"container/list"
	"countdown/src/logger"
	"errors"
	"sync"
	"time"
)

var log = logger.GetLogger("timeWheel")

type littleHandTimeWheel struct {
	wheel
	c  chan *Event
	mu sync.Mutex
}

func NewLittleHandTimeWheel(tick time.Duration, wheelSize uint64) *littleHandTimeWheel {
	_tick := uint64(tick / time.Millisecond)
	if _tick <= 0 {
		panic(errors.New("little hand's tick must be greater than 0ms"))
	}
	if wheelSize <= 0 {
		panic(errors.New("little hand's size of timeWheel must be greater than 0"))
	}
	buckets := list.New()
	for i := uint64(0); i < wheelSize; i++ {
		bucket := NewLittleHandBucket()
		buckets.PushBack(bucket)
	}
	return &littleHandTimeWheel{
		wheel: wheel{
			tick:        _tick,
			wheelSize:   wheelSize,
			buckets:     buckets,
			bucketIndex: 0,
			curBucket:   buckets.Front(),
			head:        buckets.Front(),
		},
		c: make(chan *Event),
	}
}

func (t *littleHandTimeWheel) Start() {
	ticker := time.NewTicker(time.Duration(t.tick) * time.Millisecond)
	defer ticker.Stop()
	log.Infof("little hand's ticker has started, tick is %dms", t.tick)
	for {
		select {
		case <-ticker.C:
			t.doLookup()
		}
	}
}

func (t *littleHandTimeWheel) doLookup() {
	events, ok := t.Lookup()
	if ok {
		for _, event := range events {
			t.c <- event
		}
	}
}

func (t *littleHandTimeWheel) Add(event ...*Event) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	var err error
	for _, item := range event {
		e := item.Expiration / t.tick
		item.round = e / t.wheelSize
		index := (e + t.bucketIndex) % t.wheelSize
		_bucket := t.head
		for i := uint64(0); i < index; i++ {
			_bucket = _bucket.Next()
		}
		bucket := (_bucket.Value).(*littleHandBucket)
		if index == t.bucketIndex {
			item.curRound += 1
		}
		err = bucket.Add(item)
		if err != nil {
			log.Error("add event to bucket fail")
			log.Error(err)
		}
	}
	return err
}

func (t *littleHandTimeWheel) Lookup() ([]*Event, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.bucketIndex = (t.bucketIndex + 1) % t.wheelSize
	t.curBucket = t.curBucket.Next()
	// circle queue
	if t.curBucket == nil {
		t.curBucket = t.head
	}

	bucket := (t.curBucket.Value).(*littleHandBucket)
	events, err := bucket.Lookup()
	if err != nil {
		log.Error(err)
		return nil, false
	}
	if events == nil || len(events) == 0 {
		return nil, false
	}
	return events, true
}
