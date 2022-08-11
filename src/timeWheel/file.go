package timeWheel

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

const dirName = "/tmp/countdown" + dirPrefix
const dirPrefix string = "/eventLog/"

type File struct {
	element

	name string
	tick uint64
}

func createFile(expiration time.Duration, round uint64, fileName string, tick uint64) (*File, error) {
	file, err := newFile(expiration, round, fileName, tick)
	if err != nil {
		return nil, err
	}
	exist, err := hasDir(dirName)
	if err != nil {
		return nil, err
	}
	if exist {
		return file, nil
	}
	err = os.Mkdir(dirName, os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func hasDir(path string) (bool, error) {
	_, _err := os.Stat(path)
	if _err == nil {
		return true, nil
	}
	if os.IsNotExist(_err) {
		return false, nil
	}
	return false, _err
}

func newFile(expiration time.Duration, round uint64, fileName string, tick uint64) (*File, error) {
	_expiration := uint64(expiration / time.Minute)
	if _expiration < 1 {
		return nil, errors.New("Expiration of file must be equal or greater than 1 minute")
	}
	return &File{
		element: element{
			round:      round,
			curRound:   0,
			Expiration: _expiration,
		},

		name: fileName,
		tick: tick,
	}, nil
}

func (f *File) addEvent(event *Event, tickOffset uint64) error {
	file, err := os.OpenFile(dirName+f.name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)
	if err != nil {
		return err
	}
	event.TickOffset = tickOffset
	_, err = file.WriteString(event.toString() + ",")
	if err != nil {
		return err
	}
	log.Infof("event %+v be added to file %+v", event, f)
	return nil
}

func (f *File) getEvents() ([]*Event, error) {
	data, err := os.ReadFile(dirName + f.name)
	err = os.Remove(dirName + f.name)
	if err != nil {
		log.Errorf("failed to remove file %s", f.name)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	var eventArr []*Event
	err = json.Unmarshal(wrapperFileContent(data), &eventArr)
	if err != nil {
		return nil, err
	}
	for _, event := range eventArr {
		event.Expiration = (event.Expiration%(f.tick*uint64(time.Second.Milliseconds())) + event.TickOffset) % (f.tick * uint64(time.Second.Milliseconds()))
	}
	return eventArr, nil
}

func wrapperFileContent(data []byte) []byte {
	data = append([]byte("["), data...)
	data[len(data)-1] = ']'
	return data
}
