package handler

import (
	"github.com/fuxiaohei/GoBlog/app/cmd"
	"github.com/fuxiaohei/GoBlog/app/model/content"
	"github.com/fuxiaohei/GoBlog/app/model/message"
	"github.com/fuxiaohei/GoBlog/app/model/setting"
	"github.com/fuxiaohei/GoInk"
)

// CmdBackup is backup list and operation page, pattern /cmd/backup/.
func CmdBackup(context *GoInk.Context) {
	// backup in manual
	if context.Method == "POST" {
		file, e := cmd.DoBackup(context.App(), []string{"static", "data", "upload", "theme"})
		if e != nil {
			Json(context, false).Set("msg", e.Error()).End()
			return
		}
		Json(context, true).Set("file", file).End()
		context.Do("bakcup_success", file)
		message.Create("backup", "[1]"+file)
		return
	}
	// delete backup file
	if context.Method == "DELETE" {
		file := context.String("file")
		if file == "" {
			Json(context, false).End()
			return
		}
		cmd.RemoveBackupFile(file)
		Json(context, true).End()
		context.Do("backup_delete", file)
		return
	}
	files, _ := cmd.GetBackupFiles()
	context.Layout("admin/cmd")
	context.Render("admin/cmd/backup", map[string]interface{}{
		"Files": files,
		"Title": "备份",
	})
}

// CmdBackupFile is downloading page for backup file, pattern /cmd/backup/file/.
func CmdBackupFile(context *GoInk.Context) {
	file := context.String("file")
	context.Download(cmd.GetBackupFileAbsPath(file))
	context.Do("backup_download", file)
}

// CmdMessage is message list page, pattern /cmd/message/.
func CmdMessage(context *GoInk.Context) {
	context.Layout("admin/cmd")
	context.Render("admin/cmd/message", map[string]interface{}{
		"Title":    "消息",
		"Messages": message.All(),
	})
}

// CmdLogs is logs list and operation page, pattern /cmd/logs/.
func CmdLogs(context *GoInk.Context) {
	// delete log item
	if context.Method == "DELETE" {
		cmd.RemoveLogFile(context.App(), context.String("file"))
		Json(context, true).End()
		return
	}
	context.Layout("admin/cmd")
	context.Render("admin/cmd/log", map[string]interface{}{
		"Title": "日志",
		"Logs":  cmd.GetLogs(context.App()),
	})
}

// CmdMonitor is app monitor data page, pattern /cmd/monitor/.
func CmdMonitor(ctx *GoInk.Context) {
	ctx.Layout("admin/cmd")
	ctx.Render("admin/cmd/monitor", map[string]interface{}{
		"Title": "系统监控",
		"M":     cmd.ReadMemStats(),
	})
}

// CmdTheme is theme list and operation page, pattern /cmd/theme/.
func CmdTheme(ctx *GoInk.Context) {
	if ctx.Method == "POST" {
		// set theme cache
		change := ctx.String("cache")
		if change != "" {
			cmd.SetThemeCache(ctx, change == "true")
			Json(ctx, true).End()
			return
		}
		// change theme
		theme := ctx.String("theme")
		if theme != "" {
			setting.Set("site_theme", theme)
			setting.Sync()
			Json(ctx, true).End()
			return
		}
		return
	}
	ctx.Layout("admin/cmd")
	ctx.Render("admin/cmd/theme", map[string]interface{}{
		"Title":        "主题",
		"Themes":       cmd.GetThemes(ctx.App().Get("view_dir")),
		"CurrentTheme": setting.Get("site_theme"),
	})
}

// CmdReader is reader list page, pattern /cmd/reader/.
func CmdReader(ctx *GoInk.Context) {
	if ctx.Method == "POST" {
		email := ctx.String("email")
		content.RemoveReader(email)
		Json(ctx, true).End()
		return
	}
	ctx.Layout("admin/cmd")
	ctx.Render("admin/cmd/reader", map[string]interface{}{
		"Title":   "读者",
		"Readers": content.GetReaders(),
	})
}
