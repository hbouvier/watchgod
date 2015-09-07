package main

import (
	"encoding/json"
	"fmt"
	"github.com/hbouvier/watchgod/libwatchgod"
	"net/rpc"
	"os"
)

// When starting a process Wait on it to see if it will exit int
// that time. If it does return an error, otherwise assume that it
// is stable.
var gStartTimeoutInSeconds int = 2

type Process struct {
	Name    string
	Command []string
}

type Configuration struct {
	Processes []Process
}

func main() {
	if len(os.Args) < 2 {
		usage()
	} else if len(os.Args) == 2 {
		if os.Args[1] == "boot" {
			boot(ipcServerUrl(), []Process{})
		} else if os.Args[1] == "list" {
			list(ipcServerUrl())
		} else if os.Args[1] == "terminate" {
			terminate(ipcServerUrl())
		} else {
			usage()
		}
	} else if len(os.Args) == 3 {
		if os.Args[1] == "config" {
			configuration := loadConfiguration(os.Args[2])
			fmt.Printf("config[%s]: %v\n", os.Args[2], configuration)
			boot(ipcServerUrl(), configuration.Processes)
		} else if os.Args[1] == "start" {
			start(ipcServerUrl(), os.Args[2])
		} else if os.Args[1] == "stop" {
			stop(ipcServerUrl(), os.Args[2])
		} else if os.Args[1] == "restart" {
			restart(ipcServerUrl(), os.Args[2])
		} else {
			usage()
		}
	} else if len(os.Args) > 3 && os.Args[1] == "add" {
		add(ipcServerUrl(), os.Args[2:])
	} else {
		usage()
	}
}

func usage() {
	watchgod.Fatal("USAGE %s boot|quit|add|start|stop", os.Args[0])
}

func boot(url string, processes []Process) {
	watcher := new(watchgod.Watchgod)
	watcher.Initialize(gStartTimeoutInSeconds)
	watchgod.StartIPCServer(url, watcher)
	go func() {
		response := make(chan watchgod.RPCResponse, 1)
		go func() {
			for {
				<-response
			}
		}()
		for _, process := range processes {
			watcher.Add(process.Name, process.Command, response)
			watcher.Start(process.Name, response)
		}
	}()
	watcher.MainLoop()
}

func loadConfiguration(filename string) Configuration {
	file, fileErr := os.Open(filename)
	if fileErr != nil {
		watchgod.Fatal("%s: ERROR loading configuration file %s >>> %s\n", os.Args[0], filename, fileErr)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	decoderErr := decoder.Decode(&configuration)
	if decoderErr != nil {
		watchgod.Fatal("%s: ERROR decoding configuration file %s >>> %s\n", os.Args[0], filename, decoderErr)
	}
	return configuration
}

func list(url string) {
	fmt.Println(ipcInvoke(url, "List", "*"))
}

func terminate(url string) {
	fmt.Println(ipcInvoke(url, "Terminate", "now"))
}

func start(url string, name string) {
	fmt.Println(ipcInvoke(url, "Start", name))
}

func stop(url string, name string) {
	fmt.Println(ipcInvoke(url, "Stop", name))
}

func restart(url string, name string) {
	fmt.Println(ipcInvoke(url, "Restart", name))
}

func add(url string, args []string) {
	var reply string
	err := client(url).Call("IPCServer.Add", args, &reply)
	if err != nil {
		watchgod.Fatal("Error: %s", err)
	}
	fmt.Printf("%s\n", reply)
}

func ipcInvoke(url string, method string, argument string) string {
	var reply string
	err := client(url).Call("IPCServer."+method, argument, &reply)
	if err != nil {
		watchgod.Fatal("Error: %s", err)
	}
	return fmt.Sprintf("%s", reply)
}

func client(url string) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", url)
	if err != nil {
		watchgod.Fatal("Error connecting: %s", err)
	}
	return client
}

func ipcServerUrl() string {
	url := os.Getenv("watchgod_IPC_URL")
	if url == "" {
		url = "127.0.0.1:7099"
	}
	return url
}
