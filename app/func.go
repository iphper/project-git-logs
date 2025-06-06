package app

import "sync"

var (
	app  *App
	once sync.Once
)

// @func Single 单例
func Single() *App {
	once.Do(func() {
		app = new(App)
		app.Init()
	})
	return app
}
