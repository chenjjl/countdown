package timeWheel

import (
	"container/list"
	"countdown/src/event"
)

type littleHandBucket struct {
	events *list.List
}

func NewLittleHandBucket() *littleHandBucket {
	return &littleHandBucket{
		events: list.New(),
	}
}

func (b *littleHandBucket) Add(event *event.Event) error {
	b.events.PushBack(event)
	log.Infof("add event to little hand time wheel %+v", event)
	return nil
}

func (b *littleHandBucket) Lookup() ([]*event.Event, error) {
	var eventsRes []*event.Event
	var n *list.Element
	for e := b.events.Front(); e != nil; e = n {
		event := (e.Value).(*event.Event)
		n = e.Next()
		if event.CurRound == event.Round {
			b.events.Remove(e)
			eventsRes = append(eventsRes, event)
		} else {
			event.CurRound += 1
		}
	}
	return eventsRes, nil
}
