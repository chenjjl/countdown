package timeWheel

import (
	event2 "countdown/src/event"
	"github.com/satori/go.uuid"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var timeWheel = NewTimeWheel(time.Second, 8, time.Minute, 8)

func TestTimeWheel_Lookup(t *testing.T) {
	timeWheel.Start()

	rand.Seed(time.Now().UnixNano())
	timeRandLimit := 60 * 60
	n := 10000
	for i := 0; i < n; i++ {
		i := i
		go func() {
			time.Sleep(time.Duration(rand.Intn(timeRandLimit)) * time.Second)
			randTime := rand.Intn(timeRandLimit)
			id := uuid.NewV4().String()
			event, err := event2.NewEvent("topic1", "tag"+strconv.Itoa(i), nil, id, time.Duration(randTime)*time.Second)
			if err != nil {
				t.Error(err)
			}
			err = timeWheel.Add(event)
			if err != nil {
				t.Error(err)
			}
		}()
	}

	i := 0
	totalOffset := int64(0)
	for {
		_event := <-timeWheel.lilHandTimeWheel.c
		i += 1
		now := time.Now().UnixMilli()
		offset := now - int64(_event.ExpirationUnix)
		log.Infof("event %+v, abs offset >= 500ms is %t, offset is %d", _event, math.Abs(float64(offset)) >= 500, offset)
		totalOffset += offset
	}
	t.Logf("total message num is %d, average offset is %d", n, totalOffset/int64(n))
}
