package plugin

import (
	"github.com/fuxiaohei/GoBlog/GoInk"
	"github.com/fuxiaohei/GoBlog/app/model"
	"fmt"
)

type PluginInterface interface {
	Name() string
	Key() string
	Desc() string
	ToStorage() map[string]interface {}
	Activate()
	Deactivate()
	IsActive() bool
	Version() string
}

var (
	pluginStorage map[string]map[string]interface {}
	pluginMap map[string]PluginInterface
	middleHandler map[string]GoInk.Handler
	interHandler map[string]GoInk.Handler
)

func init() {
	if pluginMap == nil {
		pluginMap = make(map[string]PluginInterface)
	}
	//pluginMap = make(map[string]PluginInterface)
	pluginStorage = make(map[string]map[string]interface {})
	middleHandler = make(map[string]GoInk.Handler)
	interHandler = make(map[string]GoInk.Handler)
}

func Init() {
	var isChanged = false
	model.Storage.Get("plugins", &pluginStorage)
	fmt.Println(pluginStorage, pluginMap)
	for k, p := range pluginMap {
		_, ok := pluginStorage[k]
		if !ok {
			pluginStorage[k] = p.ToStorage()
			isChanged = true
		}else {
			if (pluginStorage[k]["is_activate"].(bool)) {
				p.Activate()
			}else {
				p.Deactivate()
			}
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
	}else {
		middleHandler[name] = h
	}
}

func Handlers() map[string]map[string]GoInk.Handler {
	m := make(map[string]map[string]GoInk.Handler)
	m["middle"] = middleHandler
	m["inter"] = interHandler
	return m
}

func Plugins() map[string]PluginInterface {
	return pluginMap
}

func Activate(name string) {
	p, ok := pluginMap[name]
	if !ok {
		println("activate null plugin "+name)
		return
	}
	p.Activate()
	pluginStorage[p.Key()] = p.ToStorage()
	model.Storage.Set("plugins", pluginStorage)
}

func Deactivate(name string) {
	p, ok := pluginMap[name]
	if !ok {
		println("deactivate null plugin "+name)
		return
	}
	p.Deactivate()
	pluginStorage[p.Key()] = p.ToStorage()
	model.Storage.Set("plugins", pluginStorage)
}

