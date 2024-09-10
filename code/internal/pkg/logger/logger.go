package logger

import (
	"log"
	"os"
)

type AppLogger struct {
	infoLog  *log.Logger
	fatalLog *log.Logger
}

func (a *AppLogger) Info(info ...interface{}) {
	a.infoLog.Println(info...)
}

func (a *AppLogger) Fatal(fatal ...interface{}) {
	a.fatalLog.Println(fatal...)
}

func CreateNewLogger() *AppLogger {
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLog := log.New(os.Stdout, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
	return &AppLogger{infoLog: infoLog, fatalLog: fatalLog}
}
