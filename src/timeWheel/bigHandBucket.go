package timeWheel

import (
	"container/list"
)

type bigHandBucket struct {
	files *list.List
}

func NewBigHandBucket() *bigHandBucket {
	return &bigHandBucket{
		files: list.New(),
	}
}

func (b *bigHandBucket) Add(file *File) error {
	b.files.PushBack(file)
	return nil
}

func (b *bigHandBucket) Lookup() (*File, error) {
	var n *list.Element
	var fileRes *File
	var count int
	for e := b.files.Front(); e != nil; e = n {
		file := (e.Value).(*File)
		n = e.Next()
		if file.curRound == file.round {
			b.files.Remove(e)

			count++
			if count > 1 {
				log.Infof("[same round file] file %+v first file %+v", file, fileRes)
			}
			fileRes = file
		} else {
			file.curRound += 1
		}
	}
	return fileRes, nil
}

func (b *bigHandBucket) LookupFiles(fileName string) (*File, error) {
	var n *list.Element
	for e := b.files.Front(); e != nil; e = n {
		file := (e.Value).(*File)
		n = e.Next()
		if file.name == fileName {
			return file, nil
		}
	}
	return nil, nil
}
