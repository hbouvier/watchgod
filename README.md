# Simple watchdog for docker containers is go

## How to use it

```bash
package main
import  (
  "os"
  "github.com/hbouvier/system"
)

func main() {
  watcher := system.Watchdog{}
  watcher.Initialize()
  watcher.Spawn(os.Args[1:])
	watcher.WatchLoop()
}
```
