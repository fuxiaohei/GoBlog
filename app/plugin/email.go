package plugin

import (

)

type EmailPlugin struct {
	isActive            bool
	isHandlerRegistered bool
}

func init() {
	EmailPlugin := new(EmailPlugin)
	EmailPlugin.isActive = false
	EmailPlugin.isHandlerRegistered = false
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

func (p *EmailPlugin) ToStorage() map[string]interface {} {
	m := make(map[string]interface {})
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
