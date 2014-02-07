package cmd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/Unknwon/cae/zip"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"io/ioutil"
	"os"
)

var (
	tmpZipFile      = "tmp.zip"
	installLockFile = "install.lock"
)

func CheckInstall() bool {
	_, e := os.Stat(installLockFile)
	return e == nil
}

func ExtractBundleBytes() {
	// origin from https://github.com/wendal/gor/blob/master/gor/gor.go
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(zipBytes))
	b, _ := ioutil.ReadAll(decoder)
	ioutil.WriteFile(tmpZipFile, b, os.ModePerm)
	z, e := zip.Open(tmpZipFile)
	if e != nil {
		panic(e)
		os.Exit(1)
	}
	z.ExtractTo("")
	defer func() {
		z.Close()
		decoder = nil
		os.Remove(tmpZipFile)
	}()
}

func DoInstall() {
	ExtractBundleBytes()
	ioutil.WriteFile(installLockFile, []byte(fmt.Sprint(utils.Now())), os.ModePerm)
	println("install success")
}

func DoUpdateZipBytes(file string) error {
	// copy from https://github.com/wendal/gor/blob/master/gor/gor.go
	bytes, _ := ioutil.ReadFile(file)
	zipWriter, _ := os.OpenFile("app/cmd/zip.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	header := `package cmd
const zipBytes="`
	zipWriter.Write([]byte(header))
	encoder := base64.NewEncoder(base64.StdEncoding, zipWriter)
	encoder.Write(bytes)
	encoder.Close()
	zipWriter.Write([]byte(`"`))
	zipWriter.Sync()
	zipWriter.Close()
	println("update success")
	return nil
}
