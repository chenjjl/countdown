package event

import (
	"countdown/src/timeWheel"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
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

func (e *Event) Encode() string {
	var builder strings.Builder
	builder.WriteString(e.Id)
	builder.WriteByte('-')
	builder.WriteString(e.Topic)
	builder.WriteByte('-')
	builder.WriteString(e.Tags)
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatUint(e.TickRound, 10))
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatUint(e.TickOffset, 10))
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatUint(e.CurRound, 10))
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatInt(e.AddBhUnix, 10))
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatUint(e.Expiration, 10))
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatUint(e.Round, 10))
	builder.WriteByte('-')

	return builder.String()
}

func Decode(s string) *Event {
	data := strings.Split(s, "-")
	var event *Event
	event.Id = data[0]
	event.Topic = data[1]
	event.Tags = data[2]
	event.TickRound, _ = strconv.ParseUint(data[3], 10, 64)
	event.TickOffset, _ = strconv.ParseUint(data[4], 10, 64)
	event.CurRound, _ = strconv.ParseUint(data[5], 10, 64)
	event.AddBhUnix, _ = strconv.ParseInt(data[6], 10, 64)
	event.Expiration, _ = strconv.ParseUint(data[7], 10, 64)
	event.Round, _ = strconv.ParseUint(data[8], 10, 64)
	return event
}
