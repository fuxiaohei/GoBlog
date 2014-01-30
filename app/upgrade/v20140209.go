package upgrade

import (
	"github.com/fuxiaohei/GoBlog/app/cmd"
	"github.com/fuxiaohei/GoBlog/GoInk"
)

func init() {
	cmd.SetUpgradeScript(20140209, upgrade_20140209)
}

func upgrade_20140209(app *GoInk.App) bool {
	return false
}
