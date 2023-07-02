package storage

import (
	"bufio"
	"countdown/src/event"
	"countdown/src/utils"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const BhFileNamePrefix = "bh-"
const BhFileNameOldPrefix = "o-"

// BhFile event file of big hand time wheel
type BhFile struct {
	*EventFile
}

// CreateBhFile create file for big hand time wheel
func CreateBhFile(expiration time.Duration, round uint64, fileName string, tick uint64) (*BhFile, error) {
	err := utils.CreateDir(utils.EventLogDir)
	if err != nil {
		log.Errorf("can not create file for little hand time wheel, because failed to create dir %s", utils.EventLogDir)
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
		Round:      round,
		CurRound:   0,
		Expiration: _expiration,

		Name: BhFileNamePrefix + fileName,
		Tick: tick,
	}
	return &BhFile{
		EventFile: eventFile,
	}, nil
}

func (f *BhFile) AddEvent(event *event.Event, tickOffset uint64) error {
	event.TickOffset = tickOffset
	return f.addEvent(utils.EventLogDir+f.Name, event)
}

func (f *EventFile) GetEvents(handle func(*event.Event) error) error {
	file, err := os.Open(utils.EventLogDir + f.Name)
	defer file.Close()
	if err != nil {
		log.Errorf("failed to open file %s", utils.EventLogDir+f.Name)
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
	err := os.Remove(utils.EventLogDir + f.Name)
	if err != nil {
		log.Errorf("failed to remove file %s", f.Name)
		return err
	}
	return nil
}

func ReloadBhEvents() ([]*event.Event, error) {
	files, err := ioutil.ReadDir(utils.EventLogDir)
	if err != nil {
		log.Errorf("failed to read dir %s", utils.EventLogDir)
		return nil, err
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), BhFileNamePrefix) {
			err = os.Rename(utils.EventLogDir+file.Name(), utils.EventLogDir+BhFileNameOldPrefix+file.Name())
			if err != nil {
				log.Errorf("failed to rename file name %s to new file name %s", file.Name(), BhFileNameOldPrefix+file.Name())
			}
		}
	}
	files, err = ioutil.ReadDir(utils.EventLogDir)
	if err != nil {
		log.Errorf("failed to read dir %s", utils.EventLogDir)
		return nil, err
	}
	var needReload []*event.Event
	for _, file := range files {
		eventMap := make(map[string]*event.Event)
		if strings.HasPrefix(file.Name(), BhFileNameOldPrefix+BhFileNamePrefix) {
			err = getEvents(utils.EventLogDir+file.Name(), eventMap)
			if err != nil {
				log.Errorf("break loading file %s", file.Name())
				log.Error(err)
				break
			}
			for _, e := range eventMap {
				e.ResetExpiration()
				needReload = append(needReload, e)
			}
		}
	}
	return needReload, nil
}

func RemoveBhEventFile() {
	files, err := ioutil.ReadDir(utils.EventLogDir)
	if err != nil {
		log.Errorf("failed to read dir %s", utils.EventLogDir)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), BhFileNameOldPrefix+BhFileNamePrefix) {
			err = os.Remove(utils.EventLogDir + file.Name())
			if err != nil {
				log.Errorf("failed to remove file %s", file.Name())
			}
		}
	}
}
