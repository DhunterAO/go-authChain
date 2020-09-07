package log

import (
	"log"
	"os"
	"testing"
)

func TestLogger_Info(t *testing.T) {
	testLogFile, err := os.Create("test.log")
	if err != nil {
		Fatalln("open file error !")
	}
	defer testLogFile.Close()
	testLogger := New(testLogFile, "[LOG]", log.LstdFlags)
	testLogger.Info("hello")
	testLogger.Warn("warning")
	testLogger.Error("error happens")
}
