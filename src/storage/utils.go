package storage

import (
	"os"
)

const DirName = "/tmp/countdown" + dirPrefix
const dirPrefix string = "/eventLog/"

func CreateDir() error {
	if Exists(DirName) {
		return nil
	}
	err := os.Mkdir(DirName, os.ModeDir|os.ModePerm)
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
