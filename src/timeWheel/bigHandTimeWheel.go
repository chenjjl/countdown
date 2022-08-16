package timeWheel

import (
	"container/list"
	"countdown/src/event"
	"countdown/src/storage"
	"errors"
	"strconv"
	"sync"
	"time"
)

type bigHandTimeWheel struct {
	wheel

	*littleHandTimeWheel
	tickUnix int64 // unix of each tick
	mu       sync.Mutex
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
	t.tickUnix = time.Now().Unix() * time.Second.Milliseconds()
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
		events, err := file.GetEvents()
		if err != nil {
			log.Error(err)
		}
		err = t.littleHandTimeWheel.Add(events...)
		if err != nil {
			log.Error(err)
		}
	}
	t.tickUnix = time.Now().Unix() * time.Second.Milliseconds()
	t.littleHandTimeWheel.tickRound += 1
	log.Infof("big hand time wheel tick unix is %d, TickRound is %d", t.tickUnix, t.littleHandTimeWheel.tickRound)
}

func (t *bigHandTimeWheel) Add(event *event.Event) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	_expiration := event.Expiration / uint64(time.Second.Milliseconds())
	if _expiration < t.tick {
		return errors.New("fail to add event that expiration smaller than tick to big hand time wheel")
	}
	event.AddBhUnix = time.Now().UnixMilli()
	tickOffset := uint64(event.AddBhUnix - t.tickUnix)
	_tickOffset := tickOffset / uint64(time.Second.Milliseconds())
	e := (_expiration + _tickOffset) / t.tick
	// get current bucket
	index := (e + t.bucketIndex) % t.wheelSize
	_bucket := t.head
	for i := uint64(0); i < index; i++ {
		_bucket = _bucket.Next()
	}
	bucket := (_bucket.Value).(*bigHandBucket)

	fileRound := e / t.wheelSize
	fileName := strconv.FormatUint(uint64(t.tickUnix)+e*t.tick*uint64(time.Second.Milliseconds()), 10)
	// find whether bucket that has file with the same name be exist
	file, err := bucket.LookupFiles(fileName)
	if err != nil {
		log.Error("can not lookup files")
		return err
	}
	if file == nil {
		file, err = storage.CreateBhFile(time.Duration(_expiration)*time.Second, fileRound, fileName, t.tick)
		if err != nil {
			log.Error("can not create a new file")
			return err
		}
		if index == t.bucketIndex {
			file.CurRound += 1
		}
		log.Infof("create file %+v, index = %d, bucketIndex = %d", file, index, t.bucketIndex)
		err := bucket.Add(file)
		if err != nil {
			log.Error("can not add a new file")
			return err
		}
	}
	if err = file.AddEvent(event, tickOffset); err != nil {
		log.Error("can not add event to file")
		return err
	}
	return nil
}

func (t *bigHandTimeWheel) Lookup() (*storage.BhFile, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
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
