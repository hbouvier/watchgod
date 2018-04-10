package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/logutils"
	"github.com/hbouvier/watchgod/process"
)

// version is set at build time using the go `-ldflags` like this:
// 		-ldflags "-X main.version=0.0.0"
var version = "snapshot"

func main() {
	configuration := process.DefaultConfiguration()

	levelPtr := flag.String("level", "", "logging level")
	configurationFilenamePtr := flag.String("config", "", "configuration filename.json")
	commandFlags := process.CommandFlags{
		ListPtr:      flag.Bool("list", false, "list the watched processes"),
		TerminatePtr: flag.Bool("terminate", false, "watchgod will terminate all monitored processes and exit"),
		VersionPtr:   flag.Bool("version", false, "Print the version of the watchgod command"),
		StartPtr:     flag.String("start", "", "Start a watch process that is in STOP state"),
		StopPtr:      flag.String("stop", "", "Stop a watch process"),
		RestartPtr:   flag.String("restart", "", "Restart a watch process"),
		UserPtr:      flag.String("user", "", "User to impersonate"),
		AddPtr:       flag.String("add", "", "Add a new process to watch"),
	}
	flag.Parse()

	if *configurationFilenamePtr != "" {
		configuration = process.LoadConfiguration(*configurationFilenamePtr)
	}
	if *levelPtr != "" {
		configuration.LogLevel = *levelPtr
	}
	setLogger(configuration.LogLevel)

	process.ExecuteArgument(commandFlags, flag.Args(), configuration, version, usage)
}

func usage() {
	pathToExecutable := strings.Split(os.Args[0], string(os.PathSeparator))
	executableName := pathToExecutable[len(pathToExecutable)-1]
	process.FatalCli("USAGE: %s ", executableName)
}

func setLogger(level string) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel(level),
		Writer:   os.Stderr,
	}
	filter.SetMinLevel(filter.MinLevel)
	log.SetOutput(filter)
}
