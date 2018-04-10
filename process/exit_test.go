package process

import "testing"

func TestExit(t *testing.T) {
	gExited := false

	SetExitHandler(func(code int) {
		if code != 42 {
			t.Fatalf("Expect Exit(42) to exit with the value 42. Got %d instead.", code)
		}
		gExited = true
	})

	Exit(42)

	if gExited == false {
		t.Fatalf("Expected Exit(42) to be called.")
	}
}
