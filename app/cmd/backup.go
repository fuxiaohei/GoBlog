package cmd

import (
	"github.com/Unknwon/cae/zip"
	"github.com/fuxiaohei/GoBlog/GoInk"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"os"
	"path"
	"path/filepath"
	"time"
)

var backupDir = "backup"

func DoBackup(app *GoInk.App, withData bool) (string, error) {
	os.Mkdir(backupDir, os.ModePerm)
	filename := path.Join(backupDir, utils.DateTime(time.Now(), "YYYYMMDDHHmmss"))
	if withData {
		filename += ".zip"
	} else {
		filename += "_static.zip"
	}
	z, e := zip.Create(filename)
	if e != nil {
		return "", e
	}
	root, _ := os.Getwd()
	z.AddDir("static/css", path.Join(root, "static", "css"))
	z.AddDir("static/img", path.Join(root, "static", "img"))
	z.AddDir("static/js", path.Join(root, "static", "js"))
	z.AddDir("static/lib", path.Join(root, "static", "lib"))
	z.AddFile("static/favicon.ico", path.Join(root, "static", "favicon.ico"))
	if withData {
		z.AddDir("data", path.Join(root, "data"))
		z.AddDir("static/upload", path.Join(root, "static", "upload"))
	}
	z.AddDir(app.View().Dir, path.Join(root, app.View().Dir))
	e = z.Flush()
	if e != nil {
		return "", e
	}
	println("backup success in " + filename)
	return filename, nil
}

func RemoveBackupFile(file string) {
	file = path.Join(backupDir, file)
	os.Remove(file)
}

func GetBackupFiles() ([]os.FileInfo, error) {
	fi := make([]os.FileInfo, 0)
	e := filepath.Walk(backupDir, func(_ string, info os.FileInfo, _ error) error {
		if info == nil {
			return nil
		}
		if !info.IsDir() {
			fi = append([]os.FileInfo{info}, fi...)
		}
		return nil
	})
	return fi, e
}
