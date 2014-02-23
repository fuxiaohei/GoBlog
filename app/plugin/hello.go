package plugin

import (
	"fmt"
	"github.com/fuxiaohei/GoInk"
	"time"
)

type HelloPlugin struct {
	isActive            bool
	isHandlerRegistered bool
}

func init() {
	helloPlugin := new(HelloPlugin)
	helloPlugin.isActive = true
	helloPlugin.isHandlerRegistered = false
	register(helloPlugin)
}

func (p *HelloPlugin) Name() string {
	return "Sample Hello Plugin"
}

func (p *HelloPlugin) Key() string {
	return "hello_plugin"
}

func (p *HelloPlugin) Desc() string {
	return "插件样例，页面最后输出执行时间 <!-- excute time --> 注释"
}

func (p *HelloPlugin) ToStorage() map[string]interface{} {
	m := make(map[string]interface{})
	m["name"] = p.Name()
	m["description"] = p.Desc()
	m["is_activate"] = p.isActive
	return m
}

func (p *HelloPlugin) Activate() {
	if p.isHandlerRegistered {
		p.isActive = true
		return
	}
	fn := func(context *GoInk.Context) {
		now := time.Now()
		context.On(GoInk.CONTEXT_RENDERED, func() {
			if p.isActive {
				duration := time.Since(now)
				str := fmt.Sprint(duration)
				//context.Body = append(context.Body, []byte(str)...)
				context.Header["X-Exec-Time"] = str
			}
		})
	}
	Handler("hello_plugin", fn, false)
	/*Route("hello_handler", "GET", "/hello/", func(context *GoInk.Context) {
		context.Body = []byte("Hello!")
	})*/
	p.isHandlerRegistered = true
	p.isActive = true
}

func (p *HelloPlugin) Deactivate() {
	p.isActive = false
}

func (p *HelloPlugin) IsActive() bool {
	return p.isActive
}

func (p *HelloPlugin) Version() string {
	return "0.0.1"
}

func (p *HelloPlugin) HasSetting() bool {
	return false
}

func (p *HelloPlugin) Form() string {
	return ""
}

func (p *HelloPlugin) SetSetting(settings map[string]string) {

}
