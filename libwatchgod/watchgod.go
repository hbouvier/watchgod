package watchgod

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type Process struct {
	name          string
	commandTokens []string
	pid           int
	state         State
	exitcode      int
}

type Watchgod struct {
	terminating           bool
	eventQueue            chan Event
	processes             []Process
	mutex                 *sync.Mutex
	startTimeoutInSeconds int
}

func (w *Watchgod) Initialize(startTimeoutInSeconds int) {
	w.terminating = false
	w.eventQueue = make(chan Event, 1)
	w.processes = make([]Process, 0)
	w.setupSingalHandlers()
	w.startTimeoutInSeconds = startTimeoutInSeconds
	w.mutex = &sync.Mutex{}
}

func (w *Watchgod) Add(name string, commandTokens []string, response chan RPCResponse) error {
	w.mutex.Lock()
	if w.terminating == true {
		w.mutex.Unlock()
		return errors.New(fmt.Sprintf("add %s: watchgod is terminating...", name))
	}

	id := w.findProcessByName(name)
	if id >= 0 {
		w.mutex.Unlock()
		return errors.New(fmt.Sprintf("add %s: already exist", name))
	}
	w.processes = append(w.processes, Process{name: name, commandTokens: commandTokens, pid: 0, state: STOPPED})
	w.mutex.Unlock()
	response <- RPCResponse{err: nil, msg: fmt.Sprintf("Added %s", name)}
	return nil
}

func (w *Watchgod) Restart(name string, response chan RPCResponse) error {
	id := w.findProcessByName(name)
	if id < 0 {
		return errors.New(fmt.Sprintf("restart %s: does not exist", name))
	}
	w.eventQueue <- Event{eventType: RESTART, name: name, id: id, response: response}
	return nil
}

func (w *Watchgod) Start(name string, response chan RPCResponse) error {
	id := w.findProcessByName(name)
	if id < 0 {
		return errors.New(fmt.Sprintf("start %s: does not exist", name))
	}
	w.eventQueue <- Event{eventType: START, name: name, id: id, response: response}
	return nil
}

func (w *Watchgod) Stop(name string, response chan RPCResponse) error {
	id := w.findProcessByName(name)
	if id < 0 {
		return errors.New(fmt.Sprintf("stop %s: does not exist", name))
	}
	w.eventQueue <- Event{eventType: STOP, name: name, id: id, response: response}
	return nil
}

func (w *Watchgod) List(filter string, response chan RPCResponse) error {
	buffer := fmt.Sprintf(" %-4s  %-20s %-10s\n---------------------------------------\n", "PID", "NAME", "STATE")
	w.mutex.Lock()
	for _, process := range w.processes {
		buffer += fmt.Sprintf("[%4d] %-20s %-10s\n", process.pid, process.name, process.state)
	}
	w.mutex.Unlock()
	response <- RPCResponse{err: nil, msg: buffer}
	return nil
}

func (w *Watchgod) Terminate(reason string, response chan RPCResponse) error {
	fmt.Fprintf(os.Stderr, Timestamp()+"[watchgod] Terminated by operator request: '%s'\n", reason)
	w.eventQueue <- Event{eventType: TERMINATE, response: response}
	return nil
}

func (w *Watchgod) MainLoop() {
	fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] Daemon mainloop is running...\n")
	for event := range w.eventQueue {
		switch event.eventType {
		case TERMINATE:
			w.terminate()
			fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] terminated\n")
			event.response <- RPCResponse{err: nil, msg: "Terminated!"}

		case RESTART:
			if w.processes[event.id].state == RUNNING {
				w.processes[event.id].state = RESTARTING
				w.stop(event)
			} else {
				message := fmt.Sprintf("watchgod cannot %s %s [%d]: because it is %s",
					event.eventType, w.processes[event.id].name, w.processes[event.id].pid, w.processes[event.id].state)
				fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] %s\n", message)
				event.response <- RPCResponse{err: errors.New(message), msg: ""}
			}

		case RESTART_START:
			if w.processes[event.id].state == RESTARTING {
				w.start(event, w.startTimeoutInSeconds)
			} else {
				message := fmt.Sprintf("watchgod cannot %s %s [%d]: because it is %s",
					event.eventType, w.processes[event.id].name, w.processes[event.id].pid, w.processes[event.id].state)
				fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] %s\n", message)
				event.response <- RPCResponse{err: errors.New(message), msg: ""}
			}

		case START:
			if w.processes[event.id].state == DEAD || w.processes[event.id].state == STOPPED {
				w.start(event, w.startTimeoutInSeconds)
			} else {
				message := fmt.Sprintf("watchgod cannot %s %s [%d]: because it is %s",
					event.eventType, w.processes[event.id].name, w.processes[event.id].pid, w.processes[event.id].state)
				fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] %s\n", message)
				event.response <- RPCResponse{err: errors.New(message), msg: ""}
			}

		case STOP:
			if w.processes[event.id].state == RUNNING {
				w.stop(event)
			} else {
				message := fmt.Sprintf("watchgod cannot %s %s [%d]: because it is %s",
					event.eventType, w.processes[event.id].name, w.processes[event.id].pid, w.processes[event.id].state)
				fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] %s\n", message)
				event.response <- RPCResponse{err: errors.New(message), msg: ""}
			}

		case EXITED:
			if w.processes[event.id].state != STOPPED && w.processes[event.id].state != RESTARTING {
				fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] [%d] %-20s %s(%d) %s->DEAD\n",
					w.processes[event.id].pid, event.name, event.eventType, event.exitcode, w.processes[event.id].state)
				w.processes[event.id].pid = 0
				w.processes[event.id].state = DEAD
				go func() {
					time.Sleep(1000 * time.Millisecond)
					if w.terminating == false && w.processes[event.id].state != STOPPED {
						w.eventQueue <- Event{eventType: START, name: event.name, id: event.id, response: make(chan RPCResponse, 1)}
					}
				}()
			} else {
				message := fmt.Sprintf("[%d] %-20s %s(%d) RUNNING->%s",
					w.processes[event.id].pid, event.name, event.eventType, event.exitcode, w.processes[event.id].state)
				w.processes[event.id].pid = 0
				fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] %s\n", message)
				event.response <- RPCResponse{err: nil, msg: message}
			}

		default:
			Fatal(Timestamp()+"[watchgod] MainLoop unknown event >>> %v", event)
		}
	}
	fmt.Fprintf(os.Stdout, Timestamp()+"[watchgod] Daemon mainooop >>> EXITED\n")
}
