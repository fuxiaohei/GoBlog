
package handler

import "github.com/fuxiaohei/GoInk/Core"

func AdminSetting(context *Core.Context) interface {}{
	context.Render("admin:admin/setting.html",map[string]interface {}{
			"Title":"配置",
			"IsSetting":true,
		})
	return nil
}

