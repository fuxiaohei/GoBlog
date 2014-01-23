package cmd

import (
	"github.com/fuxiaohei/GoBlog/app/model"
	"strconv"
)

func CheckUpgrade(v int) bool {
	model.Init(v)
	appV := model.GetVersion()
	b := v > appV.Version
	if b {
		println("app version @ " + strconv.Itoa(v) + " is ahead of current version @ " + strconv.Itoa(appV.Version) + " , please run 'GoBlog upgrade'")
	}
	return b
}
