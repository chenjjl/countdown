package storage

import (
	"os"
)

const dirName = "/tmp/countdown" + dirPrefix
const dirPrefix string = "/eventLog/"

func CreateDir() error {
	exist, err := hasDir(dirName)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	err = os.Mkdir(dirName, os.ModeDir|os.ModePerm)
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
