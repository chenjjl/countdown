package event

import (
	"countdown/src/logger"
	"errors"
	"strconv"
	"strings"
	"time"
)

var log = logger.GetLogger("event")

type Event struct {
	Round      uint64
	CurRound   uint64
	Expiration uint64

	Topic string
	Tags  string
	Body  []byte

	AddBhUnix      int64  // unix of event been added to big hand time wheel
	TickOffset     uint64 // offset from each tick of big hand time wheel
	TickRound      uint64
	Id             string // unique event Id
	ExpirationUnix uint64 // unix of expiration
}

func NewEvent(topic string, tags string, body []byte, id string, expiration time.Duration) (*Event, error) {
	if topic == "" || tags == "" {
		return nil, errors.New("topic or Tags is empty")
	}
	_expiration := uint64(expiration / time.Millisecond)
	if _expiration <= 0 {
		return nil, errors.New("expiration of event must be greater than 0ms")
	}
	return &Event{
		Round:      0,
		CurRound:   0,
		Expiration: _expiration,

		Topic: topic,
		Tags:  tags,
		Body:  body,

		Id:             id,
		ExpirationUnix: _expiration + uint64(time.Now().UnixMilli()),
	}, nil
}

func (e *Event) ResetExpiration() {
	now := uint64(time.Now().UnixMilli())
	if now >= e.ExpirationUnix {
		// todo
		e.Expiration = 1000
	} else {
		e.Expiration = e.ExpirationUnix - now
	}
}

func (e *Event) Encode() string {
	var builder strings.Builder
	builder.WriteString(e.Id)
	builder.WriteByte('_')
	builder.WriteString(e.Topic)
	builder.WriteByte('_')
	builder.WriteString(e.Tags)
	builder.WriteByte('_')
	builder.Write(e.Body)
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatUint(e.TickRound, 10))
	builder.WriteByte('_')
	builder.WriteString(strconv.FormatUint(e.TickOffset, 10))
	builder.WriteByte('_')
	builder.WriteString(strconv.FormatUint(e.CurRound, 10))
	builder.WriteByte('_')
	builder.WriteString(strconv.FormatInt(e.AddBhUnix, 10))
	builder.WriteByte('_')
	builder.WriteString(strconv.FormatUint(e.Expiration, 10))
	builder.WriteByte('_')
	builder.WriteString(strconv.FormatUint(e.Round, 10))
	builder.WriteByte('_')
	builder.WriteString(strconv.FormatUint(e.ExpirationUnix, 10))

	return builder.String()
}

func Decode(s string) *Event {
	data := strings.Split(strings.Trim(s, ","), "_")
	event := &Event{}
	event.Id = data[0]
	event.Topic = data[1]
	event.Tags = data[2]
	event.Body = []byte(data[3])
	event.TickRound, _ = strconv.ParseUint(data[4], 10, 64)
	event.TickOffset, _ = strconv.ParseUint(data[5], 10, 64)
	event.CurRound, _ = strconv.ParseUint(data[6], 10, 64)
	event.AddBhUnix, _ = strconv.ParseInt(data[7], 10, 64)
	event.Expiration, _ = strconv.ParseUint(data[8], 10, 64)
	event.Round, _ = strconv.ParseUint(data[9], 10, 64)
	event.ExpirationUnix, _ = strconv.ParseUint(data[10], 10, 64)
	return event
}
