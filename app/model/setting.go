package model

import (
	"github.com/fuxiaohei/GoBlog/app"
	"fmt"
)

type SettingModel struct {
	settingCache map[string]string
}

func (this *SettingModel) GetAll() map[string]string {
	if len(this.settingCache) > 0 {
		return this.settingCache
	}
	sql := "SELECT * FROM blog_setting"
	res, e := app.Db.Query(sql)
	if e != nil {
		return nil
	}
	this.cacheSetting(res.Data...)
	return this.settingCache
}

func (this *SettingModel) cacheSetting(settings... map[string]string) {
	for _, s := range settings {
		this.settingCache[s["key"]] = s["value"]
	}
}

func (this *SettingModel) SaveSetting(settings... map[string]string) {
	sqlUpdate := "UPDATE blog_setting SET key = ? AND value = ?"
	sqlInsert := "INSERT INTO blog_setting(key,value) VALUES(?,?)"
	for _, s := range settings {
		_, ok := this.settingCache[s["key"]]
		if ok {
			app.Db.Exec(sqlUpdate, s["key"], s["value"])
		}else {
			app.Db.Exec(sqlInsert, s["key"], s["value"])
		}
	}
	this.cacheSetting(settings...)
}

func (this *SettingModel) GetItem(key string) string {
	return this.settingCache[key]
}

func (this *SettingModel) Reset() {
	this.settingCache = make(map[string]string)
	this.GetAll()
	fmt.Println(this.settingCache)
}

func NewSettingModel() *SettingModel {
	s := new(SettingModel)
	s.Reset()
	return s
}
