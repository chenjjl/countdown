package storage

import (
	"countdown/src/event"
	"countdown/src/logger"
	"countdown/src/timeWheel"
	"os"
)

var log = logger.GetLogger("storage")

type EventFile struct {
	timeWheel.Element

	Name string
	Tick uint64
}

func (f *EventFile) addEvent(fileName string, event *event.Event) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)
	if err != nil {
		return err
	}
	_, err = file.WriteString(event.Encode() + ",")
	if err != nil {
		return err
	}
	log.Infof("event %+v be added to file %+v", event, f)
	return nil
}
