package logger

import (
	"log"
	"os"
)

var (
	infoLog    *log.Logger
	errorLog   *log.Logger
)

func init() {
	infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(message string) {
	infoLog.Println(message)
}

func Error(message string) {
	errorLog.Println(message)
}