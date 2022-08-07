package timeWheel

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var bigTimeWheel = NewBigHandTimeWheel(time.Minute, 8, lilTimeWheel)

func TestBigHandTimeWheel_Add(t *testing.T) {
	var event, _ = NewEvent("Topic", "tag1", 1*time.Minute)
	err := bigTimeWheel.Add(event)
	if err != nil {
		t.Error(err)
	}
}

func TestBigHandTimeWheel_Lookup(t *testing.T) {
	go bigTimeWheel.Start()
	go lilTimeWheel.Start()

	eventMap := make(map[string]uint64)

	n := 10
	timeRandLimit := 60
	time.Sleep(5 * time.Second)

	rand.Seed(time.Now().UnixNano())
	start := time.Now().UnixMilli()
	for i := 0; i < n; i++ {
		randTime := rand.Intn(timeRandLimit) + 60
		event, err := NewEvent("topic1", "tag"+strconv.Itoa(i), time.Duration(randTime)*time.Second)
		if err != nil {
			t.Error(err)
		}
		eventMap[event.Tags] = event.Expiration
		err = bigTimeWheel.Add(event)
		if err != nil {
			t.Error(err)
		}
	}
	i := 0
	totalOffset := int64(0)
	for i < n {
		_event := <-bigTimeWheel.c
		i += 1
		end := time.Now().UnixMilli()
		expectExp := eventMap[_event.Tags]
		log.Infof("actual expiration is %d, offset is %d", expectExp, end-start)
		log.Infof("start is %d, end is %d", start, end)
		totalOffset += int64(math.Abs(float64(end-start))) - int64(expectExp)
	}
	t.Logf("total message num is %d, average offset is %d", n, totalOffset/int64(n))
}
