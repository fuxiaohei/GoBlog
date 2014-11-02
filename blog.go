package main

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/GoBlog/cmd"
	"github.com/fuxiaohei/GoBlog/gof/log"
	"github.com/fuxiaohei/GoBlog/vars"
	"os"
)

func init() {
	log.PREFIX = vars.LOG_PREFIX
}

func main() {
	app := cli.NewApp()
	app.Name = "Fxh.Go"
	app.Usage = "golang blog engine"
	app.Version = vars.VERSION
	app.Commands = []cli.Command{
		*cmd.CmdRegister.Get("web"),
		*cmd.CmdRegister.Get("backup"),
		*cmd.CmdRegister.Get("upgrade"),
	}
	app.Run(os.Args)
}
