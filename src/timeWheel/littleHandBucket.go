package timeWheel

import (
	"container/list"
)

type littleHandBucket struct {
	events *list.List
}

func NewLittleHandBucket() *littleHandBucket {
	return &littleHandBucket{
		events: list.New(),
	}
}

func (b *littleHandBucket) Add(event *Event) error {
	b.events.PushBack(event)
	log.Infof("add event to little hand time wheel %+v", event)
	return nil
}

func (b *littleHandBucket) Lookup() ([]*Event, error) {
	var eventsRes []*Event
	var n *list.Element
	for e := b.events.Front(); e != nil; e = n {
		event := (e.Value).(*Event)
		n = e.Next()
		if event.curRound == event.round {
			b.events.Remove(e)
			eventsRes = append(eventsRes, event)
		} else {
			event.curRound += 1
		}
	}
	return eventsRes, nil
}
