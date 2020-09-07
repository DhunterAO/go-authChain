package types

import (
	"fmt"
	"github.com/DhunterAO/goAuthChain/log"
	"testing"
)

var timer *MyTimer

func echo() {
	fmt.Println("iiiiii")
	timer.Reset(1)
}

func TestNewMyTimer(t *testing.T) {
	timer = NewMyTimer("testTimer", 2, echo, log.NewTestLogger())
	//timer.Reset(1)
	//time.Sleep(10 * time.Second)
}
