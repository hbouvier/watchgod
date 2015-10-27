package main

import (
	"flag"
	"github.com/hashicorp/logutils"
	"github.com/hbouvier/watchgod/libwatchgod"
	"log"
	"os"
)

func main() {
	configuration := watchgod.DefaultConfiguration()
	configurationFilenamePtr := flag.String("config", "", "configuration filename.json")
	levelPtr := flag.String("level", "", "logging level")
	flag.Parse()
	if *configurationFilenamePtr != "" {
		configuration = watchgod.LoadConfiguration(*configurationFilenamePtr)
	}
	if *levelPtr != "" {
		configuration.LogLevel = *levelPtr
	}
	setLogger(configuration.LogLevel)
	watchgod.ExecuteArgument(flag.Args(), configuration, usage)
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
