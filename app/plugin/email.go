package plugin

import "github.com/fuxiaohei/GoBlog/app/model"

type EmailPlugin struct {
	isActive            bool
	isHandlerRegistered bool
	settings            map[string]string
}

func init() {
	EmailPlugin := new(EmailPlugin)
	EmailPlugin.isActive = false
	EmailPlugin.isHandlerRegistered = false
	EmailPlugin.settings = make(map[string]string)
	register(EmailPlugin)
}

func (p *EmailPlugin) Name() string {
	return "邮件提醒"
}

func (p *EmailPlugin) Key() string {
	return "email_notify_plugin"
}

func (p *EmailPlugin) Desc() string {
	return "评论及回复等邮件提醒"
}

func (p *EmailPlugin) ToStorage() map[string]interface{} {
	m := make(map[string]interface{})
	m["name"] = p.Name()
	m["description"] = p.Desc()
	m["is_activate"] = p.isActive
	return m
}

func (p *EmailPlugin) Activate() {
	if p.isHandlerRegistered {
		p.isActive = true
		return
	}
	model.Storage.Get("plugin/"+p.Key(), &p.settings)
	/*fn := func(context *GoInk.Context) {
		now := time.Now()
		context.On(GoInk.CONTEXT_RENDERED, func() {
				if p.isActive {
					duration := time.Since(now)
					context.Body = append(context.Body, []byte(fmt.Sprint("\n<!-- execute ", duration, " -->"))...)
				}
			})
	}
	Handler("hello_plugin", fn, false)
	p.isHandlerRegistered = true*/
	p.isActive = true
	p.isHandlerRegistered = true
}

func (p *EmailPlugin) Deactivate() {
	p.isActive = false
}

func (p *EmailPlugin) IsActive() bool {
	return p.isActive
}

func (p *EmailPlugin) Version() string {
	return "0.1.5"
}

func (p *EmailPlugin) HasSetting() bool {
	return true
}

func (p *EmailPlugin) Form() string {
	return `<h4 class="title">
        SMTP 邮箱提醒设置
    </h4>
    <p class="item">
        <label for="smtp-addr">SMTP服务器</label>
        <input class="ipt" id="smtp-addr" type="text" name="smtp_host" placeholder="SMTP服务器地址，如 smtp.gmail.com:465"/>
    </p>
    <p class="item">
        <label for="smtp-email">SMTP邮箱</label>
        <input class="ipt" id="smtp-email" type="email" name="smtp_email_user" placeholder="使用SMTP的邮箱"/>
    </p>
    <p class="item">
        <label for="smtp-password">邮箱密码</label>
        <input class="ipt" id="smtp-password" type="password" name="smtp_email_password" placeholder="邮箱密码"/>
    </p>
    <p class="submit item">
        <label>&nbsp;</label>
        <button class="btn btn-blue">保存设置</button>
    </p>`
}

func (p *EmailPlugin) SetSetting(settings map[string]string) {
	p.settings = settings
	model.Storage.Set("plugin/"+p.Key(), p.settings)
}
