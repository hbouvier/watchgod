package watchgod

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"syscall"
	"time"
)

type MonitoredProcess struct {
	id             string
	arguments      []string
	pid            int
	state          ProcessState
	ripperChannels []chan ProcessInfo
	mutex          *sync.Mutex
}

func newProcess(id string, arguments []string, ripperChannel chan ProcessInfo) MonitoredProcess {
	ripperChannels := make([]chan ProcessInfo, 0)
	ripperChannels = append(ripperChannels, ripperChannel)
	return MonitoredProcess{id: id, arguments: arguments, pid: 0,
		state:          STOPPED,
		ripperChannels: ripperChannels,
		mutex:          &sync.Mutex{}}
}

func (p *MonitoredProcess) start(monitor chan ProcessInfo) {
	if p.state != DEAD && p.state != STOPPED {
		stateErr := errors.New(fmt.Sprintf("[ERROR] %-20s(%6d) process is not in DEAD state (%s)",
			p.id, p.pid, p.state))
		monitor <- ProcessInfo{id: p.id, state: ERROR, exitcode: 0, err: stateErr}
		return
	}

	pid, spawnErr := Spawn(p.arguments)
	if spawnErr != nil {
		monitor <- ProcessInfo{id: p.id, state: DEAD, exitcode: 0, err: spawnErr}
		return
	}
	p.pid = pid
	p.state = RUNNING
	monitor <- ProcessInfo{id: p.id, state: RUNNING, pid: pid, exitcode: 0, err: nil}

	go func() {
		exitcode, waitErr := Wait(pid)
		log.Printf("[INFO] [watchgod] %s: PID %d exit code %d", p.id, p.pid, exitcode)
		p.pid = 0
		p.state = DEAD
		processInfo := ProcessInfo{id: p.id, state: DEAD, pid: pid, exitcode: exitcode, err: waitErr}
		monitor <- processInfo
		p.sendToLastRipperChannel(processInfo)
	}()
}

func (p *MonitoredProcess) stop(monitor chan ProcessInfo) {
	if p.pid > 0 && p.state == RUNNING {
		Kill(p.pid, syscall.SIGTERM)
	} else {
		err := errors.New(
			fmt.Sprintf("[ERROR] %-20s(%6d) Process pid %d must be greater than zero and state '%s' must be RUNNING",
				p.id, p.pid, p.pid, p.state))
		monitor <- ProcessInfo{id: p.id, state: ALREADYDEAD, pid: p.pid, exitcode: 0, err: err}
	}
}

func (p *MonitoredProcess) waitForNextEvent(monitor chan ProcessInfo, timeoutInSeconds int) ProcessInfo {
	timeoutChannel := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(timeoutInSeconds) * time.Second)
		timeoutChannel <- true
	}()

	select {
	case processInfo := <-monitor:
		return processInfo
	case <-timeoutChannel:
		return ProcessInfo{id: p.id, state: TIMEOUT, pid: p.pid, exitcode: 0, err: nil}
	}
	return ProcessInfo{id: p.id, state: ERROR, exitcode: 0, pid: p.pid, err: errors.New("UNEXPECTED CASE")}
}

func (p *MonitoredProcess) interceptRipperChannel(ripperChannel chan ProcessInfo) {
	p.mutex.Lock()
	p.ripperChannels = append(p.ripperChannels, ripperChannel)
	p.mutex.Unlock()
}

func (p *MonitoredProcess) releaseRipperChannel(ripperChannel chan ProcessInfo) {
	p.mutex.Lock()
	for i := 0; i < len(p.ripperChannels); i++ {
		if p.ripperChannels[i] == ripperChannel {
			if i+1 < len(p.ripperChannels) {
				p.ripperChannels = append(p.ripperChannels[:i], p.ripperChannels[i+1:]...)
			} else {
				p.ripperChannels = p.ripperChannels[:i]
			}
			break
		}
	}
	p.mutex.Unlock()
}

func (p *MonitoredProcess) sendToLastRipperChannel(state ProcessInfo) {
	p.ripperChannels[len(p.ripperChannels)-1] <- state
}
