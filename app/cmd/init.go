package cmd

import "github.com/fuxiaohei/GoInk"

// Init initializes cmd operations.
// Now it starts backup timer.
func Init(app *GoInk.App) {
	StartBackupTimer(app, 24)
}
