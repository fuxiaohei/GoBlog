package model

import (
	"strconv"
	"strings"
)

type navItem struct {
	Order int
	Text  string
	Title string
	Link  string
}

var (
	settings   map[string]string
	navigators []*navItem
)

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

func SortNavigators() {
	l := len(navigators)
	for i := 1; i < l; i++ {
		for j := i; j > 0; j-- {
			if navigators[j].Order < navigators[j-1].Order {
				navigators[j], navigators[j-1] = navigators[j-1], navigators[j]
			}
		}
	}
}

func LoadNavigators() {
	navigators = make([]*navItem, 0)
	Storage.Get("navigators", &navigators)
	SortNavigators()
}

func SetNavigators(order []string, text []string, title []string, link []string) {
	navs := make([]*navItem, len(text))
	for i, t := range text {
		n := new(navItem)
		n.Order, _ = strconv.Atoi(order[i])
		n.Text = t
		n.Title = title[i]
		n.Link = link[i]
		navs[i] = n
	}
	navigators = navs
	SyncNavigators()
}

func DefaultNavigators() {
	n := new(navItem)
	n.Order = 1
	n.Text = "文章"
	n.Title = "文章"
	n.Link = "/"
	n2 := new(navItem)
	n2.Order = 2
	n2.Text = "关于"
	n2.Title = "关于"
	n2.Link = "/about-me.html"
	n3 := new(navItem)
	n3.Order = 3
	n3.Text = "好友"
	n3.Title = "好友"
	n3.Link = "/friends.html"
	navigators = []*navItem{n, n2, n3}
	Storage.Set("navigators", navigators)
}

func SyncNavigators() {
	Storage.Set("navigators", navigators)
	SortNavigators()
}

func GetNavigators() []*navItem {
	return navigators
}
