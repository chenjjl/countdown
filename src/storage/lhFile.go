package storage

import (
	"bufio"
	"countdown/src/event"
	"io"
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
		log.Errorf("can not create file for little hand time wheel, because failed to create dir %s", dirName)
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
	return f.addEvent(dirName+f.Name+strconv.FormatUint(tickRound, 10), event)
}

func (f *LhFile) reloadEvents() error {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Errorf("failed to read dir %s", dirName)
		return err
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

	var eventMap map[string]*event.Event

	file, err := os.Open(dirName + f.Name + strconv.FormatInt(maxRound, 10))
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
		eventMap[e.Id] = e
	}
}
