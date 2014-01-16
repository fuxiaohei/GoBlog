package app

import (
	"os"
	"git.oschina.net/fuxiaohei/GoBlog.git/app/cmd"
)

func Cmd() {
	args := os.Args
	if len(args) > 1 {
		switch args[1]{
		case "install":
			cmd.DoInstall()
		case "update":
			file, _ := cmd.DoBackup(App)
			cmd.DoUpdateZipBytes(file)
		case "backup":
			cmd.DoBackup(App)
		}
		os.Exit(1)
	}
	// do install and run server together
	if !cmd.CheckInstall(){
		cmd.DoInstall()
	}
}
