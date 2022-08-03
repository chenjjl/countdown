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
