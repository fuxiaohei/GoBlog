package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/GoBlog/gof"
	"github.com/fuxiaohei/GoBlog/gof/log"
	"github.com/fuxiaohei/GoBlog/vars"
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
	log.PREFIX = vars.LOG_PREFIX + "[upg]"

	cfg, err := gof.NewConfig(vars.CONFIG_FILE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	version := cfg.String("version")
	log.Info("upgrade process from %s to %s", version, vars.VERSION)
}

func upgradeCheck(ctx *cli.Context) {
	if ctx.Command.Name == "upgrade" {
		return
	}

	log.PREFIX = vars.LOG_PREFIX + "[upg]"

	cfg, err := gof.NewConfig(vars.CONFIG_FILE)
	if err != nil {
		log.Error("%v", err)
		os.Exit(1)
	}
	version := cfg.String("version")
	if version == "" {
		log.Fatal("read blog engine version fail")
		//os.Exit(1)
	}
	if version != vars.VERSION {
		log.Info("the new blog engine version is %s, please upgrade your blog %s", vars.VERSION, version)
		log.Info("run ./blog(.exe) upgrade !!")
		//os.Exit(1)
	}
}
