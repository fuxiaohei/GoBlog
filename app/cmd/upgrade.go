package cmd

import (
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoInk"
	"sort"
	"strconv"
)

var upgradeScript map[int]func(app *GoInk.App) bool

func init() {
	upgradeScript = make(map[int]func(app *GoInk.App) bool)
}

func SetUpgradeScript(v int, script func(app *GoInk.App) bool) {
	upgradeScript[v] = script
}

func CheckUpgrade(v int, print bool) bool {
	model.Init(v)
	appV := model.GetVersion()
	b := v > appV.Version
	if b && print {
		println("app version @ " + strconv.Itoa(v) + " is ahead of current version @ " + strconv.Itoa(appV.Version) + " , please run 'GoBlog upgrade'")
	}
	return b
}

func DoUpgrade(v int, app *GoInk.App) {
	if !CheckUpgrade(v, false) {
		println("app version @", v, "is updated")
		return
	}
	oldVersion := model.GetVersion().Version
	scriptIndex := []int{}
	for vr, _ := range upgradeScript {
		if vr <= v && vr > oldVersion {
			scriptIndex = append(scriptIndex, vr)
		}
	}
	sort.Sort(sort.IntSlice(scriptIndex))
	for _, cv := range scriptIndex {
		upgradeScript[cv](app)
		println("upgrade @", cv, "success")
	}
	model.GetVersion().Version = v
	model.SyncVersion()
	println("app has upgraded to version @", v, "successfully, restart and keep enjoy !!")
}
