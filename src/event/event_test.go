package event

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getEvents() (*Event, error) {
	id := uuid.NewV4().String()
	event, err := NewEvent("topic", []byte("test"), id, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func TestEvent_Decode(t *testing.T) {

}

func TestEvent_Encode(t *testing.T) {
	event, err := getEvents()
	if err != nil {
		t.Error("fail to create a new event")
		t.Error(err)
	}
	s := event.Encode()
	t.Logf("encode event is %s", s)

	_event := Decode(s)
	assert.Equal(t, event, _event)
	t.Logf("event is %+v", _event)
}
