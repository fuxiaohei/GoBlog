package upgrade

import (
	"github.com/fuxiaohei/GoBlog/app/cmd"
	"github.com/fuxiaohei/GoBlog/app/model/setting"
	"github.com/fuxiaohei/GoInk"
)

func init() {
	cmd.SetUpgradeScript(20140228, upgrade_20140228)
}

func upgrade_20140228(_ *GoInk.App) bool {

	// change settings
	setting.Load()
	setting.Set("popular_size", "4")
	setting.Set("recent_comment_size", "3")
	setting.Set("theme_cache", "false")
	setting.Sync()

	// overwrite zip bundle bytes
	cmd.ExtractBundleBytes()
	return true
}
