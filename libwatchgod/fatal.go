package watchgod

import (
	"fmt"
	"log"
)

// Die with an error message.
func Fatal(msg string, args ...interface{}) {
	log.Printf("[FATAL] "+msg, args...)
	Exit(2)
}

// Cli dies with an error message, but does not print the timestamp
func FatalCli(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	Exit(2)
}
