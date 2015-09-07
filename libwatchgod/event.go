package watchgod

import "fmt"

type EventType int

const ( // iota is reset to 0
	START EventType = 1 << iota
	STOP
	RESTART
	RESTART_START
	TERMINATE
	EXITED
)

type Event struct {
	id        int
	name      string
	eventType EventType
	exitcode  int
	response  chan RPCResponse
}

func (event EventType) String() string {
	var s string
	switch event {
	case START:
		s = "START"
	case STOP:
		s = "STOP"
	case RESTART:
		s = "RESTART"
	case RESTART_START:
		s = "RESTART_START"
	case TERMINATE:
		s = "TERMINATE"
	case EXITED:
		s = "EXITED"
	default:
		s = fmt.Sprintf("[ERROR] Unknown event-type: %d", event)
	}
	return s
}
