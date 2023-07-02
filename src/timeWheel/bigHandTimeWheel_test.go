package timeWheel

import (
	"sync"
	"time"
)

var bigTimeWheel = newBigHandTimeWheel(time.Minute, 8, lilTimeWheel)
var mu = sync.Mutex{}

//func TestBigHandTimeWheel_Lookup(t *testing.T) {
//	go bigTimeWheel.start()
//	go lilTimeWheel.start()
//
//	eventMap := make(map[string]uint64)
//
//	n := 10
//	timeRandLimit := 500
//	time.Sleep(3 * time.Second)
//
//	rand.Seed(time.Now().UnixNano())
//	for i := 0; i < n; i++ {
//		i := i
//		go func() {
//			time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
//			randTime := rand.Intn(timeRandLimit) + 60
//			event, err := event2.NewEvent("topic1", nil, uuid.NewV4().String(), time.Duration(randTime)*time.Second)
//			if err != nil {
//				t.Error(err)
//			}
//			err = bigTimeWheel.add(event)
//			if err != nil {
//				t.Error(err)
//			}
//		}()
//	}
//
//	i := 0
//	totalOffset := int64(0)
//	for i < n {
//		_event := <-lilTimeWheel.c
//		i += 1
//		end := time.Now().UnixMilli()
//		log.Infof("expected expiration is %d, actual expiration is %d", expectExp, end-_event.AddBhUnix)
//		totalOffset += int64(math.Abs(float64(end-_event.AddBhUnix))) - int64(expectExp)
//	}
//	t.Logf("total message num is %d, average offset is %d", n, totalOffset/int64(n))
//}
