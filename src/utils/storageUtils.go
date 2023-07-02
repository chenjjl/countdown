package utils

import (
	"os"
)

const EventLogDir = prefixDir + "/eventLog/"
const LogDir string = prefixDir + "/log/"

const prefixDir string = "/tmp/countdown"

func CreateDir(fileName string) error {
	if Exists(fileName) {
		return nil
	}
	err := os.Mkdir(fileName, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
