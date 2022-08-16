package storage

import (
	"countdown/src/event"
	"os"
)

// LhFile event file of little hand time wheel
type LhFile struct {
	*EventFile
}

// CreateLhFile create file for little hand time wheel
func CreateLhFile() (*LhFile, error) {
	err := CreateDir()
	if err != nil {
		log.Errorf("can not create file for little hand time wheel, because failed to create dir %s", DirName)
	}
	file, err := newLhFile()
	if err != nil {
		return nil, err
	}
	return file, nil
}

func newLhFile() (*LhFile, error) {
	eventFile := &EventFile{Name: lhFileName}
	return &LhFile{
		EventFile: eventFile,
	}, nil
}

func (f *LhFile) AddEvent(event *event.Event, tickRound uint64) error {
	event.TickRound = tickRound
	// todo compress event file
	return f.addEvent(event)
}

func (f *LhFile) GetEventsBehind() ([]*event.Event, error) {
	file, err := os.OpenFile(DirName+f.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)
	if err != nil {
		return nil, err
	}

}
