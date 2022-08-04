package timeWheel

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var timeWheel = NewTimeWheel(time.Second, 8)

func TestAdd(t *testing.T) {
	var event, _ = NewEvent("topic", "tag1", 10*time.Second)
	err := timeWheel.Add(event)
	if err != nil {
		t.Error(err)
	}
}

func TestLookUp(t *testing.T) {
	go timeWheel.Start()

	n := 5
	timeRandLimit := 10
	startTimeMap := make(map[string]int)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < n; i++ {
		randTime := rand.Intn(timeRandLimit) + 1
		_randTime := time.Duration(randTime) * time.Second
		name := "topic1" + "-" + "tag" + string(rune(randTime))
		event, err := NewEvent("topic1", "tag"+strconv.FormatUint(uint64(_randTime.Milliseconds()), 10), time.Duration(randTime)*time.Second)
		if err != nil {
			t.Error(err)
		}
		err = timeWheel.Add(event)
		if err != nil {
			t.Error(err)
		}
		start := time.Now().Second()
		startTimeMap[name] = start
	}

	for {
		_event := <-timeWheel.c

		end := time.Now().Second()
		tag := "tag" + strconv.FormatUint(_event.expiration, 10)
		name := "topic1" + "-" + tag
		start := startTimeMap[name]
		assert.Equal(t, "topic1", _event.topic)
		assert.Equal(t, tag, _event.tags)
		assert.Equal(t, _event.expiration, time.Duration(end-start)*time.Millisecond)
	}
}
