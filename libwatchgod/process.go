package watchgod

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func Spawn(tokens []string) (int, error) {
	cmd, lookError := exec.LookPath(tokens[0])
	if lookError != nil {
		log.Printf("[ERROR] [process] Spawn.exec.LookPath(%s) >>> %s\n", tokens[0], lookError)
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
	if forkError != nil {
		log.Printf("[ERROR] [process] Spwan.syscall.ForkExec(%s) >>> %s\n", cmd, forkError)
		return 0, forkError
	}
	return pid, nil
}

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
