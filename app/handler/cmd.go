package handler

import (
	"git.oschina.net/fuxiaohei/GoBlog.git/GoInk"
	"git.oschina.net/fuxiaohei/GoBlog.git/app/cmd"
)

func CmdBackup(context *GoInk.Context) {
	if context.Method == "POST" {
		file , e := cmd.DoBackup(context.App())
		if e != nil {
			Json(context, false).Set("msg", e.Error()).End()
			return
		}
		Json(context, true).Set("file", file).End()
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
		return
	}
	files, _ := cmd.GetBackupFiles()
	context.Layout("cmd")
	context.Render("cmd/backup", map[string]interface {}{
			"Files":files,
		})
}
