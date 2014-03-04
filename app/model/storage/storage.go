package storage

import (
	"encoding/json"
	"github.com/fuxiaohei/GoBlog/app/model/comment"
	"github.com/fuxiaohei/GoBlog/app/model/content"
	"github.com/fuxiaohei/GoBlog/app/model/message"
	"github.com/fuxiaohei/GoBlog/app/model/setting"
	"github.com/fuxiaohei/GoBlog/app/model/user"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

var (
	AppVersion int
	// global data storage instance
	Storage *JsonStorage
	// global tmp data storage instance. Temp data are generated for special usages, will not backup.
	TmpStorage *JsonStorage
)

type JsonStorage struct {
	dir string
}

func (jss *JsonStorage) Init(dir string) {
	jss.dir = dir
}

func (jss *JsonStorage) Has(key string) bool {
	file := path.Join(jss.dir, key+".json")
	_, e := os.Stat(file)
	return e == nil
}

func (jss *JsonStorage) Get(key string, v interface{}) {
	file := path.Join(jss.dir, key+".json")
	bytes, e := ioutil.ReadFile(file)
	if e != nil {
		println("read storage '" + key + "' error")
		return
	}
	e = json.Unmarshal(bytes, v)
	if e != nil {
		println("json decode '" + key + "' error")
	}
}

func (jss *JsonStorage) Set(key string, v interface{}) {
	locker.Lock()
	defer locker.Unlock()

	bytes, e := json.Marshal(v)
	if e != nil {
		println("json encode '" + key + "' error")
		return
	}
	file := path.Join(jss.dir, key+".json")
	e = ioutil.WriteFile(file, bytes, 0777)
	if e != nil {
		println("write storage '" + key + "' error")
	}
}

func (jss *JsonStorage) Dir(name string) {
	os.MkdirAll(path.Join(jss.dir, name), os.ModePerm)
}

// TimeInc returns time step value devided by d int with time unix stamp.
func (jss *JsonStorage) TimeInc(d int) int {
	return int(utils.Now())%d + 1
}
