package watchgod

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func (w *Watchgod) setupSingalHandlers() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Printf("\n")
		log.Printf("[FATAL] [watchgod] TERMINATE >>> due to signal '%s'\n", sig)
		w.eventChannel <- Event{eventType: TERMINATE, response: make(chan RPCResponse, 1)}
	}()
}

func (w *Watchgod) startResponseSink() chan RPCResponse {
	responseSink := make(chan RPCResponse, 1)
	go func() {
		for w.terminating == false {
			<-responseSink
		}
	}()
	return responseSink
}

func (w *Watchgod) startMonitor(responseSink chan RPCResponse) {
	go func() {
		for w.terminating == false {
			state := <-w.ripperChannel
			if state.state == DEAD && w.terminating == false {
				w.eventChannel <- Event{eventType: START, id: state.id, response: responseSink}
			}
		}
	}()
}

func (w *Watchgod) runEventLoop(responseSink chan RPCResponse) {
	for event := range w.eventChannel {
		switch event.eventType {
		case LIST:
			buffer := fmt.Sprintf(" %-4s  %-20s %-10s\n---------------------------------------\n", "PID", "NAME", "STATE")
			for _, process := range w.processes {
				buffer += fmt.Sprintf("[%4d] %-20s %-10s\n", process.pid, process.id, process.state)
			}
			event.response <- RPCResponse{err: nil, msg: buffer} // Do not use sendResponse() to avoid poluting the logs

		case TERMINATE:
			w.terminate(event)

		case CREATE:
			w.create(event)

		case START:
			w.start(event, responseSink)

		case STOP:
			w.stop(event)

		case RESTART:
			w.restart(event)

		default:
			Fatal("[watchgod] MainLoop unknown event >>> %v", event)
		}
	}
}

func (w *Watchgod) terminate(event Event) {
	w.terminating = true
	for _, process := range w.processes {
		if process.pid > 0 && process.state != DEAD {
			Kill(process.pid, syscall.SIGTERM)
		}
	}
	for _, process := range w.processes {
		if process.pid > 0 && process.state != DEAD {
			Wait(process.pid)
		}
	}
	close(w.eventChannel)
	log.Printf("[INFO] [watchgod] terminated\n")
	sendResponse(event.response, RPCResponse{err: nil, msg: "WatchGOd: terminated"})
}

func (w *Watchgod) create(event Event) {
	_, err := w.findById(event.id)
	if err != nil {
		process := newProcess(event.id, event.arguments, w.ripperChannel)
		w.processes = append(w.processes, &process)
		sendResponse(event.response, RPCResponse{err: nil, msg: fmt.Sprintf("%s: created", event.id)})
	} else {
		sendResponse(event.response, RPCResponse{err: errors.New(fmt.Sprintf("%s: already exist", event.id))})
	}
}

func (w *Watchgod) start(event Event, responseSink chan RPCResponse) {
	process, err := w.findById(event.id)
	if err != nil {
		sendResponse(event.response, RPCResponse{err: err})
		return
	}
	monitor := make(chan ProcessInfo, 1)

	go func(event Event, process *MonitoredProcess, monitor chan ProcessInfo) {
		processInfo := process.waitForNextEvent(monitor, w.startTimeoutInSeconds)
		if processInfo.state == RUNNING {
			processInfo = process.waitForNextEvent(monitor, w.startTimeoutInSeconds)
		}
		switch processInfo.state {
		case DEAD:
			w.eventChannel <- Event{eventType: STOP, id: processInfo.id, response: responseSink}
			sendResponse(event.response, RPCResponse{err: errors.New(fmt.Sprintf("%s: [DISABLED] PID %d exited at launch with code %d", event.id, processInfo.pid, processInfo.exitcode))})
		case TIMEOUT:
			sendResponse(event.response, RPCResponse{msg: fmt.Sprintf("%s: started with PID %d", event.id, processInfo.pid)})
		default:
			sendResponse(event.response, RPCResponse{err: processInfo.err})
		}
	}(event, process, monitor)
	process.start(monitor)
}

func (w *Watchgod) stop(event Event) {
	process, err := w.findById(event.id)
	if err != nil {
		sendResponse(event.response, RPCResponse{err: err})
		return
	}
	monitor := make(chan ProcessInfo, 1)
	process.interceptRipperChannel(monitor) // to avoid the default restart policy
	go func(event Event, process *MonitoredProcess, monitor chan ProcessInfo) {
		processInfo := process.waitForNextEvent(monitor, w.stopTimeoutInSeconds)
		switch processInfo.state {
		case DEAD:
			sendResponse(event.response, RPCResponse{msg: fmt.Sprintf("%s: stopped", event.id)})
		case TIMEOUT:
			sendResponse(event.response, RPCResponse{err: errors.New(fmt.Sprintf("%s: [TIMEOUT] is still running", event.id))})
		default:
			sendResponse(event.response, RPCResponse{err: processInfo.err})
		}
		process.releaseRipperChannel(monitor)
	}(event, process, monitor)
	process.stop(monitor)
}

func (w *Watchgod) restart(event Event) {
	process, err := w.findById(event.id)
	if err != nil {
		sendResponse(event.response, RPCResponse{err: err})
		return
	}
	monitor := make(chan ProcessInfo, 1)
	process.interceptRipperChannel(monitor) // to restart if it's already DEAD or when it dies instead of the default policy
	go func(event Event, process *MonitoredProcess, monitor chan ProcessInfo) {
		processInfo := process.waitForNextEvent(monitor, w.stopTimeoutInSeconds)
		process.releaseRipperChannel(monitor)
		switch processInfo.state {
		case DEAD:
			w.eventChannel <- Event{eventType: START, id: processInfo.id, response: event.response}
		case ALREADYDEAD:
			w.eventChannel <- Event{eventType: START, id: processInfo.id, response: event.response}
		case TIMEOUT:
			sendResponse(event.response, RPCResponse{err: errors.New(fmt.Sprintf("%s: [TIMEOUT] is still running", event.id))})
		default:
			sendResponse(event.response, RPCResponse{err: processInfo.err})
		}
	}(event, process, monitor)
	process.stop(monitor)
}

func sendResponse(output chan RPCResponse, response RPCResponse) {
	if response.err != nil {
		log.Printf("[ERROR] [watchgod] %s", response.err)
	} else {
		log.Printf("[INFO] [watchgod] %s", response.msg)
	}
	output <- response
}

func (w *Watchgod) findById(id string) (*MonitoredProcess, error) {
	for i := 0; i < len(w.processes); i++ {
		if w.processes[i].id == id {
			return w.processes[i], nil
		}
	}
	return nil, errors.New(fmt.Sprintf("%s: not found", id))
}
