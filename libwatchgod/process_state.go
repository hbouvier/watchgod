package watchgod

import "fmt"

type ProcessState int

const ( // iota is reset to 0
	DEAD ProcessState = 1 << iota
	ALREADYDEAD
	RUNNING
	STOPPED
	RESTARTING
	TIMEOUT
	ERROR
)

func (state ProcessState) String() string {
	var s string
	switch state {
	case DEAD:
		s = "DEAD"
	case ALREADYDEAD:
		s = "ALREADYDEAD"
	case RUNNING:
		s = "RUNNING"
	case STOPPED:
		s = "STOPPED"
	case RESTARTING:
		s = "RESTARTING"
	case TIMEOUT:
		s = "TIMEOUT"
	case ERROR:
		s = "ERROR"
	default:
		s = fmt.Sprintf("[ERROR] Unknown ProcessState: %d", state)
	}
	return s
}
