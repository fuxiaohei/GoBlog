package upgrade

import (
	"github.com/fuxiaohei/GoBlog/app/cmd"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoBlog/app/plugin"
	"github.com/fuxiaohei/GoInk"
	"os"
	"path"
)

func init() {
	cmd.SetUpgradeScript(20140130, upgrade_20140130)
}

func upgrade_20140130(app *GoInk.App) bool {

	// change settings
	model.LoadSettings()
	model.SetSetting("c_footer_ga", "<!-- google analytics or other -->")
	model.SetSetting("enable_go_markdown", "false")
	model.SetSetting("enable_go_markdown_def", "false")
	model.SetSetting("site_theme", "default")
	model.SetSetting("site_theme_def", "default")
	model.SetSetting("c_home_avatar", "/static/img/site.png")
	model.SyncSettings()

	// init plugin
	plugin.Init()
	model.Storage.Dir("plugin")

	// remove static files
	os.RemoveAll(app.Get("view_dir"))
	os.RemoveAll(path.Join(app.Get("static_dir"), "less"))
	os.RemoveAll(path.Join(app.Get("static_dir"), "css"))
	os.RemoveAll(path.Join(app.Get("static_dir"), "img"))
	os.RemoveAll(path.Join(app.Get("static_dir"), "js"))
	os.RemoveAll(path.Join(app.Get("static_dir"), "lib"))
	os.Remove(path.Join(app.Get("static_dir"), "favicon.ico"))

	// extract current static files
	cmd.ExtractBundleBytes()

	// "c_footer_ga":        "<!-- google analytics or other -->",
	// "enable_go_markdown": "true",
	// "enable_go_markdown_def": "false",
	// "site_theme": "ling",
	// "site_theme_def": "default",
	return true
}
