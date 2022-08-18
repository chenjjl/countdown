package timeWheel

import (
	"container/list"
	"countdown/src/storage"
)

type bigHandBucket struct {
	files *list.List
}

func NewBigHandBucket() *bigHandBucket {
	return &bigHandBucket{
		files: list.New(),
	}
}

func (b *bigHandBucket) Add(file *storage.BhFile) error {
	b.files.PushBack(file)
	return nil
}

func (b *bigHandBucket) Lookup() (*storage.BhFile, error) {
	var n *list.Element
	var fileRes *storage.BhFile
	var count int
	for e := b.files.Front(); e != nil; e = n {
		file := (e.Value).(*storage.BhFile)
		n = e.Next()
		if file.CurRound == file.Round {
			b.files.Remove(e)

			count++
			if count > 1 {
				log.Infof("[same Round file] file %+v first file %+v", file, fileRes)
			}
			fileRes = file
		} else {
			file.CurRound += 1
		}
	}
	return fileRes, nil
}

func (b *bigHandBucket) LookupFiles(fileName string) (*storage.BhFile, error) {
	var n *list.Element
	for e := b.files.Front(); e != nil; e = n {
		file := (e.Value).(*storage.BhFile)
		n = e.Next()
		if file.Name == storage.BhFileNamePrefix+fileName {
			return file, nil
		}
	}
	return nil, nil
}
