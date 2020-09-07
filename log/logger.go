package log

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	l *log.Logger
}

var LstdFlags = log.LstdFlags

func New(out io.Writer, prefix string, flag int) *Logger {
	logger := log.New(out, prefix, flag)
	return &Logger{l: logger}
}

func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}

func (logger *Logger) Info(msg string) {
	logger.l.SetPrefix("INFO:")
	logger.l.Println(msg)
}

func (logger *Logger) Warn(msg string) {
	logger.l.SetPrefix("WARN:")
	logger.l.Println(msg)
}

func (logger *Logger) Error(msg string) {
	logger.l.SetPrefix("ERROR:")
	logger.l.Println(msg)
}

//func (logger *Logger)Close(msg string)  {
//	logger.l.SetPrefix("Closing:")
//	logger.l.Println(msg)
//	logger.l.p
//}

func NewTestLogger() *Logger {
	testLogFile, err := os.Create("test.log")
	if err != nil {
		Fatalln("open blockchain log file error !")
	}
	defer testLogFile.Close()
	testLogger := New(testLogFile, "[LOG]", log.LstdFlags)
	return testLogger
}

func LogError(err error, logger Logger) {
	if err != nil {
		logger.Error(err.Error())
	}
}
