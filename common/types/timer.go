package types

import (
	"github.com/DhunterAO/goAuthChain/log"
	"sync"
	"time"
)

type MyTimer struct {
	Name   string
	Delay  uint64
	Target func()
	Timer  *time.Timer

	logger *log.Logger

	Running   bool
	RunningMu sync.Mutex
}

func NewMyTimer(name string, delay uint64, target func(), logger *log.Logger) *MyTimer {
	newTimer := &MyTimer{
		Name:      name,
		Delay:     delay,
		Target:    target,
		Timer:     new(time.Timer),
		logger:    logger,
		Running:   false,
		RunningMu: sync.Mutex{},
	}
	return newTimer
}

// The target of timer will start in {delay} seconds, when delay is 0, it means that it will start in the default time.
func (mt *MyTimer) Start(delay uint64) {
	mt.RunningMu.Lock()
	defer func() {
		mt.RunningMu.Unlock()
		mt.Running = true
	}()

	if delay == uint64(0) {
		mt.logger.Info(mt.Name + "start, " + string(mt.Delay) + " seconds later call target")
		mt.Timer = time.AfterFunc(time.Duration(mt.Delay)*time.Second, mt.Target)
	} else {
		mt.logger.Info(mt.Name + "start, " + string(delay) + " seconds later call target")
		mt.Timer = time.AfterFunc(time.Duration(delay)*time.Second, mt.Target)
	}
}

// Stop() stops the timer and returns if it is running now, true for running and false for not running.
func (mt *MyTimer) Stop() bool {
	mt.logger.Info(mt.Name + " end.")

	mt.RunningMu.Lock()
	defer func() {
		mt.RunningMu.Unlock()
		mt.Running = false
	}()

	if mt.Running {
		return mt.Timer.Stop()
	} else {
		return false
	}
}

func (mt *MyTimer) Reset(delay uint64) {
	mt.Stop()
	mt.Start(delay)
}
