package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/GoBlog/gof"
	"github.com/fuxiaohei/GoBlog/gof/log"
)

var CmdWeb = cli.Command{
	Name:        "web",
	Usage:       "start blog web server",
	Description: "start blog web server",
	Action:      DoWeb,
	Flags:       []cli.Flag{},
}

func DoWeb(ctx *cli.Context) {

	server := gof.NewHttpServer()

	addr := server.ConfigInterface.String("web.addr", "0.0.0.0")
	port := server.ConfigInterface.Int("web.port", 8989)

	if err := server.Listen(addr, int(port)); err != nil {
		log.Fatal("%s listen to %s:%d fail - %v", server.LogPrefix, addr, port, err)
	}
}
