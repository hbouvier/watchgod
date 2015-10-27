package watchgod

type ProcessInfo struct {
	id             string
	state          ProcessState
	pid            int
	requestId      int
	exitcode       int
	pauseInSeconds int
	err            error
}
