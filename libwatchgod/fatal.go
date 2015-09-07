package watchgod

import (
	"fmt"
	"os"
	"time"
)

// Die with an error message.
func Fatal(msg string, args ...interface{}) {
	timestamp := fmt.Sprintf("%-42s ", time.Now())
	fmt.Fprintf(os.Stderr, timestamp+msg+"\n", args...)
	Exit(2)
}
