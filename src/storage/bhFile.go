package storage

import (
	"bufio"
	"countdown/src/event"
	"countdown/src/timeWheel"
	"io"
	"os"
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
		log.Errorf("can not create file for little hand time wheel, because failed to create dir %s", dirName)
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
	return f.addEvent(dirName+f.Name, event)
}

func (f *EventFile) GetEvents(handle func(*event.Event) error) error {
	file, err := os.Open(dirName + f.Name)
	defer file.Close()
	if err != nil {
		log.Errorf("failed to open file %s", dirName+f.Name)
		return err
	}
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString(',')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Errorf("failed to read line from big hand time wheel file")
			return err
		}
		e := event.Decode(line)
		e.Expiration = (e.Expiration%(f.Tick*uint64(time.Second.Milliseconds())) + e.TickOffset) % (f.Tick * uint64(time.Second.Milliseconds()))
		err = handle(e)
		if err != nil {
			return err
		}
	}
}

func (f *EventFile) Remove() error {
	err := os.Remove(dirName + f.Name)
	if err != nil {
		log.Errorf("failed to remove file %s", f.Name)
		return err
	}
	return nil
}
