package internal

import (
    "log"
    "os"
)

var (
    LogFile *os.File
    Logger  *log.Logger
)

func InitLogger() {
    var err error
    LogFile, err = os.OpenFile("chat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("Failed to open log file:", err)
    }
    Logger = log.New(LogFile, "", log.LstdFlags)
}

func CloseLogger() {
    LogFile.Close()
}
