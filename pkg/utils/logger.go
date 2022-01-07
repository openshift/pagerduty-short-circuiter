package utils

import (
	"io"
	"log"
)

var InfoLogger *log.Logger
var ErrorLogger *log.Logger

func InitLogger(logWriter io.Writer) {
	InfoLogger = log.New(logWriter, "[INFO]  ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(logWriter, "[ERROR] ", log.Ldate|log.Ltime)
}
