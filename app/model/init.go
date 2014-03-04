
package model


func loadAllData() {
	loadVersion()
	LoadSettings()
	LoadNavigators()
	LoadUsers()
	LoadTokens()
	LoadContents()
	LoadMessages()
	LoadReaders()
	LoadComments()
	LoadFiles()
}

// Init does model initialization.
// If first run, write default data.
// v means app.Version number. It's needed for version data.
func Init(v int) {
	appVersion = v
	Storage = new(jsonStorage)
	Storage.Init("data")
	TmpStorage = new(jsonStorage)
	TmpStorage.dir = "tmp/data"
	if !Storage.Has("version") {
		os.Mkdir(Storage.dir, os.ModePerm)
		os.Mkdir(path.Join(Storage.dir, "content"), os.ModePerm)
		os.Mkdir(path.Join(Storage.dir, "plugin"), os.ModePerm)
		writeDefaultData()
	}
}

// All loads all data from storage to memory.
// Start timers for content, comment and message.
func All() {
	loadAllData()
	// generate indexes
	SyncIndexes()
	// start model timer, do all timer stuffs
	StartModelTimer()
}

func SyncIndexes() {
	// generate indexes
	generatePublishArticleIndex()
	generateContentTmpIndexes()
}

// SyncAll writes all current memory data to storage files.
func SyncAll() {
	SyncContents()
	SyncMessages()
	SyncFiles()
	SyncReaders()
	SyncSettings()
	SyncNavigators()
	SyncTokens()
	SyncUsers()
	SyncVersion()
}
