package watchgod

import "fmt"

type EventType int

const ( // iota is reset to 0
	CREATE  EventType = 1 << iota // New process without starting it
	START                         // Start an already created process
	STOP                          // Stop a running process
	RESTART                       // Restart a running or stopped process
	LIST
	TERMINATE // Terminate the Watch GO Deamon (and all processes)
)

type Event struct {
	id        string
	arguments []string
	eventType EventType
	exitcode  int
	response  chan RPCResponse
}

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
