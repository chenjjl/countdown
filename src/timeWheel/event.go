package timeWheel

import (
	"errors"
	"time"
)

type Event struct {
	bucket     *bucket
	round      uint64
	curRound   uint64
	expiration uint64

	topic string
	tags  string
}

func NewEvent(topic string, tags string, expiration time.Duration) (*Event, error) {
	if topic == "" || tags == "" {
		return nil, errors.New("topic or tags is empty")
	}
	_expiration := uint64(expiration / time.Millisecond)
	if _expiration <= 0 {
		return nil, errors.New("expiration of event must be greater than 0ms")
	}
	return &Event{
		bucket:     nil,
		round:      0,
		curRound:   0,
		expiration: _expiration,

		topic: topic,
		tags:  tags,
	}, nil
}
