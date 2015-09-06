package system

import (
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "time"
)

type Process struct {
  commandTokens []string
  pid int
}

type Watchdog struct {
  terminating bool
  graveyard chan Process
  processes []Process
}

func (w * Watchdog) Initialize() {
  w.terminating = false
  w.graveyard = make(chan Process, 1)
  w.processes = make([]Process, 0)
  w.setupSingalHandlers()
}

func (w * Watchdog) Spawn(commandTokens []string) (int, error) {
  pid, spawnErr := Spawn(commandTokens)
  if spawnErr != nil {
    fmt.Fprintf(os.Stderr, timestamp() + "[watchdog] Spawn(%v) >>> %s\n", commandTokens, spawnErr)
    return -1, spawnErr
  }
  w.processes = append(w.processes, Process{commandTokens: commandTokens, pid: pid})

  go func() {
    status, waitErr := Wait(pid)
    if waitErr != nil {
      fmt.Fprintf(os.Stderr, 
                  timestamp() + "[watchdog] Spawn(%v) [%d] process was killed >>> %s\n",
                  commandTokens,
                  pid,
                  waitErr)
    } else {
      fmt.Fprintf(os.Stdout, timestamp() + "[watchdog] Spawn(%v) process exited with code %d\n", commandTokens, status)
    }
    w.graveyard <- Process{commandTokens: commandTokens, pid: pid}
  }()
  return pid, nil
}

func (w * Watchdog) WatchLoop() {
  fmt.Fprintf(os.Stdout, timestamp() + "[watchdog] WatchLoop\n")
  for w.terminating == false {
    process := <- w.graveyard
    w.terminateOnSignal(process)
    w.remove(process.pid)
    if w.terminating == false {
      time.Sleep(1000 * time.Millisecond)
      commandTokens := process.commandTokens
      _, err := w.Spawn(commandTokens)
      if err != nil {
        Fatal(timestamp() + "[watchdog] WatchLoop >>> %s", err)
      }
    }
  }
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (w * Watchdog) setupSingalHandlers() {
  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
  go func() {
    sig := <-sigs
    w.terminating = true
    fmt.Fprintf(os.Stderr, "\n" + timestamp() + "[watchdog] signalHandler() >>> %s\n", sig)
    w.graveyard <- Process{pid: 1}
  }()
}

func (w * Watchdog) remove(pid int) {
  var index = -1
  for i := 0; i < len(w.processes); i++ {
    if w.processes[i].pid == pid {
      index = i
      break
    }
  }
  if index == -1 {
    return
  }
  w.processes = append(w.processes[:index], w.processes[index+1:]...)
}

func (w * Watchdog) terminateOnSignal(process Process) {
  if process.pid != 1 {
    return
  }
  fmt.Fprintf(os.Stderr, timestamp() + "[watchdog] terminating all childrens due to SIGTERM\n")
  for _, process := range w.processes {
    Kill(process.pid, syscall.SIGTERM)
  }
}

