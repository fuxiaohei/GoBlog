package cmd

import (
	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/GoBlog/gof"
	"github.com/fuxiaohei/GoBlog/gof/log"
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
	log.Info("install -- begin to install Fxh.Go blog engine")
	cfg, _ := gof.NewConfig("")

	cfg.Set("version", Version)

	if err := cfg.ToFile("config.json"); err != nil {
		log.Fatal("install -- install error : %v", err)
		os.Exit(1)
	}

	log.Info("install -- success !")
}
