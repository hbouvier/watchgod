package watchgod

type ProcessInfo struct {
	id       string
	state    ProcessState
	pid      int
	exitcode int
	err      error
}
