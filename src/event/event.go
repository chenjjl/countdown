package event

import (
	"countdown/src/timeWheel"
	"encoding/json"
	"errors"
	"time"
)

type Event struct {
	timeWheel.Element

	Topic string
	Tags  string

	AddBhUnix  int64  // unix of event been added to big hand time wheel
	TickOffset uint64 // offset from each tick of big hand time wheel
	TickRound  uint64
	Id         string // unique event Id
}

func NewEvent(topic string, tags string, id string, expiration time.Duration) (*Event, error) {
	if topic == "" || tags == "" {
		return nil, errors.New("topic or Tags is empty")
	}
	_expiration := uint64(expiration / time.Millisecond)
	if _expiration <= 0 {
		return nil, errors.New("expiration of event must be greater than 0ms")
	}
	return &Event{
		Element: timeWheel.Element{
			Round:      0,
			CurRound:   0,
			Expiration: _expiration,
		},

		Topic: topic,
		Tags:  tags,
		Id:    id,
	}, nil
}

func (e *Event) ToString() string {
	data, _ := json.Marshal(e)
	return string(data)
}
