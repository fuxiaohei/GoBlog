package plugin

import (
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoInk"
)

type PluginInterface interface {
	Name() string
	Key() string
	Desc() string
	Version() string

	Activate()
	Deactivate()
	IsActive() bool

	HasSetting() bool
	Form() string
	SetSetting(settings map[string]string)

	ToStorage() map[string]interface{}
}

type pluginRoute struct {
	Method  string
	Pattern string
	Handler GoInk.Handler
}

var (
	pluginStorage map[string]map[string]interface{}
	pluginMap     map[string]PluginInterface
	middleHandler map[string]GoInk.Handler
	interHandler  map[string]GoInk.Handler
	usedHandler   map[string]map[string]bool
	routeHandler  map[string]pluginRoute
)

func init() {
	if pluginMap == nil {
		pluginMap = make(map[string]PluginInterface)
	}
	//pluginMap = make(map[string]PluginInterface)
	pluginStorage = make(map[string]map[string]interface{})
	middleHandler = make(map[string]GoInk.Handler)
	routeHandler = make(map[string]pluginRoute)
	interHandler = make(map[string]GoInk.Handler)
	usedHandler = make(map[string]map[string]bool)
	usedHandler["middle"] = make(map[string]bool)
	usedHandler["inter"] = make(map[string]bool)
	usedHandler["route"] = make(map[string]bool)
}

func Init() {
	var isChanged = false
	if model.Storage.Has("plugins") {
		model.Storage.Get("plugins", &pluginStorage)
	}
	// activate
	for k, p := range pluginMap {
		_, ok := pluginStorage[k]
		if !ok {
			pluginStorage[k] = p.ToStorage()
			isChanged = true
		}
		if pluginStorage[k]["is_activate"].(bool) {
			p.Activate()
		} else {
			p.Deactivate()
		}
	}
	// clean deleted
	for k, _ := range pluginStorage {
		if pluginMap[k] == nil {
			delete(pluginStorage, k)
			isChanged = true
		}
	}
	if isChanged {
		model.Storage.Set("plugins", pluginStorage)
	}
}

func register(plugin PluginInterface) {
	if pluginMap == nil {
		pluginMap = make(map[string]PluginInterface)
	}
	pluginMap[plugin.Key()] = plugin
}

func Handler(name string, h GoInk.Handler, inter bool) {
	if inter {
		interHandler[name] = h
	} else {
		middleHandler[name] = h
	}
}

func Route(name string, method string, pattern string, handler GoInk.Handler) {
	pr := pluginRoute{}
	pr.Method = method
	pr.Handler = handler
	pr.Pattern = pattern
	routeHandler[name] = pr
}

func Handlers() (map[string]map[string]GoInk.Handler, map[string]pluginRoute) {
	m := make(map[string]map[string]GoInk.Handler)
	m["middle"] = middleHandler
	m["inter"] = interHandler
	return m, routeHandler
}

func GetPlugins() map[string]PluginInterface {
	return pluginMap
}

func GetPluginByKey(key string) PluginInterface {
	return pluginMap[key]
}

func Activate(name string) {
	p, ok := pluginMap[name]
	if !ok {
		println("activate null plugin " + name)
		return
	}
	p.Activate()
	pluginStorage[p.Key()] = p.ToStorage()
	model.Storage.Set("plugins", pluginStorage)
	println("activate", p.Key())
}

func Deactivate(name string) {
	p, ok := pluginMap[name]
	if !ok {
		println("deactivate null plugin " + name)
		return
	}
	p.Deactivate()
	pluginStorage[p.Key()] = p.ToStorage()
	model.Storage.Set("plugins", pluginStorage)
	println("deactivate", p.Key())
}

func Update(app *GoInk.App) {
	pluginHandlers, routeHandlers := Handlers()

	if len(routeHandlers) > 0 {
		for n, h := range routeHandlers {
			if usedHandler["route"][n] {
				continue
			}
			app.Route(h.Method, h.Pattern, h.Handler)
			usedHandler["route"][n] = true
		}
	}

	if len(pluginHandlers["middle"]) > 0 {
		for n, h := range pluginHandlers["middle"] {
			if usedHandler["middle"][n] {
				continue
			}
			app.Use(h)
			usedHandler["middle"][n] = true
			//println("use plugin middle handler",n)
		}
		//fmt.Println(usedHandler)
	}

	if len(pluginHandlers["inter"]) > 0 {
		for name, h := range pluginHandlers["inter"] {
			if usedHandler["inter"][name] {
				continue
			}
			if name == "static" {
				app.Static(h)
				usedHandler["inter"][name] = true
				continue
			}
			if name == "recover" {
				app.Recover(h)
				usedHandler["inter"][name] = true
				continue
			}
			if name == "notfound" {
				app.NotFound(h)
				usedHandler["inter"][name] = true
				continue
			}
		}
	}
}
