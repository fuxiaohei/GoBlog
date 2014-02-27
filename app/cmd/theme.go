package cmd

import (
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoInk"
	"io/ioutil"
	"path/filepath"
)

var adminTheme = "admin"

type themeItem struct {
	Name       string
	Files      []string
	ErrorFiles []string
	Layout     []string
}

func SetThemeCache(ctx *GoInk.Context, cache bool) {
	ctx.App().View().NoCache()
	ctx.App().View().IsCache = cache
	if cache {
		model.SetSetting("theme_cache", "true")
	} else {
		model.SetSetting("theme_cache", "false")
	}
	model.SyncSettings()
}

func GetThemes(dir string) map[string]*themeItem {
	m := make(map[string]*themeItem)
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		panic(e)
	}
	for _, fi := range files {
		if fi.IsDir() && fi.Name() != adminTheme {
			theme, e := createThemeItem(filepath.Join(dir, fi.Name()))
			if e != nil {
				continue
			}
			theme.Name = fi.Name()
			m[fi.Name()] = theme
		}
	}
	return m
}

func createThemeItem(dir string) (*themeItem, error) {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	theme := new(themeItem)
	theme.Files = make([]string, 0)
	theme.Layout = make([]string, 0)
	for _, fi := range files {
		if fi.IsDir() {
			if fi.Name() == "error" {
				theme.ErrorFiles, _ = filepath.Glob(filepath.Join(dir, fi.Name(), "*.html"))
				for i, f := range theme.ErrorFiles {
					theme.ErrorFiles[i] = filepath.Join(fi.Name(), filepath.Base(f))
				}
			} else {
				f, _ := filepath.Glob(filepath.Join(dir, fi.Name(), "*.html"))
				for _, ff := range f {
					theme.Files = append(theme.Files, filepath.Join(fi.Name(), filepath.Base(ff)))
				}
			}
		} else {
			ext := filepath.Ext(fi.Name())
			if ext == ".html" {
				theme.Files = append(theme.Files, fi.Name())
				continue
			}
			if ext == ".layout" {
				theme.Layout = append(theme.Layout, fi.Name())
			}
		}
	}
	return theme, nil
}
