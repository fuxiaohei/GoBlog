package cmd

import (
	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/GoBlog/app/model"
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

func installConfig() *gof.Config {
	cfg, _ := gof.NewConfig("")
	cfg.Set("version", vars.VERSION)

	// set web config
	cfg.Set("web.addr", vars.WEB_ADDR)
	cfg.Set("web.port", vars.WEB_PORT)

	// set db config
	cfg.Set("database.driver", vars.DB_DRIVER)
	cfg.Set("database.file", vars.DB_FILE)

	if err := cfg.ToFile(vars.CONFIG_FILE); err != nil {
		log.Fatal("[x] create config json error : %v", err)
		os.Exit(1)
	}
	log.Info("[√] create config json file !")
	return cfg
}

func installDatabase(cfg *gof.Config) {
	if err := model.CreateDB(cfg); err != nil {
		log.Fatal("[x] create database error : %v", err)
		os.Exit(1)
	}
	log.Info("[√] create database !")
}

func doInstall() {
	log.PREFIX = vars.LOG_PREFIX + "[ins]"

	log.Info("begin to install Fxh.Go blog engine")

	cfg := installConfig()

	installDatabase(cfg)

	log.Info("install success !")
}
