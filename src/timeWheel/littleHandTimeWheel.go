package timeWheel

import (
	"container/list"
	"countdown/src/event"
	"countdown/src/storage"
	"errors"
	"sync"
	"time"
)

type littleHandTimeWheel struct {
	wheel
	lhFile      *storage.LhFile
	eventIdFile *storage.IdFile
	tickRound   uint64
	c           chan *event.Event
	mu          sync.Mutex
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
	lhFile, err := storage.CreateLhFile()
	if err != nil {
		panic(err)
	}
	eventIdFile, err := storage.CreateEventIdFile()
	if err != nil {
		panic(err)
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
		lhFile:      lhFile,
		eventIdFile: eventIdFile,
		tickRound:   0,
		c:           make(chan *event.Event),
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
			err := t.eventIdFile.Add(event)
			if err != nil {
				log.Errorf("failed to add event Id %s to file %s", event.Id, t.eventIdFile)
			}
		}
	}
}

func (t *littleHandTimeWheel) Add(event *event.Event) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	var err error
	e := event.Expiration / t.tick
	event.Round = e / t.wheelSize
	index := (e + t.bucketIndex) % t.wheelSize
	_bucket := t.head
	for i := uint64(0); i < index; i++ {
		_bucket = _bucket.Next()
	}
	bucket := (_bucket.Value).(*littleHandBucket)
	if index == t.bucketIndex {
		event.CurRound += 1
	}
	err = t.lhFile.AddEvent(event, t.tickRound)
	if err != nil {
		log.Errorf("failed to add event %+v to little hand file %+v", event, t.lhFile)
		return err
	}
	err = bucket.Add(event)
	if err != nil {
		log.Error("add event to bucket fail")
		return err
	}
	return err
}

func (t *littleHandTimeWheel) Lookup() ([]*event.Event, bool) {
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
