package storage

import (
	"bufio"
	"countdown/src/event"
	"countdown/src/logger"
	"io"
	"os"
)

var log = logger.GetLogger("storage")

type EventFile struct {
	Round      uint64
	CurRound   uint64
	Expiration uint64

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

func getEvents(fileName string, eventMap map[string]*event.Event) error {
	if !Exists(fileName) {
		return nil
	}
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Errorf("failed to open file %s", fileName)
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
		eventMap[e.Id] = e
	}
}
