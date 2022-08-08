package timeWheel

import (
	"encoding/json"
	"errors"
	"time"
)

type Event struct {
	element

	Topic string
	Tags  string

	AddBhUnix  int64  // unix of event been added to big hand time wheel
	TickOffset uint64 // offset from each tick of big hand time wheel
}

func NewEvent(topic string, tags string, expiration time.Duration) (*Event, error) {
	if topic == "" || tags == "" {
		return nil, errors.New("Topic or Tags is empty")
	}
	_expiration := uint64(expiration / time.Millisecond)
	if _expiration <= 0 {
		return nil, errors.New("Expiration of event must be greater than 0ms")
	}
	return &Event{
		element: element{
			round:      0,
			curRound:   0,
			Expiration: _expiration,
		},

		Topic: topic,
		Tags:  tags,
	}, nil
}

func (e *Event) toString() string {
	data, _ := json.Marshal(e)
	return string(data)
}
