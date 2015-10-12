package watchgod

import (
	"encoding/json"
	"log"
	"os"
)

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

func DefaultConfiguration() Configuration {
	return Configuration{StartTimeoutInSeconds: gStartTimeoutInSeconds,
		StopTimeoutInSeconds: gStopTimeoutInSeconds,
		IPCServerURL:         gIPCServerURL,
		LogLevel:             gLogLevel,
		Processes:            make([]Process, 0)}
}

func LoadConfiguration(filename string) Configuration {
	log.Printf("[INFO][watchgod] Loading %s\n", filename)
	file, fileErr := os.Open(filename)
	if fileErr != nil {
		Fatal("%s: ERROR loading configuration file %s >>> %s\n", os.Args[0], filename, fileErr)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := DefaultConfiguration()
	decoderErr := decoder.Decode(&configuration)
	if decoderErr != nil {
		Fatal("%s: ERROR decoding configuration file %s >>> %s\n", os.Args[0], filename, decoderErr)
	}
	return configuration
}
