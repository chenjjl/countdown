package storage

import (
	"bufio"
	"countdown/src/event"
	"io"
	"os"
	"strconv"
	"strings"
)

const eventIdFilePrefixName = "event_id-"

type IdFile struct {
	Name string
}

type idItem struct {
	Id        string
	TickRound uint64
}

// CreateEventIdFile create file for recording event Id in the current little hand time wheel
func CreateEventIdFile() (*IdFile, error) {
	err := CreateDir()
	if err != nil {
		log.Errorf("can not create file for little hand time wheel, because failed to create dir %s", DirName)
	}
	file, err := newEventIdFile()
	if err != nil {
		return nil, err
	}
	return file, nil
}

func newEventIdFile() (*IdFile, error) {
	return &IdFile{
		Name: eventIdFilePrefixName,
	}, nil
}

func (f *IdFile) Add(event *event.Event) error {
	file, err := os.OpenFile(DirName+f.Name+strconv.FormatUint(event.TickRound, 10), os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)
	if err != nil {
		return err
	}
	item := &idItem{Id: event.Id, TickRound: event.TickRound}
	_, err = file.WriteString(item.Encode() + ",")
	if err != nil {
		return err
	}
	log.Infof("Id %+v be added to idFile %+v", item, f)
	return nil
}

func getIdItems(fileName string, idItemMap map[string]*idItem) error {
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
			log.Errorf("failed to read line from id file")
			return err
		}
		item := Decode(line)
		idItemMap[item.Id] = item
	}
}

func (i *idItem) Encode() string {
	var builder strings.Builder
	builder.WriteString(i.Id)
	builder.WriteByte('_')
	builder.WriteString(strconv.FormatUint(i.TickRound, 10))

	return builder.String()
}

func Decode(s string) *idItem {
	data := strings.Split(s, "_")
	idItem := &idItem{}
	idItem.Id = data[0]
	idItem.TickRound, _ = strconv.ParseUint(data[1], 10, 64)

	return idItem
}
