package system

import "os"

// Ugly global exit handler!
var exit_impl = os.Exit

// Die with an error message.
func Exit(code int) {
	exit_impl(code)
}

/**
 *  Mainly for tests, replate the exit function with your
 * own implementation:
 *
 *    exited := false
 *    system.SetExitHandler(func (code int) {
 *      if code != 2 {
 *        t.Fatalf("expected exit code 2, got %d", code)
 *      }
 *      exited = true
 *    })
 *
 *    system.Exit(2)
 *
 *    if exited == false {
 *      t.Fatalf("Failed to exit")
 *    }
 *
 **/
func SetExitHandler(impl func(int)) {
	exit_impl = impl
}
