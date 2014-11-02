package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/GoBlog/gof"
	"github.com/fuxiaohei/GoBlog/gof/log"
	"github.com/fuxiaohei/GoBlog/gof/middle"
	"github.com/fuxiaohei/GoBlog/vars"
)

var CmdWeb = cli.Command{
	Name:        "web",
	Usage:       "start blog web server",
	Description: "start blog web server",
	Action:      DoWeb,
	Flags:       []cli.Flag{},
}

func DoWeb(ctx *cli.Context) {

	log.PREFIX = vars.LOG_PREFIX + "[web]"

	server := gof.NewHttpServer(vars.CONFIG_FILE)

	server.Use(middle.Static(nil))

	addr := server.ConfigInterface.String("web.addr", "0.0.0.0")
	port := server.ConfigInterface.Int("web.port", 8989)

	if err := server.Listen(addr, int(port)); err != nil {
		log.Fatal("listen to %s:%d fail - %v", addr, port, err)
	}
}
