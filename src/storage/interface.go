package storage

import (
	"countdown/src/event"
	"countdown/src/logger"
	"countdown/src/timeWheel"
	"encoding/json"
	"os"
	"time"
)

const lhFileName = "lh.txt"

var log = logger.GetLogger("storage")

type EventFile struct {
	timeWheel.Element

	Name string
	Tick uint64
}

func (f *EventFile) addEvent(event *event.Event) error {
	file, err := os.OpenFile(DirName+f.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)
	if err != nil {
		return err
	}
	_, err = file.WriteString(encodeEvent(event) + ",")
	if err != nil {
		return err
	}
	log.Infof("event %+v be added to file %+v", event, f)
	return nil
}

func (f *EventFile) GetEvents() ([]*event.Event, error) {
	data, err := os.ReadFile(DirName + f.Name)
	err = os.Remove(DirName + f.Name)
	if err != nil {
		log.Errorf("failed to remove file %s", f.Name)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	var eventArr []*event.Event
	err = json.Unmarshal(WrapperFileContent(data), &eventArr)
	if err != nil {
		return nil, err
	}
	for _, event := range eventArr {
		event.Expiration = (event.Expiration%(f.Tick*uint64(time.Second.Milliseconds())) + event.TickOffset) % (f.Tick * uint64(time.Second.Milliseconds()))
	}
	return eventArr, nil
}
