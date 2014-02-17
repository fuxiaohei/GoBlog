package app

import (
	"github.com/fuxiaohei/GoBlog/app/cmd"
	_ "github.com/fuxiaohei/GoBlog/app/upgrade"
	"os"
)

// Cmd starts command line application.
// It captures command line arguments and executes proper operation.
// Some operations will exit application when finished.
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
		case "upgrade":
			cmd.DoUpgrade(VERSION, App)
		}
		os.Exit(1)
	}
	// do install and run server together
	if !cmd.CheckInstall() {
		cmd.DoInstall()
		return
	}
	// check app version
	if cmd.CheckUpgrade(VERSION) {
		os.Exit(1)
		return
	}

	// begin cmd init
	cmd.Init(App)
}
