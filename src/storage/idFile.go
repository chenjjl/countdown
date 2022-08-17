package storage

import (
	"countdown/src/event"
	"encoding/json"
	"os"
)

const eventIdFileName = "event_id.txt"

type IdFile struct {
	Name string
}

type item struct {
	Id        string
	TickRound uint64
}

// CreateEventIdFile create file for recording event Id in the current little hand time wheel
func CreateEventIdFile() (*IdFile, error) {
	err := CreateDir()
	if err != nil {
		log.Errorf("can not create file for little hand time wheel, because failed to create dir %s", dirName)
	}
	file, err := newEventIdFile()
	if err != nil {
		return nil, err
	}
	return file, nil
}

func newEventIdFile() (*IdFile, error) {
	return &IdFile{
		Name: eventIdFileName,
	}, nil
}

func (f *IdFile) Add(event *event.Event) error {
	file, err := os.OpenFile(dirName+f.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)
	if err != nil {
		return err
	}
	item := &item{Id: event.Id, TickRound: event.TickRound}
	data, _ := json.Marshal(item)
	_, err = file.WriteString(string(data) + ",")
	if err != nil {
		return err
	}
	log.Infof("Id %+v be added to idFile %+v", item, f)
	return nil
}
