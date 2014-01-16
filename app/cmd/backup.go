package cmd

import (
	"os"
	"time"
	"git.oschina.net/fuxiaohei/GoBlog.git/app/utils"
	"path"
	"github.com/Unknwon/cae/zip"
	"git.oschina.net/fuxiaohei/GoBlog.git/GoInk"
	"path/filepath"
)

var backupDir = "backup"

func DoBackup(app *GoInk.App) (string, error) {
	os.Mkdir(backupDir, os.ModePerm)
	filename := path.Join(backupDir, utils.DateTime(time.Now(), "YYYYMMDDHHmmss.zip"))
	z, e := zip.Create(filename)
	if e != nil {
		return "", e
	}
	root, _ := os.Getwd()
	z.AddDir("static", path.Join(root, "static"))
	z.AddDir("data", path.Join(root, "data"))
	z.AddDir(app.View().Dir, path.Join(root, app.View().Dir))
	e = z.Flush()
	if e != nil {
		return "", e
	}
	println("backup success in "+filename)
	return filename, nil
}

func RemoveBackupFile(file string) {
	file = path.Join(backupDir, file)
	os.Remove(file)
}

func GetBackupFiles() ([]os.FileInfo, error) {
	fi := make([]os.FileInfo, 0)
	e := filepath.Walk(backupDir, func(_ string, info os.FileInfo, _ error) error {
			if info == nil{
				return nil
			}
			if !info.IsDir() {
				fi = append([]os.FileInfo{info}, fi...)
			}
			return nil
		})
	return fi, e
}
