package timeWheel

import (
	"container/list"
	"countdown/src/event"
)

type littleHandBucket struct {
	events *list.List
}

func newLittleHandBucket() *littleHandBucket {
	return &littleHandBucket{
		events: list.New(),
	}
}

func (b *littleHandBucket) add(event *event.Event) error {
	b.events.PushBack(event)
	log.Infof("add event to little hand time wheel %+v", event)
	return nil
}

func (b *littleHandBucket) lookup() ([]*event.Event, error) {
	var eventsRes []*event.Event
	var n *list.Element
	for e := b.events.Front(); e != nil; e = n {
		event := (e.Value).(*event.Event)
		n = e.Next()
		if event.CurRound >= event.Round {
			b.events.Remove(e)
			eventsRes = append(eventsRes, event)
		} else {
			event.CurRound += 1
		}
	}
	return eventsRes, nil
}
