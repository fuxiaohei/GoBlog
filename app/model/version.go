package model

var ver *version

type version struct {
	Name           string
	BuildTime      int64
	Version        int
	CurrentVersion int
}

func loadVersion() {
	ver = new(version)
	Storage.Get("version", ver)
}

func GetVersion() *version {
	if ver == nil {
		loadVersion()
	}
	return ver
}

func SyncVersion() {
	Storage.Set("version", ver)
}
