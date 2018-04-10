package process

import (
	"fmt"
	"os"
)

// Die with an error message.
// func Fatal(msg string, args ...interface{}) {
// 	log.Printf("[FATAL] "+msg, args...)
// 	os.Exit(2)
// }

// FatalCli dies with an error message, but does not print the timestamp
func FatalCli(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(2)
}
