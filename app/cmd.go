package app

import (
	"github.com/fuxiaohei/GoBlog/app/cmd"
	"os"
)

func Cmd() {
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "install":
			cmd.DoInstall()
		case "update":
			file, _ := cmd.DoBackup(App, false)
			cmd.DoUpdateZipBytes(file)
		case "backup":
			cmd.DoBackup(App, true)
		}
		os.Exit(1)
	}
	// do install and run server together
	if !cmd.CheckInstall() {
		cmd.DoInstall()
	}
}
