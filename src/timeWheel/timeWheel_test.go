package timeWheel

import (
	"testing"
	"time"
)

func TestNewTimeWheel(t *testing.T) {
	event, err := NewEvent("topic", "tag1", 10*time.Second)
	if err != nil {
		t.Error(err)
	}

	timeWheel := NewTimeWheel(time.Second, 8)
	err = timeWheel.Add(event)
	if err != nil {
		t.Error(err)
	}
}
