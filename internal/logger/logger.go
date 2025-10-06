package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func Init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(msg string) {
	infoLogger.Println(msg)
}

func Error(msg string) {
	errorLogger.Println(msg)
}
