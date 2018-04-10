package process

import (
	"os"
	"strings"
)

// Boot ...
func Boot(url string, arguments []string, configuration Configuration, version string) {
	watcher := new(Watchgod)
	watcher.Initialize(version, configuration.StartTimeoutInSeconds, configuration.StopTimeoutInSeconds)
	StartIPCServer(url, watcher)
	go func(watcher *Watchgod) {
		response := make(chan RPCResponse, 1)
		go func() {
			for {
				<-response
			}
		}()
		for _, process := range configuration.Processes {
			watcher.Add(process.Name, process.Command, response)
			watcher.Start(process.Name, response)
		}
		if len(arguments) > 0 {
			pathToExecutable := strings.Split(arguments[0], string(os.PathSeparator))
			executableName := pathToExecutable[len(pathToExecutable)-1]
			watcher.Add(executableName, arguments, response)
			watcher.Start(executableName, response)
		}
	}(watcher)
	watcher.MainLoop()
}

// IpcServerURL ...
func IpcServerURL(configURL string) string {
	var url string

	if envURL := os.Getenv("WATCHGOD_IPC_URL"); envURL != "" {
		url = envURL
	} else {
		url = configURL
	}
	return url
}
