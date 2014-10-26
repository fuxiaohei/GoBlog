package main

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/GoBlog/cmd"
	"os"
)

const VERSION = "0.3.1 alpha"

func init() {
	cmd.Version = VERSION
}

func main() {
	app := cli.NewApp()
	app.Name = "Fxh.Go"
	app.Usage = "golang blog engine"
	app.Version = VERSION
	app.Commands = []cli.Command{
		*cmd.CmdRegister.Get("web"),
		*cmd.CmdRegister.Get("backup"),
		*cmd.CmdRegister.Get("upgrade"),
	}
	app.Run(os.Args)
}
