package service

import "log"

type MyLog struct {
	PrintDebug bool
}

func (m *MyLog) Debug(args ...interface{}) {
	if m.PrintDebug {
		m.Print(args...)
	}
}

func (m *MyLog) Print(args ...interface{}) {
	log.Print(args...)
}
