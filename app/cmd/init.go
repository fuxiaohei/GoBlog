package cmd

import "github.com/fuxiaohei/GoInk"

func Init(app *GoInk.App) {
	StartBackupTimer(app, 24)
}
