package process

// ProcessInfo ...
type ProcessInfo struct {
	ID             string
	State          ProcessState
	pid            int
	requestID      int
	exitcode       int
	pauseInSeconds int
	err            error
}
