package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/logutils"
	"github.com/hbouvier/watchgod/libwatchgod"
	"log"
	"net/rpc"
	"os"
)

var VERSION = "1.0.1"

// When starting a process Wait on it to see if it will exit int
// that time. If it does return an error, otherwise assume that it
// is stable.
var gStartTimeoutInSeconds int = 1
var gStopTimeoutInSeconds int = 5
var gIPCServerURL string = "127.0.0.1:7099"
var gLogLevel string = "INFO"

type Process struct {
	Name    string
	Command []string
}

type Configuration struct {
	LogLevel              string
	IPCServerURL          string
	StartTimeoutInSeconds int // Wait in seconds to see if process exit, before returning OK
	StopTimeoutInSeconds  int // Wait in seconds for process to die, before returning FAILURE
	Processes             []Process
}

func main() {
	setLogger(defaultConfiguration().LogLevel)
	if len(os.Args) < 2 {
		usage()
	} else if len(os.Args) == 2 {
		if os.Args[1] == "boot" {
			configuration := defaultConfiguration()
			boot(ipcServerUrl(configuration.IPCServerURL), configuration)
		} else if os.Args[1] == "list" {
			list(ipcServerUrl(defaultConfiguration().IPCServerURL))
		} else if os.Args[1] == "terminate" {
			terminate(ipcServerUrl(defaultConfiguration().IPCServerURL))
		} else if os.Args[1] == "version" {
			fmt.Printf("Client version %s\n", VERSION)
			version(ipcServerUrl(defaultConfiguration().IPCServerURL))
		} else {
			usage()
		}
	} else if len(os.Args) == 3 {
		if os.Args[1] == "boot" {
			configuration := loadConfiguration(os.Args[2])
			setLogger(configuration.LogLevel)
			boot(ipcServerUrl(configuration.IPCServerURL), configuration)
		} else if os.Args[1] == "start" {
			start(ipcServerUrl(defaultConfiguration().IPCServerURL), os.Args[2])
		} else if os.Args[1] == "stop" {
			stop(ipcServerUrl(defaultConfiguration().IPCServerURL), os.Args[2])
		} else if os.Args[1] == "restart" {
			restart(ipcServerUrl(defaultConfiguration().IPCServerURL), os.Args[2])
		} else {
			usage()
		}
	} else if len(os.Args) > 3 && os.Args[1] == "add" {
		add(ipcServerUrl(defaultConfiguration().IPCServerURL), os.Args[2:])
	} else {
		usage()
	}
}

func usage() {
	watchgod.FatalCli("USAGE %s list|start {name}|stop {name}|restart {name}|boot {config.json}|add {id} {args...}|terminate|version", os.Args[0])
}

func setLogger(level string) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel(level),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}

func boot(url string, configuration Configuration) {
	watcher := new(watchgod.Watchgod)
	watcher.Initialize(VERSION, configuration.StartTimeoutInSeconds, configuration.StopTimeoutInSeconds)
	watchgod.StartIPCServer(url, watcher)
	go func(watcher *watchgod.Watchgod) {
		response := make(chan watchgod.RPCResponse, 1)
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

func defaultConfiguration() Configuration {
	return Configuration{StartTimeoutInSeconds: gStartTimeoutInSeconds,
		StopTimeoutInSeconds: gStopTimeoutInSeconds,
		IPCServerURL:         gIPCServerURL,
		LogLevel:             gLogLevel,
		Processes:            make([]Process, 0)}
}

func loadConfiguration(filename string) Configuration {
	file, fileErr := os.Open(filename)
	if fileErr != nil {
		watchgod.Fatal("%s: ERROR loading configuration file %s >>> %s\n", os.Args[0], filename, fileErr)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := defaultConfiguration()
	decoderErr := decoder.Decode(&configuration)
	if decoderErr != nil {
		watchgod.Fatal("%s: ERROR decoding configuration file %s >>> %s\n", os.Args[0], filename, decoderErr)
	}
	return configuration
}

func version(url string) {
	fmt.Println(ipcInvoke(url, "Version", "*"))
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
		watchgod.FatalCli("Error: %s", err)
	}
	fmt.Printf("%s\n", reply)
}

func ipcInvoke(url string, method string, argument string) string {
	var reply string
	err := client(url).Call("IPCServer."+method, argument, &reply)
	if err != nil {
		watchgod.FatalCli("Error: %s", err)
	}
	return fmt.Sprintf("%s", reply)
}

func client(url string) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", url)
	if err != nil {
		watchgod.FatalCli("Error connecting: %s", err)
	}
	return client
}

func ipcServerUrl(config_url string) string {
	var url string

	if env_url := os.Getenv("WATCHGOD_IPC_URL"); env_url != "" {
		url = env_url
	} else {
		url = config_url
	}
	return url
}
