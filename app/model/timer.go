package model

import (
	"time"
)

type timerFunc struct {
	Fn     func()
	Ticker int
}

var (
	timerCount int
	timerFuncs map[string]*timerFunc
)

func init() {
	timerCount = 0
	timerFuncs = make(map[string]*timerFunc)
}

// SetTimerFunc adds timer func for time ticker.
// Ticker means step time, after ticker size step passed, do function.
// Name is unique name of func.If set same name func, use the last one.
func SetTimerFunc(name string, ticker int, fn func()) {
	tfn := new(timerFunc)
	tfn.Fn = fn
	tfn.Ticker = ticker
	timerFuncs[name] = tfn
}

// ChangeTimerFunc can change timer func by given name.
// If the func of name is none, do not change anything, print error message.
func ChangeTimerFunc(name string, ticker int, fn func()) {
	if _, ok := timerFuncs[name]; ok {
		timerFuncs[name].Fn = fn
		timerFuncs[name].Ticker = ticker
	} else {
		println("change invalid timer func : " + name)
	}
}

// DelTimerFunc deletes timer func.
func DelTimerFunc(name string) {
	delete(timerFuncs, name)
}

// GetTimerFuncs returns registered timer func with its name and ticker int.
func GetTimerFuncs() map[string]int {
	m := make(map[string]int)
	for n, f := range timerFuncs {
		m[n] = f.Ticker
	}
	return m
}

// StartModelTimer adds models' timer and starts time ticker.
// The default step is 10 min once.
func StartModelTimer() {
	// start all timers
	startCommentsTimer()
	startContentSyncTimer()
	startContentTmpIndexesTimer()
	startMessageTimer()
	// start time ticker
	ticker := time.NewTicker(time.Duration(10) * time.Minute)
	go doTimers(ticker.C)
}

func doTimers(c <-chan time.Time) {
	for {
		<-c
		timerCount++
		for _, tfn := range timerFuncs {
			if timerCount%tfn.Ticker == 0 {
				tfn.Fn()
			}
		}
		if timerCount > 999 {
			timerCount = 0
		}
	}
}
