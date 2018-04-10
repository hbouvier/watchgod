package process

import "os"

// Ugly global exit handler!
var exitImpl = os.Exit

// Exit die with an error message.
func Exit(code int) {
	exitImpl(code)
}

// SetExitHandler Mainly for tests, replate the exit function with your
// own implementation:
//
//    exited := false
//    system.SetExitHandler(func (code int) {
//      if code != 2 {
//        t.Fatalf("expected exit code 2, got %d", code)
//      }
//      exited = true
//    })
//
//    system.Exit(2)
//
//    if exited == false {
//      t.Fatalf("Failed to exit")
//    }
//
//
func SetExitHandler(impl func(int)) {
	exitImpl = impl
}
