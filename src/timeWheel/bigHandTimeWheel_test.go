package timeWheel

import (
	event2 "countdown/src/event"
	uuid "github.com/satori/go.uuid"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

var bigTimeWheel = NewBigHandTimeWheel(time.Minute, 8, lilTimeWheel)
var mu = sync.Mutex{}

func TestBigHandTimeWheel_Add(t *testing.T) {
	var event, _ = event2.NewEvent("Topic", "tag1", nil, uuid.NewV4().String(), 1*time.Minute)
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
	timeRandLimit := 500
	time.Sleep(3 * time.Second)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		i := i
		go func() {
			time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
			randTime := rand.Intn(timeRandLimit) + 60
			event, err := event2.NewEvent("topic1", "tag"+strconv.Itoa(i), nil, uuid.NewV4().String(), time.Duration(randTime)*time.Second)
			if err != nil {
				t.Error(err)
			}
			mu.Lock()
			eventMap[event.Tags] = event.Expiration
			mu.Unlock()
			err = bigTimeWheel.Add(event)
			if err != nil {
				t.Error(err)
			}
		}()
	}

	i := 0
	totalOffset := int64(0)
	for i < n {
		_event := <-lilTimeWheel.c
		i += 1
		end := time.Now().UnixMilli()
		expectExp := eventMap[_event.Tags]
		log.Infof("expected expiration is %d, actual expiration is %d", expectExp, end-_event.AddBhUnix)
		totalOffset += int64(math.Abs(float64(end-_event.AddBhUnix))) - int64(expectExp)
	}
	t.Logf("total message num is %d, average offset is %d", n, totalOffset/int64(n))
}
