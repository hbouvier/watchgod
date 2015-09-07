package watchgod

import (
	"fmt"
	"time"
)

func Timestamp() string {
	return fmt.Sprintf("%-42s ", time.Now())
}
