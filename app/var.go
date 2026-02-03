package app

import "sync"

var (
	App *Application
	wg  = &sync.WaitGroup{}
)
