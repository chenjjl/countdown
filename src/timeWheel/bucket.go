package timeWheel

import (
	"container/list"
	"errors"
)

type bucket struct {
	events *list.List
}

func NewBucket() *bucket {
	return &bucket{
		events: list.New(),
	}
}

func (b *bucket) Add(event *Event) error {
	if event.bucket != nil {
		return errors.New("event has added")
	}
	b.events.PushBack(event)
	event.bucket = b
	return nil
}

func (b *bucket) Lookup() (*Event, error) {
	for e := b.events.Front(); e != nil; e = e.Next() {
		event := (e.Value).(*Event)
		if event.curRound == event.round {
			return event, nil
		}
		event.curRound += 1
	}
	return nil, nil
}
