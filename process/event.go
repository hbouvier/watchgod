package process

import "fmt"

// EventType ...
type EventType int

const ( // iota is reset to 0
	// CREATE a new process without starting it
	CREATE EventType = 1 << iota

	// START  an already created process
	START
	// STOP  a running process
	STOP
	// RESTART a running or stopped process
	RESTART
	// LIST ...
	LIST
	// TERMINATE the Watch GO Deamon (and all processes)
	TERMINATE
)

// Event ...
type Event struct {
	id        string
	arguments []string
	eventType EventType
	exitcode  int
	requestID int
	response  chan RPCResponse
}

// String ...
func (event EventType) String() string {
	var s string
	switch event {
	case CREATE:
		s = "CREATE"
	case START:
		s = "START"
	case STOP:
		s = "STOP"
	case RESTART:
		s = "RESTART"
	case LIST:
		s = "LIST"
	case TERMINATE:
		s = "TERMINATE"
	default:
		s = fmt.Sprintf("[ERROR] Unknown event-type: %d", event)
	}
	return s
}
