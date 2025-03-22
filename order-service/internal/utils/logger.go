package utils

import (
    "log"
    "os"
)

var (
    infoLogger    *log.Logger
    warningLogger *log.Logger
    errorLogger   *log.Logger
)

func init() {
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }

    infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    warningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(message string) {
    infoLogger.Println(message)
}

func Warning(message string) {
    warningLogger.Println(message)
}

func Error(message string) {
    errorLogger.Println(message)
}