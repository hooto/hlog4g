// Copyright 2013 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hlog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hooto/hflag4g/hflag"
)

const (
	Version            = "0.9.2"
	printDefault uint8 = iota
	printFormat
)

var (
	locker    sync.Mutex
	levels    = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	levelChar = "DIWEF"
	levelMap  = map[string]int{}
	levelOut  = map[string]int{}
	bufs      = make(chan *entry, 100000)
	onceCmd   sync.Once
	log_ws    = map[int]*logFileWriter{}

	// If non-empty, write log files in this directory
	logDir = ""
	// log to standard error instead of files
	logToStderr = false
	// Messages logged at a lower level than this don't actually get logged anywhere
	minLogLevel = 1
	// Write log to multi level files
	logToLevels = false
)

type entry struct {
	ptype      uint8
	level      int
	format     string
	fileName   string
	lineNumber int
	ltime      time.Time
	args       []interface{}
}

func init() {

	if v, ok := hflag.ValueOK("log_dir"); ok {
		logDir = v.String()
	}

	if hflag.Value("logtostderr").String() == "true" {
		logToStderr = true
	}

	if v, ok := hflag.ValueOK("minloglevel"); ok {
		minLogLevel = v.Int()
	}

	if hflag.Value("logtolevels").String() == "true" {
		logToLevels = true
	}

	levelInit()
}

func LogDirSet(path string) {
	logDir = path
}

func LevelConfig(ls []string) {

	if len(ls) < 1 {
		return
	}

	levels = []string{}
	for _, v := range ls {
		levels = append(levels, strings.ToUpper(v))
	}

	levelInit()
}

func levelInit() {

	locker.Lock()
	defer locker.Unlock()

	//
	for _, wr := range log_ws {
		wr.Close()
	}
	log_ws = map[int]*logFileWriter{}

	//
	levelMap = map[string]int{}
	levelChar = ""
	for _, tag := range levels {

		if _, ok := levelMap[tag]; !ok {
			levelMap[tag] = len(levelMap)
			levelChar += tag[0:1]
		}
	}

	if minLogLevel < 0 {
		minLogLevel = 0
	} else if minLogLevel >= len(levelMap) {
		minLogLevel = len(levelMap) - 1
	}

	//
	for tag, level := range levelMap {

		//
		if (logToLevels == true && level >= minLogLevel) ||
			(logToLevels == false && level == minLogLevel) {

			levelOut[tag] = level
		}
	}

	onceCmd.Do(outputAction)
}

func (e *entry) line() string {

	logLine := fmt.Sprintf("%s %s %s:%d] ", string(levelChar[e.level]),
		e.ltime.Format("2006-01-02 15:04:05.000000"), e.fileName, e.lineNumber)

	if e.ptype == printDefault {
		logLine += fmt.Sprint(e.args...)
	} else if e.ptype == printFormat {
		logLine += fmt.Sprintf(e.format, e.args...)
	}

	return logLine + "\n"
}

func newEntry(ptype uint8, level_tag, format string, a ...interface{}) {

	level_tag = strings.ToUpper(level_tag)

	level, ok := levelMap[level_tag]
	if !ok || level < minLogLevel {
		return
	}

	// It's always the same number of frames to the user's call.
	_, fileName, lineNumber, ok := runtime.Caller(2)
	if !ok {
		fileName = "?"
		lineNumber = 1
	} else {
		slash := strings.LastIndex(fileName, "/")
		if slash >= 0 {
			fileName = fileName[slash+1:]
		}
	}
	if lineNumber < 0 {
		lineNumber = 0 // not a real line number, but acceptable to someDigits
	}

	bufs <- &entry{
		ptype:      ptype,
		level:      level,
		format:     format,
		fileName:   fileName,
		lineNumber: lineNumber,
		args:       a,
		ltime:      time.Now(),
	}
}

func Print(level string, a ...interface{}) {
	newEntry(printDefault, level, "", a...)
}

func Printf(level, format string, a ...interface{}) {
	newEntry(printFormat, level, format, a...)
}

func Flush() error {

	var e error

	for {

		if len(bufs) > 0 {
			time.Sleep(10e6)
			continue
		}

		for _, ws := range log_ws {
			if err := ws.Sync(); err != nil {
				if e == nil {
					e = err
				}
			}
		}

		time.Sleep(10e6)
		break
	}

	return e
}

func outputAction() {

	go func() {

		for logEntry := range bufs {

			bs := []byte(logEntry.line())

			if logToStderr {
				os.Stderr.Write(bs)
			}

			if len(logDir) > 0 {

				for level_tag, level := range levelOut {

					if logEntry.level < level {
						continue
					}

					locker.Lock()
					wr, _ := log_ws[level]

					if wr == nil {
						wr = newLogFileWriter(level_tag)
						log_ws[level] = wr
					}

					locker.Unlock()

					wr.Write(bs)
				}
			}
		}
	}()
}
