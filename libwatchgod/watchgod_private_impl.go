package watchgod

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (w *Watchgod) setupSingalHandlers() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Fprintf(os.Stderr, "\n"+Timestamp()+"[watchgod] TERMINATE >>> due to signal '%s'\n", sig)
		w.eventQueue <- Event{eventType: TERMINATE, response: make(chan RPCResponse, 1)}
	}()
}

func (w *Watchgod) spawn(id int, process Process, stable *chan error) (int, error) {
	pid, spawnErr := Spawn(process.commandTokens)
	if spawnErr != nil {
		fmt.Fprintf(os.Stderr, Timestamp()+"[watchgod] [0] %-20s ERROR spawning process (%v) >>> %s\n",
			process.name, process.commandTokens, spawnErr)
		return -1, spawnErr
	}

	go func() {
		exitcode, waitErr := Wait(pid)
		if stable != nil {
			*stable <- errors.New(fmt.Sprintf("%s [%d]: exited with %d", process.name, pid, exitcode))
		}
		if waitErr != nil {
			fmt.Fprintf(os.Stderr,
				Timestamp()+"[watchgod] [%d] %-20s ERROR monitoring process >>> %s\n", pid, process.name, waitErr)
		}
		if w.terminating == false {
			w.eventQueue <- Event{id: id, name: process.name, eventType: EXITED, exitcode: exitcode,
				response: make(chan RPCResponse, 1)}
		} else {
			fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] [%d] %-20s TERMINATING(%d) %s->DEAD\n",
				pid, process.name, exitcode, process.state)

		}
	}()
	return pid, nil
}

func (w *Watchgod) terminate() {
	w.mutex.Lock()
	w.terminating = true
	for _, process := range w.processes {
		if process.pid > 0 && process.state != DEAD {
			Kill(process.pid, syscall.SIGTERM)
		}
	}
	w.mutex.Unlock()
	for _, process := range w.processes {
		if process.pid > 0 && process.state != DEAD {
			Wait(process.pid)
		}
	}
	close(w.eventQueue)
}

func (w *Watchgod) start(event Event, timeoutInSeconds int) {
	var message string
	stable := make(chan error, 1)
	pid, err := w.spawn(event.id, w.processes[event.id], &stable)
	if err != nil {
		message = fmt.Sprintf("[%d] %-20s %s %s >>> %s",
			pid, event.name, event.eventType, w.processes[event.id].state, err)
		fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] %s\n", message)
		event.response <- RPCResponse{err: errors.New(message), msg: ""}
	} else {
		fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] [%d] %-20s %s %s->...staring...\n",
			pid, event.name, event.eventType, w.processes[event.id].state)
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(time.Duration(timeoutInSeconds) * time.Second)
			timeout <- true
		}()
		select {
		case e := <-stable:
			message = fmt.Sprintf("[%d] %-20s %s %s >>> %s",
				pid, event.name, event.eventType, w.processes[event.id].state, e)
			fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] %s\n", message)
			event.response <- RPCResponse{err: errors.New(message), msg: ""}
		case <-timeout:
			message = fmt.Sprintf("[%d] %-20s %s ...starting...->RUNNING",
				pid, event.name, event.eventType)
			w.processes[event.id].pid = pid
			w.processes[event.id].state = RUNNING
			fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] %s\n", message)
			event.response <- RPCResponse{err: nil, msg: message}
		}
	}
}

func (w *Watchgod) stop(event Event) {
	pid := w.processes[event.id].pid
	oldstate := RUNNING
	var newstate State
	if event.eventType == STOP {
		newstate = STOPPED
	} else if event.eventType == RESTART {
		newstate = RESTARTING
	} else {
		Fatal("Unknow event type %s (expecting STOP or RESTART)", event.eventType)
	}
	message := fmt.Sprintf("[%d] %-20s %s %s->%s",
		pid, event.name, event.eventType, oldstate, newstate)
	w.processes[event.id].state = newstate
	if pid > 0 {
		Kill(pid, syscall.SIGTERM)
		go func(event Event, pid int, message string) {
			Wait(pid)
			time.Sleep(1000 * time.Millisecond)
			if w.terminating == false {
				switch w.processes[event.id].state {
				case RESTARTING:
					w.eventQueue <- Event{eventType: RESTART_START, name: event.name, id: event.id, response: event.response}
				case STOPPED:
					message := fmt.Sprintf("[%d] %-20s %s DEAD->STOPPED",
						w.processes[event.id].pid, event.name, event.eventType)
					fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] %s\n", message)
					event.response <- RPCResponse{err: nil, msg: message}
				default:
					Fatal("Unknow state: %s", w.processes[event.id].state)
				}
			}
		}(event, pid, message)
	} else {
		event.response <- RPCResponse{err: errors.New(message), msg: ""}
	}
	fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] %s\n", message)
}

func (w *Watchgod) findProcessByName(name string) int {
	for i := 0; i < len(w.processes); i++ {
		if w.processes[i].name == name {
			return i
		}
	}
	return -1
}
