package process

import (
	"fmt"
	"log"
	"time"
)

// Watchgod ...
type Watchgod struct {
	version               string
	terminating           bool
	eventChannel          chan Event
	ripperChannel         chan ProcessInfo
	processes             []*MonitoredProcess
	startTimeoutInSeconds int
	stopTimeoutInSeconds  int
}

// Initialize ...
func (w *Watchgod) Initialize(version string, startTimeoutInSeconds int, stopTimeoutInSeconds int) {
	w.version = version
	w.terminating = false
	w.eventChannel = make(chan Event, 1)
	w.ripperChannel = make(chan ProcessInfo, 1)
	w.processes = make([]*MonitoredProcess, 0)
	w.setupSingalHandlers()
	w.startTimeoutInSeconds = startTimeoutInSeconds
	w.stopTimeoutInSeconds = stopTimeoutInSeconds
}

// MainLoop ...
func (w *Watchgod) MainLoop() {
	log.Printf("[INFO] [watchgod] Daemon version %s started...\n", w.version)
	responseSink := w.startResponseSink()
	w.startMonitor(responseSink)
	w.runEventLoop(responseSink)
	time.Sleep(time.Duration(1) * time.Second / 10) // Give 0.1 sec to the IPC message to be sent back to the client
}

// Add ...
func (w *Watchgod) Add(id string, arguments []string, response chan RPCResponse) {
	w.eventChannel <- Event{eventType: CREATE, id: id, arguments: arguments, response: response}
}

// Restart ...
func (w *Watchgod) Restart(id string, response chan RPCResponse) {
	w.eventChannel <- Event{eventType: RESTART, id: id, response: response}
}

// Start ...
func (w *Watchgod) Start(id string, response chan RPCResponse) {
	w.eventChannel <- Event{eventType: START, id: id, response: response}
}

// Stop ...
func (w *Watchgod) Stop(id string, response chan RPCResponse) {
	w.eventChannel <- Event{eventType: STOP, id: id, response: response}
}

// List ...
func (w *Watchgod) List(filter string, response chan RPCResponse) {
	w.eventChannel <- Event{eventType: LIST, id: "", response: response}
}

// Version ...
func (w *Watchgod) Version(response chan RPCResponse) {
	response <- RPCResponse{err: nil, msg: fmt.Sprintf("Deamon version %s", w.version)}
}

// Terminate ...
func (w *Watchgod) Terminate(reason string, response chan RPCResponse) {
	log.Printf("[DEBUG] [watchgod] Terminated by operator request: '%s'\n", reason)
	w.eventChannel <- Event{eventType: TERMINATE, response: response}
}
