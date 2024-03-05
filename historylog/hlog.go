package historylog

import (
	"log"
)

type HLog struct {
	*log.Logger
	File *HFile
}

func HLogWithFile(out *HFile) *HLog {
	return &HLog{
		Logger: log.New(out, "", log.LstdFlags),
		File:   out,
	}
}

func (l *HLog) TakeOutHistory() string {
	return l.File.TakeOutHistory()
}
