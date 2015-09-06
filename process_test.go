package system

import (
  "testing"
  "syscall"
)

func TestProcessSpawnFailure(t *testing.T) {
	_, err := Spawn([]string{"/does_not_exist"})
	if err == nil {
		t.Fatalf("Sawn failed to report an error when the process does not exist.")
	}
}

func TestProcessSpawnAndWait(t *testing.T) {
	args := []string{"/bin/ls", "-l"}
	pid, err := Spawn(args)
	if err != nil {
		t.Fatalf("Spawn reported an unexpect error: %s", err)
	}
	if pid <= 1 {
		t.Fatalf("Spawn returned an invalid pid (must be greater than 1): pid=%d", pid)
	}

	status, err := Wait(pid)
	if err != nil {
		t.Fatalf("Wait reported an unexpect error: %s", err)
	}
	if status != 0 {
		t.Fatalf("Spawn expect exit code to be zeros: status=%d", status)
	}
}

func TestProcessKill(t *testing.T) {
  args := []string{"sleep", "10"}
  pid, err := Spawn(args)
  if err != nil {
    t.Fatalf("Spawn reported an unexpect error: %s", err)
  }
  if pid <= 1 {
    t.Fatalf("Spawn returned an invalid pid (must be greater than 1): pid=%d", pid)
  }

  Kill(pid, syscall.SIGTERM)

  status, err := Wait(pid)
  if err != nil {
    t.Fatalf("Wait reported an unexpect error: %s", err)
  }
  if status != -1 {
    t.Fatalf("Spawn expect exit code to be zeros: status=%d", status)
  }
}
