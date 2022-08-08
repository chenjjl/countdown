package timeWheel

import "time"

type TimeWheel struct {
	lilHandTimeWheel *littleHandTimeWheel
	bigHandTimeWheel *bigHandTimeWheel
}

func NewTimeWheel(lhTick time.Duration, lhWheelSize uint64, bhTick time.Duration, bhWheelSize uint64) *TimeWheel {
	lilHandTimeWheel := NewLittleHandTimeWheel(lhTick, lhWheelSize)
	return &TimeWheel{
		lilHandTimeWheel: lilHandTimeWheel,
		bigHandTimeWheel: NewBigHandTimeWheel(bhTick, bhWheelSize, lilHandTimeWheel),
	}
}

func (t *TimeWheel) Start() {
	go t.bigHandTimeWheel.Start()
	go t.lilHandTimeWheel.Start()
	time.Sleep(3 * time.Second) // wait to start up
}

func (t *TimeWheel) Add(event *Event) error {
	if event.Expiration/uint64(time.Second.Milliseconds()) >= t.bigHandTimeWheel.tick {
		err := t.bigHandTimeWheel.Add(event)
		if err != nil {
			log.Error("add event to big hand time wheel fail")
			return err
		}
	} else {
		err := t.lilHandTimeWheel.Add(event)
		if err != nil {
			log.Error("add event to little hand time wheel fail")
			return err
		}
	}
	return nil
}
