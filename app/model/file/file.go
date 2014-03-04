package file

import (
	"fmt"
	. "github.com/fuxiaohei/GoBlog/app/model/storage"
	"github.com/fuxiaohei/GoBlog/app/model/timer"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"os"
	"path"
	"strconv"
)

var (
	files        []*File
	fileMaxId    int
	FileLinkMode = "local"
)

// File struct contains file name, type and upload time data.
type File struct {
	Id          int
	Name        string
	UploadTime  int64
	Url         string
	ContentType string
	Author      int
	IsUsed      bool
	Size        int64
	Type        string
	Hits        int
	Links       map[string]string
}

// Link returns top level file path for uploaded file.
func (f *File) Link() string {
	if _, ok := f.Links[FileLinkMode]; ok {
		return f.Links[FileLinkMode]
	}
	return f.Links["local"]
}

func (f *File) getLocalLink() string {
	return fmt.Sprintf("/upload/%d/%s", f.Id, path.Base(f.Url))
}

// CreateFile saves a file instance to json storage.
func Create(f *File) *File {
	fileMaxId += Storage.TimeInc(3)
	f.Id = fileMaxId
	f.UploadTime = utils.Now()
	f.IsUsed = true
	f.Hits = 0
	f.Links = make(map[string]string)
	f.Links["local"] = f.getLocalLink()
	files = append([]*File{f}, files...)
	go Sync()
	return f
}

// CreateFilePath generates a file path for new uploading file.
func CreatePath(dir string, f *File) string {
	os.MkdirAll(dir, os.ModePerm)
	name := utils.DateInt64(utils.Now(), "YYYYMMDDHHmmss")
	name += strconv.Itoa(Storage.TimeInc(10)) + path.Ext(f.Name)
	return path.Join(dir, name)
}

// GetFileList returns a uploaded file instance list with page and size int.
func List(page, size int) ([]*File, *utils.Pager) {
	pager := utils.NewPager(page, size, len(files))
	f := make([]*File, 0)
	if page > pager.Pages || len(files) < 1 {
		return f, pager
	}
	for i := pager.Begin; i <= pager.End; i++ {
		f = files[pager.Begin-1 : pager.End]
	}
	return f, pager
}

// GetFileById returns a file instance by given id.
func ById(id int) *File {
	for _, f := range files {
		if f.Id == id {
			return f
		}
	}
	return nil
}

// RemoveFile removes file by id.
func Remove(id int) {
	for i, f2 := range files {
		if id == f2.Id {
			files = append(files[:i], files[i+1:]...)
			os.Remove(f2.Url)
		}
	}
	go Sync()
}

// SyncFiles saves all files data to json storage.
func Sync() {
	Storage.Set("files", files)
}

// LoadFiles loads all files data from json storage.
func Load() {
	files = make([]*File, 0)
	fileMaxId = 0
	Storage.Get("files", &files)
	for _, f := range files {
		if f.Id > fileMaxId {
			fileMaxId = f.Id
		}
	}
	startFileSyncTimer()
}

func Len()int{
	return len(files)
}

func startFileSyncTimer() {
	timer.SetFunc("files-sync", 72, func() {
		println("write media in 12 hour timer")
		Sync()
	})
}
