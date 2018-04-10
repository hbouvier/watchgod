package process

import "fmt"

// ProcessState ...
type ProcessState int

const ( // iota is reset to 0
	// DEAD ...
	DEAD ProcessState = 1 << iota
	// ALREADYDEAD ...
	ALREADYDEAD
	// RUNNING ...
	RUNNING
	// STOPPED ...
	STOPPED
	// TIMEOUT ...
	TIMEOUT
	// ERROR ...
	ERROR
)

// String ...
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
	case TIMEOUT:
		s = "TIMEOUT"
	case ERROR:
		s = "ERROR"
	default:
		s = fmt.Sprintf("[ERROR] Unknown ProcessState: %d", state)
	}
	return s
}
