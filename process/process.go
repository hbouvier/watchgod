package process

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func init() {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

// Spawn ...
func Spawn(tokens []string) (int, error) {
	cmd, lookError := exec.LookPath(tokens[0])
	if lookError != nil {
		log.Printf("[ERROR] [process] Spawn.exec.LookPath(%s) >>> %s\n", tokens[0], lookError)
		return 0, lookError
	}

	var sysAttr syscall.SysProcAttr
	var procAttr syscall.ProcAttr
	procAttr.Sys = &sysAttr
	procAttr.Env = os.Environ()
	procAttr.Files = []uintptr{uintptr(syscall.Stdin),
		uintptr(syscall.Stdout),
		uintptr(syscall.Stderr)}

	pid, forkError := syscall.ForkExec(cmd, tokens, &procAttr)
	if forkError != nil {
		log.Printf("[ERROR] [process] Spwan.syscall.ForkExec(%s) >>> %s\n", cmd, forkError)
		return 0, forkError
	}
	return pid, nil
}

// Wait ...
func Wait(pid int) (int, error) {
	var wstat syscall.WaitStatus
	_, err := syscall.Wait4(pid, &wstat, 0, nil)
	if err != nil {
		if err.Error() != "no child processes" {
			log.Printf("[ERROR] [process] [%d] Wait >>> %s\n", pid, err)
			return -1, err
		}
		log.Printf("[DEBUG] [process] [%d] Wait >>> %s\n", pid, err)
		return -1, nil
	}

	status := wstat.ExitStatus()
	return status, nil
}

// Kill ...
func Kill(pid int, signal syscall.Signal) error {
	err := syscall.Kill(pid, signal)
	if err != nil {
		if err.Error() != "no such process" {
			log.Printf("[ERROR] [process] [%d] Kill(signal:%s) >>> %s\n", pid, signal, err)
			return err
		}
		log.Printf("[DEBUG] [process] [%d] Kill(signal:%s) >>> %s\n", pid, signal, err)
	}
	return nil
}
