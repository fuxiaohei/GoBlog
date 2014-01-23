package plugin

import (
	"github.com/fuxiaohei/GoBlog/GoInk"
	"github.com/fuxiaohei/GoBlog/app/model"
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

var (
	pluginStorage map[string]map[string]interface{}
	pluginMap     map[string]PluginInterface
	middleHandler map[string]GoInk.Handler
	interHandler  map[string]GoInk.Handler
)

func init() {
	if pluginMap == nil {
		pluginMap = make(map[string]PluginInterface)
	}
	//pluginMap = make(map[string]PluginInterface)
	pluginStorage = make(map[string]map[string]interface{})
	middleHandler = make(map[string]GoInk.Handler)
	interHandler = make(map[string]GoInk.Handler)
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

func Handlers() map[string]map[string]GoInk.Handler {
	m := make(map[string]map[string]GoInk.Handler)
	m["middle"] = middleHandler
	m["inter"] = interHandler
	return m
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
}
