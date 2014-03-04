package setting

import (
	. "github.com/fuxiaohei/GoBlog/app/model/storage"
)

var ver *Version

type Version struct {
	Name           string
	BuildTime      int64
	Version        int
	CurrentVersion int
}

func LoadVersion() {
	ver = new(Version)
	Storage.Get("Version", ver)
}

func GetVersion() *Version {
	if ver == nil {
		LoadVersion()
	}
	return ver
}

func SyncVersion() {
	Storage.Set("version", ver)
}
