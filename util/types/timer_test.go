package types

import (
	"fmt"
	"testing"
	"time"
)

var timer *MyTimer

func echo() {
	fmt.Println("iiiiii")
	timer.Reset(1)
}

func TestNewMyTimer(t *testing.T) {
	timer = NewMyTimer("testTimer", 2, echo)
	timer.Reset(1)
	time.Sleep(10 * time.Second)
}
