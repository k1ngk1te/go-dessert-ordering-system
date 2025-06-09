package app

import (
	"log"
	"os"
)

type ApplicationLoggers struct {
	Error *log.Logger
	Info  *log.Logger
}

func NewApplicationLoggers() *ApplicationLoggers {
	return &ApplicationLoggers{
		Error: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime),
		Info:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}
}