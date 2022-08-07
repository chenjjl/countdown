package timeWheel

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var lilTimeWheel = NewLittleHandTimeWheel(time.Second, 8)

func TestLittleHandTimeWheel_Add(t *testing.T) {
	var event, _ = NewEvent("Topic", "tag1", 8*time.Second)
	err := lilTimeWheel.Add(event)
	if err != nil {
		t.Error(err)
	}
}

func TestLittleHandTimeWheel_Lookup(t *testing.T) {
	go lilTimeWheel.Start()

	n := 10
	timeRandLimit := 10
	time.Sleep(4 * time.Second)
	rand.Seed(time.Now().UnixNano())
	start := time.Now().UnixMilli()
	for i := 0; i < n; i++ {
		randTime := rand.Intn(timeRandLimit) + 1
		event, err := NewEvent("topic1", "tag"+strconv.Itoa(i), time.Duration(randTime)*time.Second)
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
		t.Logf("Expiration of event is %d, offset is %d", _event.Expiration, end-start)
		totalOffset += int64(math.Abs(float64(end-start))) - int64(_event.Expiration)
	}
	t.Logf("total message num is %d, average offset is %d", n, totalOffset/int64(n))
}
