package timeWheel

import (
	"container/list"
)

type wheel struct {
	tick        uint64
	wheelSize   uint64
	buckets     *list.List
	curBucket   *list.Element
	bucketIndex uint64
	head        *list.Element
}

type element struct {
	round      uint64
	curRound   uint64
	Expiration uint64
}
