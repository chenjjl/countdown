package storage

import (
	"countdown/src/event"
	"os"
	"strconv"
	"strings"
)

const DirName = "/tmp/countdown" + dirPrefix
const dirPrefix string = "/eventLog/"

func CreateDir() error {
	exist, err := hasDir(DirName)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	err = os.Mkdir(DirName, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}
	return nil
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

func WrapperFileContent(data []byte) []byte {
	data = append([]byte("["), data...)
	data[len(data)-1] = ']'
	return data
}

func encodeEvent(event *event.Event) string {
	var builder strings.Builder
	builder.WriteString(event.Id)
	builder.WriteByte('-')
	builder.WriteString(event.Topic)
	builder.WriteByte('-')
	builder.WriteString(event.Tags)
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatUint(event.TickRound, 10))
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatUint(event.TickOffset, 10))
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatUint(event.CurRound, 10))
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatInt(event.AddBhUnix, 10))
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatUint(event.Expiration, 10))
	builder.WriteByte('-')
	builder.WriteString(strconv.FormatUint(event.Round, 10))
	builder.WriteByte('-')

	return builder.String()
}
