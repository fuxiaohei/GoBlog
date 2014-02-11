package upgrade

import (
	"github.com/fuxiaohei/GoBlog/app/cmd"
	"github.com/fuxiaohei/GoInk"
	"os"
	"path"
)

func init() {
	cmd.SetUpgradeScript(20140131, upgrade_20140131)
}

func upgrade_20140131(app *GoInk.App) bool {

	// re-write all data to non-indent json
	/*model.All()
	model.SyncContents()
	model.SyncFiles()
	model.SyncReaders()
	model.SyncSettings()
	model.SyncTokens()
	model.SyncUsers()
	model.SyncVersion()*/

	// update ling template
	os.RemoveAll(path.Join(app.Get("view_dir"), "ling"))
	cmd.ExtractBundleBytes()

	return true
}
