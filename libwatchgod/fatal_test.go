package watchgod

import "testing"

func TestFatal(t *testing.T) {
	gExited := false

	SetExitHandler(func(code int) {
		if code != 2 {
			t.Fatalf("When calling Fatal, we expect the process to exit with a status code of two (2)."+
				" Process exited with status code of %d instead.", code)
		}
		gExited = true
	})

	Fatal("Invoking Fatal() to exit the process with an exit code of two (2).")

	if gExited == false {
		t.Fatalf("Expected Exit(2) to be called.")
	}
}
