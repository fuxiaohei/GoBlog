package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/GoBlog/gof"
	"github.com/fuxiaohei/GoBlog/gof/log"
	"os"
)

var CmdUpgrade = cli.Command{
	Name:        "upgrade",
	Usage:       "upgrade blog engine",
	Description: "upgrade blog engine",
	Action:      DoUpgrade,
	Flags:       []cli.Flag{},
}

func DoUpgrade(ctx *cli.Context) {
	cfg, err := gof.NewConfig("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	version := cfg.String("version")
	log.Info("upgrade -- upgrade process from %s to %s", version, Version)
}

func upgradeCheck(ctx *cli.Context) {
	if ctx.Command.Name == "upgrade" {
		return
	}
	cfg, err := gof.NewConfig("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	version := cfg.String("version")
	if version == "" {
		fmt.Println("read blog engine version fail")
		os.Exit(1)
	}
	if version != Version {
		fmt.Printf("the new blog engine version is %s, please upgrade your blog %s\n", Version, version)
		fmt.Println("run ./blog(.exe) upgrade !!")
		os.Exit(1)
	}
}
