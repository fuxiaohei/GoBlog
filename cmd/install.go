package cmd

import (
	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/GoBlog/gof"
	"github.com/fuxiaohei/GoBlog/gof/log"
	"github.com/fuxiaohei/GoBlog/vars"
	"os"
)

func installCheck(_ *cli.Context) {
	configFile := "config.json"
	if com.IsFile(configFile) {
		return
	}
	doInstall()
}

func doInstall() {
	log.PREFIX = vars.LOG_PREFIX + "[ins]"

	log.Info("begin to install Fxh.Go blog engine")

	cfg, _ := gof.NewConfig("")
	cfg.Set("version", vars.VERSION)

	cfg.Set("web.addr", "0.0.0.0")
	cfg.Set("web.port", 8989)

	if err := cfg.ToFile(vars.CONFIG_FILE); err != nil {
		log.Fatal("install error : %v", err)
		os.Exit(1)
	}

	log.Info("install success !")
}
