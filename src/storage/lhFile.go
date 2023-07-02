package storage

import (
	"countdown/src/event"
	"countdown/src/utils"
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
	err := utils.CreateDir(utils.EventLogDir)
	if err != nil {
		log.Errorf("can not create file for little hand time wheel, because failed to create dir %s", utils.EventLogDir)
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
	return f.addEvent(utils.EventLogDir+f.Name+strconv.FormatUint(tickRound, 10), event)
}

func ReloadLhEvents() ([]*event.Event, error) {
	files, err := ioutil.ReadDir(utils.EventLogDir)
	if err != nil {
		log.Errorf("failed to read dir %s", utils.EventLogDir)
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

	eventFileName := utils.EventLogDir + lhFilePrefixName + strconv.FormatInt(maxRound, 10)
	idFileName := utils.EventLogDir + eventIdFilePrefixName + strconv.FormatInt(maxRound, 10)
	err = getEvents(eventFileName, eventMap)
	if err != nil {
		return nil, err
	}
	err = getIdItems(idFileName, idItemMap)
	if err != nil {
		return nil, err
	}

	eventFileName = utils.EventLogDir + lhFilePrefixName + strconv.FormatInt(maxRound-1, 10)
	idFileName = utils.EventLogDir + eventIdFilePrefixName + strconv.FormatInt(maxRound-1, 10)
	err = getEvents(eventFileName, eventMap)
	if err != nil {
		return nil, err
	}
	err = getIdItems(idFileName, idItemMap)
	if err != nil {
		return nil, err
	}

	var needReload []*event.Event
	for _, e := range eventMap {
		if _, exist := idItemMap[e.Id]; !exist {
			e.ResetExpiration()
			needReload = append(needReload, e)
		}
	}
	return needReload, nil
}

func RemoveAllLhEventFiles() {
	files, err := ioutil.ReadDir(utils.EventLogDir)
	if err != nil {
		log.Errorf("failed to read dir %s", utils.EventLogDir)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), lhFilePrefixName) {
			err = os.Remove(utils.EventLogDir + file.Name())
			if err != nil {
				log.Errorf("failed to remove file %s", file.Name())
			}
		}
		if strings.HasPrefix(file.Name(), eventIdFilePrefixName) {
			err = os.Remove(utils.EventLogDir + file.Name())
			if err != nil {
				log.Errorf("failed to remove file %s", file.Name())
			}
		}
	}
}

func RemoveUnusedLhEventFiles(curTickRound uint64) {
	if curTickRound <= 1 {
		return
	}
	files, err := ioutil.ReadDir(utils.EventLogDir)
	if err != nil {
		log.Errorf("failed to read dir %s", utils.EventLogDir)
	}
	for _, file := range files {
		if isUnusedLhEventFiles(file.Name(), curTickRound) || isUnusedEventIdFiles(file.Name(), curTickRound) {
			err = os.Remove(utils.EventLogDir + file.Name())
			if err != nil {
				log.Errorf("failed to remove file %s", file.Name())
			}
		}
	}
}

func isUnusedLhEventFiles(fileName string, curTickRound uint64) bool {
	if !strings.HasPrefix(fileName, lhFilePrefixName) {
		return false
	}
	prevTickRound := curTickRound - 1
	_curTickRound := strconv.FormatUint(curTickRound, 10)
	_prevTickRound := strconv.FormatUint(prevTickRound, 10)

	lhFileRound := strings.TrimPrefix(fileName, lhFilePrefixName)
	if lhFileRound == _curTickRound || lhFileRound == _prevTickRound {
		return false
	} else {
		return true
	}
}

func isUnusedEventIdFiles(fileName string, curTickRound uint64) bool {
	if !strings.HasPrefix(fileName, eventIdFilePrefixName) {
		return false
	}
	prevTickRound := curTickRound - 1
	_curTickRound := strconv.FormatUint(curTickRound, 10)
	_prevTickRound := strconv.FormatUint(prevTickRound, 10)

	eventFileRound := strings.TrimPrefix(fileName, eventIdFilePrefixName)
	if eventFileRound == _curTickRound || eventFileRound == _prevTickRound {
		return false
	} else {
		return true
	}
}
