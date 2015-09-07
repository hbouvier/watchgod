package watchgod

import "fmt"

type State int

const ( // iota is reset to 0
	DEAD State = 1 << iota
	RUNNING
	STOPPED
	RESTARTING
)

func (state State) String() string {
	var s string
	switch state {
	case DEAD:
		s = "DEAD"
	case RUNNING:
		s = "RUNNING"
	case STOPPED:
		s = "STOPPED"
	case RESTARTING:
		s = "RESTARTING"
	default:
		s = fmt.Sprintf("[ERROR] Unknown state: %d", state)
	}
	return s
}
