package handler

import (
	"github.com/fuxiaohei/GoBlog/app/cmd"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoInk"
)

func CmdBackup(context *GoInk.Context) {
	if context.Method == "POST" {
		file, e := cmd.DoBackup(context.App(), true)
		if e != nil {
			Json(context, false).Set("msg", e.Error()).End()
			return
		}
		Json(context, true).Set("file", file).End()
		context.Do("bakcup_success", file)
		model.CreateMessage("backup", "[1]"+file)
		return
	}
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

func CmdBackupFile(context *GoInk.Context) {
	file := context.String("file")
	context.Download(cmd.GetBackupFileAbsPath(file))
	context.Do("backup_download", file)
}

func CmdMessage(context *GoInk.Context) {
	context.Layout("admin/cmd")
	context.Render("admin/cmd/message", map[string]interface{}{
		"Title":    "消息",
		"Messages": model.GetMessages(),
	})
}

func CmdLogs(context *GoInk.Context) {
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

func CmdMonitor(ctx *GoInk.Context) {
	ctx.Layout("admin/cmd")
	ctx.Render("admin/cmd/monitor", map[string]interface{}{
		"Title": "系统监控",
		"M":     cmd.ReadMemStats(),
	})
}
