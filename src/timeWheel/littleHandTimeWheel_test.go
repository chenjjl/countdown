package timeWheel

import (
	event2 "countdown/src/event"
	uuid "github.com/satori/go.uuid"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var lilTimeWheel = NewLittleHandTimeWheel(time.Second, 8)

func TestLittleHandTimeWheel_Add(t *testing.T) {
	var event, _ = event2.NewEvent("Topic", "tag1", nil, uuid.NewV4().String(), 8*time.Second)
	err := lilTimeWheel.Add(event)
	if err != nil {
		t.Error(err)
	}
}

func TestLittleHandTimeWheel_Lookup(t *testing.T) {
	go lilTimeWheel.Start()

	n := 100
	timeRandLimit := 30
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
	start := time.Now().UnixMilli()
	for i := 0; i < n; i++ {
		randTime := rand.Intn(timeRandLimit) + 1
		event, err := event2.NewEvent("topic1", "tag"+strconv.Itoa(i), nil, uuid.NewV4().String(), time.Duration(randTime)*time.Second)
		if err != nil {
			t.Error(err)
		}
		err = lilTimeWheel.Add(event)
		if err != nil {
			t.Error(err)
		}
	}
	i := 0
	totalOffset := int64(0)
	for i < n {
		_event := <-lilTimeWheel.c
		i += 1
		end := time.Now().UnixMilli()
		t.Logf("expected expiration is %d, actual expiration is %d", _event.Expiration, end-start)
		totalOffset += int64(math.Abs(float64(end-start))) - int64(_event.Expiration)
	}
	t.Logf("total message num is %d, average offset is %d", n, totalOffset/int64(n))
}
