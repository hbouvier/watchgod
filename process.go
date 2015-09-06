package system

import (
	"fmt"
  "time"
	"os"
	"os/exec"
	"syscall"
)

func Spawn(tokens []string) (int, error) {
  fmt.Fprintf(os.Stdout, timestamp() + "[process] Spawn(%s)\n", tokens[0])
	cmd, lookError := exec.LookPath(tokens[0])
	if lookError != nil {
		fmt.Fprintf(os.Stderr, timestamp() + "[process] Spawn.exec.LookPath(%s) >>> %s\n", tokens[0], lookError)
		return 0, lookError
	}

	var sys_attr syscall.SysProcAttr
	var proc_attr syscall.ProcAttr
	proc_attr.Sys = &sys_attr
	proc_attr.Env = os.Environ()
	proc_attr.Files = []uintptr{uintptr(syscall.Stdin),
		uintptr(syscall.Stdout),
		uintptr(syscall.Stderr)}

	pid, forkError := syscall.ForkExec(cmd, tokens, &proc_attr)
  fmt.Fprintf(os.Stdout, timestamp() + "[process] Spawn(%s) [%d]\n", tokens[0], pid)
	if forkError != nil {
		fmt.Fprintf(os.Stderr, timestamp() + "[process] Spwan.syscall.ForkExec(%s) >>> %s\n", cmd, forkError)
		return 0, forkError
	}
	fmt.Fprintf(os.Stdout, timestamp() + "[process] Spawn(%v) >>> PID: %d\n", tokens, pid)
	return pid, nil
}

func Wait(pid int) (int, error) {
  fmt.Fprintf(os.Stdout, timestamp() + "[process] Wait.Wait4(%d)\n", pid)
	var wstat syscall.WaitStatus
	_, err := syscall.Wait4(pid, &wstat, 0, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, timestamp() + "[process] Wait.Wait4(%d) >>> failed %s\n", pid, err)
		return -1 << 63, err
	}

	status := wstat.ExitStatus()
	fmt.Fprintf(os.Stdout, timestamp() + "[process] Wait.Wait4(%d) >>> exit-code: %d\n", pid, status)
	return status, nil
}

func Kill(pid int, signal syscall.Signal) error {
  proc, err := sendSignal(pid, signal)
  if err != nil {
    fmt.Fprintf(os.Stderr, timestamp() + "[process] Kill(pid:%d, signal:%s) >>> %s\n", pid, signal, err)
    return err
  }

  go func () {
    _, err := proc.Wait()
    fmt.Fprintf(os.Stderr, timestamp() + "[process] Kill(pid:%d, signal:%s) terminated >>> %s\n", pid, signal, err)
  }()
  return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func sendSignal(pid int, signal syscall.Signal) (*os.Process, error) {
  proc, err := os.FindProcess(pid)
  if err != nil {
    fmt.Fprintf(os.Stderr, 
                timestamp() + "[process] Signal(pid:%d, signal:%s) FindProcess(%d) >>> %s\n",
                pid, signal, pid, err)
    return proc, err
  }
  fmt.Fprintf(os.Stderr, timestamp() + "[process] Signal(pid:%d, signal:%s)\n", pid, signal)
  proc.Signal(signal)
  return proc, nil
}

func timestamp() string {
  return fmt.Sprintf("%-42s ", time.Now())
}