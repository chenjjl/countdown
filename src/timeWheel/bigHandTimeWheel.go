package timeWheel

import (
	"container/list"
	"errors"
	"time"
)

type bigHandTimeWheel struct {
	wheel

	*littleHandTimeWheel
	tickUnix int64 // unix of each tick
}

func NewBigHandTimeWheel(tick time.Duration, wheelSize uint64, lilHandTimeWheel *littleHandTimeWheel) *bigHandTimeWheel {
	_tick := uint64(tick / time.Second)
	if _tick < 60 {
		panic(errors.New("big hand's tick must be equal or greater than 1 minute"))
	}
	if wheelSize <= 0 {
		panic(errors.New("big hand's size of timeWheel must be greater than 0"))
	}
	buckets := list.New()
	for i := uint64(0); i < wheelSize; i++ {
		bucket := NewBigHandBucket()
		buckets.PushBack(bucket)
	}
	return &bigHandTimeWheel{
		wheel: wheel{
			tick:        _tick,
			wheelSize:   wheelSize,
			buckets:     buckets,
			bucketIndex: 0,
			curBucket:   buckets.Front(),
			head:        buckets.Front(),
		},
		littleHandTimeWheel: lilHandTimeWheel,
	}
}

func (t *bigHandTimeWheel) Start() {
	ticker := time.NewTicker(time.Duration(t.tick) * time.Second)
	defer ticker.Stop()
	t.tickUnix = time.Now().UnixMilli()
	log.Infof("big hand's ticker has started, tick is %d sec", t.tick)
	for {
		select {
		case <-ticker.C:
			t.doLookup()
		}
	}
}

func (t *bigHandTimeWheel) doLookup() {
	file, ok := t.Lookup()
	if ok {
		events, err := file.getEvents()
		if err != nil {
			log.Error(err)
		}
		err = t.littleHandTimeWheel.Add(events...)
		if err != nil {
			log.Error(err)
		}
	}
	t.tickUnix = time.Now().UnixMilli()
}

func (t *bigHandTimeWheel) Add(event *Event) error {
	_expiration := event.Expiration / uint64(time.Second.Milliseconds())
	if _expiration < t.tick {
		return nil
	}
	e := _expiration / t.tick
	index := (e + t.bucketIndex) % t.wheelSize
	_bucket := t.head
	for i := uint64(0); i < index; i++ {
		_bucket = _bucket.Next()
	}
	bucket := (_bucket.Value).(*bigHandBucket)
	fileRound := e / t.wheelSize
	file, err := bucket.LookupFiles(fileRound)
	if err != nil {
		log.Error("can not lookup files")
		return err
	}
	if file == nil {
		file, err = createFile(time.Duration(_expiration)*time.Second, fileRound, t.tickUnix, t.tick)
		if err != nil {
			log.Error("can not create a new file")
			return err
		}
		if index == t.bucketIndex {
			file.curRound += 1
		}
		log.Infof("create file %+v, index = %d, bucketIndex = %d, event = %+v", file, index, t.bucketIndex, event)
		err := bucket.Add(file)
		if err != nil {
			log.Error("can not add a new file")
			return err
		}
	}
	if err = file.addEvent(event, t.tickUnix); err != nil {
		log.Error("can not add event to file")
		return err
	}
	return nil
}

func (t *bigHandTimeWheel) Lookup() (*File, bool) {
	t.bucketIndex = (t.bucketIndex + 1) % t.wheelSize
	t.curBucket = t.curBucket.Next()

	// circle queue
	if t.curBucket == nil {
		t.curBucket = t.head
	}

	bucket := (t.curBucket.Value).(*bigHandBucket)
	file, err := bucket.Lookup()
	if err != nil {
		log.Error(err)
		return nil, false
	}
	if file == nil {
		return nil, false
	}
	log.Infof("big hand time wheel lookup bucketIndex = %d, file = %+v", t.bucketIndex, file)
	return file, true
}
