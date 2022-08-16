package timeWheel

import (
	"container/list"
	"countdown/src/logger"
)

var log = logger.GetLogger("timeWheel")

type wheel struct {
	tick        uint64
	wheelSize   uint64
	buckets     *list.List
	curBucket   *list.Element
	bucketIndex uint64
	head        *list.Element
}

type Element struct {
	Round      uint64
	CurRound   uint64
	Expiration uint64
}
