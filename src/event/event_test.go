package event

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getEvents() (*Event, error) {
	id := uuid.NewV4().String()
	event, err := NewEvent("topic", "tag1", id, time.Duration(10)*time.Second)
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

	_event := Decode("93abe957-2500-4af5-96b6-be56bc7aa83c_topic1_tag7303_0_40484_0_1660814669484_17612000_0_1660832281484")
	assert.Equal(t, event, _event)
	t.Logf("event is %+v", _event)
}
