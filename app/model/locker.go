package model

import "sync"

var locker sync.Mutex

func init() {
	locker = sync.Mutex{}
}
