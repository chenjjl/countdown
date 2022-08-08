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
	timeRandLimit := 700
	n := 10
	for i := 0; i < n; i++ {
		randTime := rand.Intn(timeRandLimit) + 1
		event, err := NewEvent("topic1", "tag"+strconv.Itoa(i), time.Duration(randTime)*time.Second)
		if err != nil {
			t.Error(err)
		}
		err = timeWheel.Add(event)
		if err != nil {
			t.Error(err)
		}
		startUnixMap[event.Topic+"-"+event.Tags] = time.Now().UnixMilli()
		eventMap[event.Topic+"-"+event.Tags] = event.Expiration
	}

	i := 0
	totalOffset := int64(0)
	for i < n {
		_event := <-timeWheel.lilHandTimeWheel.c
		i += 1
		end := time.Now().UnixMilli()
		expectExp := eventMap[_event.Topic+"-"+_event.Tags]
		startUnix := startUnixMap[_event.Topic+"-"+_event.Tags]
		log.Infof("expected expiration is %d, actual expiration is %d", expectExp, end-startUnix)
		totalOffset += int64(math.Abs(float64(end-startUnix))) - int64(expectExp)
	}
	t.Logf("total message num is %d, average offset is %d", n, totalOffset/int64(n))
}
