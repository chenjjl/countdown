package storage

import (
	"countdown/src/event"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const lhFilePrefixName = "lh-"

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
	eventFile := &EventFile{Name: lhFilePrefixName}
	return &LhFile{
		EventFile: eventFile,
	}, nil
}

func (f *LhFile) AddEvent(event *event.Event, tickRound uint64) error {
	event.TickRound = tickRound
	return f.addEvent(DirName+f.Name+strconv.FormatUint(tickRound, 10), event)
}

func ReloadLhEvents() ([]*event.Event, error) {
	files, err := ioutil.ReadDir(DirName)
	if err != nil {
		log.Errorf("failed to read dir %s", DirName)
		return nil, err
	}
	maxRound := int64(0)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), lhFilePrefixName) {
			round, _ := strconv.ParseInt(strings.TrimPrefix(file.Name(), lhFilePrefixName), 10, 64)
			if round > maxRound {
				maxRound = round
			}
		}
	}

	eventMap := make(map[string]*event.Event)
	idItemMap := make(map[string]*idItem)
	eventFileName := DirName + lhFilePrefixName + strconv.FormatInt(maxRound, 10)
	idFileName := DirName + eventIdFilePrefixName + strconv.FormatInt(maxRound, 10)
	err = getEvents(eventFileName, eventMap)
	if err != nil {
		return nil, err
	}
	err = getIdItems(idFileName, idItemMap)
	if err != nil {
		return nil, err
	}
	eventFileName = DirName + lhFilePrefixName + strconv.FormatInt(maxRound-1, 10)
	idFileName = DirName + eventIdFilePrefixName + strconv.FormatInt(maxRound-1, 10)
	if Exists(eventFileName) {
		err = getEvents(eventFileName, eventMap)
		if err != nil {
			return nil, err
		}
		err = getIdItems(idFileName, idItemMap)
		if err != nil {
			return nil, err
		}
	}
	log.Infof("eventMap is %+v", eventMap)
	log.Infof("idItemMap is %+v", idItemMap)
	var needReload []*event.Event
	for _, e := range eventMap {
		if _, exist := idItemMap[e.Id]; !exist {
			e.ResetExpiration()
			needReload = append(needReload, e)
		}
	}
	return needReload, nil
}

func RemoveLhEventFiles() {
	files, err := ioutil.ReadDir(DirName)
	if err != nil {
		log.Errorf("failed to read dir %s", DirName)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), lhFilePrefixName) {
			err = os.Remove(DirName + file.Name())
			if err != nil {
				log.Errorf("failed to remove file %s", file.Name())
			}
		}
		if strings.HasPrefix(file.Name(), eventIdFilePrefixName) {
			err = os.Remove(DirName + file.Name())
			if err != nil {
				log.Errorf("failed to remove file %s", file.Name())
			}
		}
	}
}
