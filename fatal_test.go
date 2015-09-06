package system

import "testing"

func TestFatal(t *testing.T) {
	exited := false
	SetExitHandler(func(code int) {
		if code != 2 {
			t.Fatalf("expected exit code 2, got %d", code)
		}
		exited = true
	})

	Fatal("El Kaboum!")

	if exited == false {
		t.Fatalf("Failed to exit")
	}
}
