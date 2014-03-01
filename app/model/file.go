package model

import (
	"github.com/fuxiaohei/GoBlog/app/utils"
	"os"
	"path"
	"strconv"
)

var (
	files     []*File
	fileMaxId int
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
}

// CreateFile saves a file instance to json storage.
func CreateFile(f *File) *File {
	fileMaxId += Storage.TimeInc(3)
	f.Id = fileMaxId
	f.UploadTime = utils.Now()
	f.IsUsed = true
	f.Hits = 0
	files = append([]*File{f}, files...)
	go SyncFiles()
	return f
}

// CreateFilePath generates a file path for new uploading file.
func CreateFilePath(dir string, f *File) string {
	os.MkdirAll(dir, os.ModePerm)
	name := utils.DateInt64(utils.Now(), "YYYYMMDDHHmmss")
	name += strconv.Itoa(Storage.TimeInc(10)) + path.Ext(f.Name)
	return path.Join(dir, name)
}

// GetFileList returns a uploaded file instance list with page and size int.
func GetFileList(page, size int) ([]*File, *utils.Pager) {
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

// RemoveFile removes file by id.
func RemoveFile(id int) {
	for i, f2 := range files {
		if id == f2.Id {
			files = append(files[:i], files[i+1:]...)
			os.Remove(f2.Url)
		}
	}
	go SyncFiles()
}

// SyncFiles saves all files data to json storage.
func SyncFiles() {
	Storage.Set("files", files)
}

// LoadFiles loads all files data from json storage.
func LoadFiles() {
	files = make([]*File, 0)
	fileMaxId = 0
	Storage.Get("files", &files)
	for _, f := range files {
		if f.Id > fileMaxId {
			fileMaxId = f.Id
		}
	}
}
