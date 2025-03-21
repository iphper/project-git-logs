package src

import "sync"

var (
	once      sync.Once
	group     sync.WaitGroup
	app       *gitLog
	writeChan = make(chan []string, 1)
)
