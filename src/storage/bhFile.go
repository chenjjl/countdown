package storage

import (
	"countdown/src/event"
	"countdown/src/timeWheel"
	"time"
)

// BhFile event file of big hand time wheel
type BhFile struct {
	*EventFile
}

// CreateBhFile create file for big hand time wheel
func CreateBhFile(expiration time.Duration, round uint64, fileName string, tick uint64) (*BhFile, error) {
	err := CreateDir()
	if err != nil {
		log.Errorf("can not create file for little hand time wheel, because failed to create dir %s", DirName)
		return nil, err
	}
	file, err := newBhFile(expiration, round, fileName, tick)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func newBhFile(expiration time.Duration, round uint64, fileName string, tick uint64) (*BhFile, error) {
	_expiration := uint64(expiration / time.Minute)
	eventFile := &EventFile{
		Element: timeWheel.Element{
			Round:      round,
			CurRound:   0,
			Expiration: _expiration,
		},

		Name: fileName,
		Tick: tick,
	}
	return &BhFile{
		EventFile: eventFile,
	}, nil
}

func (f *BhFile) AddEvent(event *event.Event, tickOffset uint64) error {
	event.TickOffset = tickOffset
	return f.addEvent(event)
}
