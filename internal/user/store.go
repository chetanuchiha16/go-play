package user

import "sync"

var (
	CacheMutex sync.RWMutex
	UserCache  = make(map[int]User)
)
