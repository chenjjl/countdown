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
	for e := b.files.Front(); e != nil; e = n {
		file := (e.Value).(*File)
		n = e.Next()
		if file.curRound == file.round {
			b.files.Remove(e)
			return file, nil
		}
		file.curRound += 1
	}
	return nil, nil
}

func (b *bigHandBucket) LookupFiles(round uint64) (*File, error) {
	var n *list.Element
	for e := b.files.Front(); e != nil; e = n {
		file := (e.Value).(*File)
		n = e.Next()
		if file.round == round {
			return file, nil
		}
	}
	return nil, nil
}
