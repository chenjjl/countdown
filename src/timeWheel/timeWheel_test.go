package timeWheel

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var timeWheel = NewTimeWheel(time.Second, 8, time.Minute, 8)

func TestTimeWheel_Lookup(t *testing.T) {
	timeWheel.Start()

	startUnixMap := make(map[string]int64)
	eventMap := make(map[string]uint64)
	rand.Seed(time.Now().UnixNano())
	timeRandLimit := 5 * 60 * 60
	n := 10000
	for i := 0; i < n; i++ {
		i := i
		go func() {
			time.Sleep(time.Duration(rand.Intn(5*60*60)) * time.Second)
			randTime := rand.Intn(timeRandLimit) + 1
			event, err := NewEvent("topic1", "tag"+strconv.Itoa(i), time.Duration(randTime)*time.Second)
			if err != nil {
				t.Error(err)
			}
			err = timeWheel.Add(event)
			if err != nil {
				t.Error(err)
			}
			mu.Lock()
			startUnixMap[event.Topic+"-"+event.Tags] = time.Now().UnixMilli()
			eventMap[event.Topic+"-"+event.Tags] = event.Expiration
			mu.Unlock()
		}()
	}

	i := 0
	totalOffset := int64(0)
	for i < n {
		_event := <-timeWheel.lilHandTimeWheel.c
		i += 1
		end := time.Now().UnixMilli()
		expectExp := eventMap[_event.Topic+"-"+_event.Tags]
		startUnix := startUnixMap[_event.Topic+"-"+_event.Tags]
		actualExp := float64(end - startUnix)
		offset := actualExp - float64(expectExp)
		log.Infof("event %+v, expected expiration is %d, actual expiration is %f, abs offset >= 500ms is %t, offset is %f", _event, expectExp, actualExp, math.Abs(offset) >= 500, offset)
		totalOffset += int64(offset)
	}
	t.Logf("total message num is %d, average offset is %d", n, totalOffset/int64(n))
}
