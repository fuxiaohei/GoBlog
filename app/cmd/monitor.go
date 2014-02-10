package cmd

import (
	"fmt"
	"github.com/fuxiaohei/GoBlog/app/utils"
	"runtime"
	"time"
)

type monitorStats struct {
	NumGoroutine int
	MemAllocated string
	MemMalloc    string
	MemTotal     string
	MemSys       string
	MemHeap      string
	MemGc        string
	LastGcTime   string
}

func ReadMemStats() *monitorStats {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	ms := new(monitorStats)
	ms.NumGoroutine = runtime.NumGoroutine()
	ms.MemAllocated = utils.FileSize(int64(m.Alloc))
	ms.MemTotal = utils.FileSize(int64(m.TotalAlloc))
	ms.MemSys = utils.FileSize(int64(m.Sys))
	ms.MemHeap = utils.FileSize(int64(m.HeapAlloc))
	ms.MemMalloc = utils.FileSize(int64(m.Mallocs))
	ms.LastGcTime = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	ms.MemGc = utils.FileSize(int64(m.NextGC))
	return ms
}
