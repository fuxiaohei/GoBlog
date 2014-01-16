package model

import "strings"

var settings map[string]string

func GetSetting(key string) string {
	return settings[key]
}

func GetCustomSettings() map[string]string {
	m := make(map[string]string)
	for k, v := range settings {
		if strings.HasPrefix(k, "c_") {
			m[strings.TrimPrefix(k, "c_")] = v
		}
	}
	return m
}

func SetSetting(key string, v string) {
	settings[key] = v
}

func SyncSettings() {
	Storage.Set("settings", settings)
}

func LoadSettings() {
	settings = make(map[string]string)
	Storage.Get("settings", &settings)
}
