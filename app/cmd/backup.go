package cmd

import (
	"fmt"
	"github.com/Unknwon/cae/zip"
	"github.com/fuxiaohei/GoBlog/app/model/message"
	"github.com/fuxiaohei/GoBlog/app/model/setting"
	"github.com/fuxiaohei/GoBlog/app/model/timer"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"github.com/fuxiaohei/GoInk"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var backupDir = "backup"

func init() {
	// close zip terminal output
	zip.Verbose = false
}

// DoBackup backups whole files to zip archive.
// If withData is false, it compresses static files to zip archive without data files, config files and install lock file.
func DoBackup(app *GoInk.App, cs string) (string, error) {
	com := strings.Split(cs, ",")
	os.Mkdir(backupDir, os.ModePerm)
	// create zip file name from time unix
	filename := path.Join(backupDir, utils.DateTime(time.Now(), "YYYYMMDDHHmmss"))
	filename += "_" + strings.Join(com, "_") + ".zip"
	z, e := zip.Create(filename)
	if e != nil {
		return "", e
	}
	root, _ := os.Getwd()
	for _, c := range com {
		if c == "static" {
			z.AddDir("static/css", path.Join(root, "static", "css"))
			z.AddDir("static/img", path.Join(root, "static", "img"))
			z.AddDir("static/js", path.Join(root, "static", "js"))
			z.AddDir("static/lib", path.Join(root, "static", "lib"))
			z.AddFile("static/favicon.ico", path.Join(root, "static", "favicon.ico"))
		}
		if c == "data" {
			lockFile := path.Join(root, "install.lock")
			if utils.IsFile(lockFile) {
				z.AddFile("install.lock", lockFile)
			}
			configFile := path.Join(root, "config.json")
			if utils.IsFile(configFile) {
				z.AddFile("config.json", configFile)
			}
			z.AddDir("data", path.Join(root, "data"))
		}
		if c == "upload" {
			z.AddDir("static/upload", path.Join(root, "static", "upload"))
		}
		if c == "theme" {
			z.AddDir(app.View().Dir, path.Join(root, app.View().Dir))
		}
	}
	e = z.Flush()
	if e != nil {
		return "", e
	}
	println("backup success in " + filename)
	return filename, nil
}

// RemoveBackupFile removes backup zip file with filename(not filepath).
func RemoveBackupFile(file string) {
	file = path.Join(backupDir, file)
	os.Remove(file)
}

// GetBackupFileAbsPath returns backup zip absolute filepath by filename.
func GetBackupFileAbsPath(name string) string {
	return path.Join(backupDir, name)
}

// GetBackupFile returns fileinfo slice of all backup files.
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

// StartBackupTimer starts backup operation timer for auto backup stuff.
func StartBackupTimer(app *GoInk.App, t int) {
	timer.SetFunc("backup-data", 144, func() {
		filename, e := DoBackup(app, setting.Get("backup_setting"))
		if e != nil {
			message.Create("backup", "[0]"+fmt.Sprint(e))
		} else {
			message.Create("backup", "[1]"+filename)
		}
		println("backup files in", t, "hours")
	})
}
