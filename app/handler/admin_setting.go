package handler

import (
	"github.com/fuxiaohei/GoInk/Core"
	"github.com/fuxiaohei/GoBlog/app/model"
)

func AdminSetting(context *Core.Context) interface {} {
	context.Render("admin:admin/setting.html", map[string]interface {}{
			"Title":"配置",
			"IsSetting":true,
		})
	return nil
}

func AdminSettingPost(context *Core.Context) interface {} {
	data := context.Input()
	settingsMap := make([]map[string]string, len(data))
	i := 0
	for key, v := range data {
		settingsMap[i] = map[string]string{
			"key":key,
			"value":v,
		}
		i++
	}
	model.SettingM.SaveSetting(settingsMap...)
	context.Json(map[string]interface {}{
		"res":true,
	})
	return nil
}

