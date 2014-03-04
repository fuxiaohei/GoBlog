package upgrade

import (
	"github.com/fuxiaohei/GoBlog/app/cmd"
	"github.com/fuxiaohei/GoBlog/app/model/setting"
	"github.com/fuxiaohei/GoBlog/app/model/storage"
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
	setting.Load()
	setting.Set("c_footer_ga", "<!-- google analytics or other -->")
	setting.Set("enable_go_markdown", "false")
	setting.Set("enable_go_markdown_def", "false")
	setting.Set("site_theme", "default")
	setting.Set("site_theme_def", "default")
	setting.Set("c_home_avatar", "/static/img/site.png")
	setting.Sync()

	// init plugin
	plugin.Init()
	storage.Storage.MkDir("plugin")

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
