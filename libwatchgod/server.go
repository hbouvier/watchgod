package watchgod

import (
	"os"
)

func Boot(url string, configuration Configuration) {
	watcher := new(Watchgod)
	watcher.Initialize(VERSION, configuration.StartTimeoutInSeconds, configuration.StopTimeoutInSeconds)
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
	}(watcher)
	watcher.MainLoop()
}

func IpcServerUrl(config_url string) string {
	var url string

	if env_url := os.Getenv("WATCHGOD_IPC_URL"); env_url != "" {
		url = env_url
	} else {
		url = config_url
	}
	return url
}
