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

func SetVersion(v int) {
	ver.CurrentVersion = v
}

func CheckVersion() bool {
	return ver.Version >= ver.CurrentVersion
}
