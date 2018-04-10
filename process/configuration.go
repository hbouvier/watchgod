package process

import (
	"encoding/json"
	"log"
	"os"
)

// When starting a process Wait on it to see if it will exit int
// that time. If it does return an error, otherwise assume that it
// is stable.
var gStartTimeoutInSeconds = 1
var gStopTimeoutInSeconds = 5
var gIPCServerURL = "127.0.0.1:7099"
var gLogLevel = "INFO"

// Process ...
type Process struct {
	Name    string
	Command []string
}

// Configuration ...
type Configuration struct {
	LogLevel              string
	IPCServerURL          string
	StartTimeoutInSeconds int // Wait in seconds to see if process exit, before returning OK
	StopTimeoutInSeconds  int // Wait in seconds for process to die, before returning FAILURE
	Processes             []Process
}

// DefaultConfiguration ...
func DefaultConfiguration() Configuration {
	ipcURL := os.Getenv("IPC_SERVER_URL")
	if ipcURL == "" {
		ipcURL = gIPCServerURL
	}
	return Configuration{StartTimeoutInSeconds: gStartTimeoutInSeconds,
		StopTimeoutInSeconds: gStopTimeoutInSeconds,
		IPCServerURL:         ipcURL,
		LogLevel:             gLogLevel,
		Processes:            make([]Process, 0)}
}

// LoadConfiguration ...
func LoadConfiguration(filename string) Configuration {
	log.Printf("[INFO][watchgod] Loading %s\n", filename)
	file, fileErr := os.Open(filename)
	if fileErr != nil {
		log.Fatalf("%s: ERROR loading configuration file %s >>> %s\n", os.Args[0], filename, fileErr)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := DefaultConfiguration()
	decoderErr := decoder.Decode(&configuration)
	if decoderErr != nil {
		log.Fatalf("%s: ERROR decoding configuration file %s >>> %s\n", os.Args[0], filename, decoderErr)
	}
	return configuration
}
