package logger

import "log"

type Logger struct{}

func (l Logger) Debug(v ...interface{}) {
	log.Println(v...)
}

func (l Logger) Error(v ...interface{}) {
	log.Println(v...)
}
