package upgrade

import (
	"github.com/fuxiaohei/GoBlog/app/cmd"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoInk"
	"os"
	"path"
)

func init() {
	cmd.SetUpgradeScript(20140209, upgrade_20140209)
}

func upgrade_20140209(app *GoInk.App) bool {
	// clean template
	vDir := app.Get("view_dir")
	os.Remove(path.Join(vDir, "admin.layout"))
	os.Remove(path.Join(vDir, "cmd.layout"))

	// write default menu setting
	model.DefaultNavigators()

	// write message storage
	model.Storage.Set("messages",[]*model.Message{})

	cmd.ExtractBundleBytes()
	return true
}
