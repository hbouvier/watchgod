package process

import (
	"testing"
)

func TestCreateProcess(t *testing.T) {
	monitor := make(chan ProcessInfo, 1)
	process := newProcess("greetings", []string{"echo", "hello world"}, monitor)
	if process.id != "greetings" || process.pid != 0 || process.state != STOPPED {
		t.Fatalf("Create process did not retrun a valid process")
	}
}

func TestStartProcess(t *testing.T) {
	defaultMonitor := make(chan ProcessInfo, 1)
	go func() {
		for {
			<-defaultMonitor
		}
	}()
	process := newProcess("sleeper", []string{"sleep", "5"}, defaultMonitor)
	monitor := make(chan ProcessInfo, 1)
	process.start(monitor)
	processInfo := <-monitor

	if processInfo.State != RUNNING {
		t.Fatalf("Start process state is not RUNNING >>> %s", processInfo.State)
	}
	if processInfo.pid <= 0 {
		t.Fatalf("Start process PID is not greater than zero >>> %d", processInfo.pid)
	}
}

func TestStartAndDieProcess(t *testing.T) {
	defaultMonitor := make(chan ProcessInfo, 1)
	process := newProcess("greetings", []string{"echo", "hello world"}, defaultMonitor)
	monitor := make(chan ProcessInfo, 1)
	process.start(monitor)

	processInfo := <-monitor
	if processInfo.State != RUNNING {
		t.Fatalf("Start process state is not RUNNING >>> %s", processInfo.State)
	}

	processInfo = process.waitForNextEvent(defaultMonitor, 1)
	if processInfo.State != DEAD {
		t.Fatalf("Process did not died within one second (%s)", processInfo.State)
	}
}

func TestStartAndStableProcess(t *testing.T) {
	defaultMonitor := make(chan ProcessInfo, 1)
	process := newProcess("sleeper", []string{"sleep", "3"}, defaultMonitor)
	monitor := make(chan ProcessInfo, 1)
	process.start(monitor)

	processInfo := <-monitor
	if processInfo.State != RUNNING {
		t.Fatalf("Start process state is not RUNNING >>> %s", processInfo.State)
	}

	processInfo = process.waitForNextEvent(defaultMonitor, 2)
	if processInfo.State != TIMEOUT {
		t.Fatalf("Process probably died (%s), expected TIMEOUT", processInfo.State)
	}
}

func TestStartAndKillProcess(t *testing.T) {
	defaultMonitor := make(chan ProcessInfo, 1)
	process := newProcess("sleeper", []string{"sleep", "60"}, defaultMonitor)
	monitor := make(chan ProcessInfo, 1)
	process.start(monitor)

	processInfo := <-monitor
	if processInfo.State != RUNNING {
		t.Fatalf("Start process state is not RUNNING >>> %s", processInfo.State)
	}
	process.stop(monitor)
	processInfo = <-defaultMonitor
	if processInfo.State != DEAD {
		t.Fatalf("Process probably died (%s), expected TIMEOUT", processInfo.State)
	}
}
