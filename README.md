[![Build Status](https://travis-ci.org/hbouvier/watchgod.png)](https://travis-ci.org/hbouvier/watchgod)
[![Coverage Status](https://coveralls.io/repos/hbouvier/watchgod/badge.svg?branch=master)](https://coveralls.io/r/hbouvier/watchgod?branch=master)

# WatchGOd

WatchGOd is a Deamon process WATCHer written in GO that can be useful as the **init** process of docker containers. It will automatically restart a dying process. With WatchGOd, you can also manage the lifecycle of a process: start, retstart and stop.

## Simple Use cases

- Manually configure WatchGOd

```bash
$ watchgod boot &
$ watchgod add greeter /bin/sh -c 'echo Hello world && sleep 10 && echo Good bye'
$ watchgod start greeter
$ sleep 11 && watchgod terminate
```

- Using a configuration file

```bash
$ echo '{"Processes":[{"Name" : "sleeper","Command" : ["/bin/sh", "-c", "echo ZZZZZ && sleep 10"]}]}' > /tmp/configuration.json
$ watchgod -config /tmp/configuration.json boot &
$ sleep 11 && watchgod terminate
```

## Other useful commands

```bash
$ echo '{"Processes":[{"Name" : "sleeper","Command" : ["/bin/sh", "-c", "echo ZZZZZ && sleep 60"]}]}' > /tmp/configuration.json
$ watchgod -config /tmp/configuration.json boot &
015/10/03 09:43:31 [INFO] [ipcserver] Starting IPC Server on 127.0.0.1:7099
2015/10/03 09:43:31 [INFO] [watchgod] Daemon version 1.0.1 started...
2015/10/03 09:43:31 [INFO] [watchgod] sleeper: created
ZZZZZ
2015/10/03 09:43:32 [INFO] [watchgod] sleeper: started with PID 48806
$ watchgod version
Client version 1.0.1
Deamon version 1.0.1
$ watchgod list
 PID   NAME                 STATE
---------------------------------------
[48806] sleeper              RUNNING

$ watchgod stop sleeper
2015/10/03 09:44:50 [INFO] [watchgod] sleeper: PID 48806 exit code -1
2015/10/03 09:44:50 [INFO] [watchgod] sleeper: stopped

$ watchgod list
 PID   NAME                 STATE
---------------------------------------
[   0] sleeper              DEAD

$ watchgod start sleeper
ZZZZZ
2015/10/03 09:45:28 [INFO] [watchgod] sleeper: started with PID 48983
sleeper: started with PID 48983
$ watchgod list
PID   NAME                 STATE
---------------------------------------
[48983] sleeper              RUNNING

$ watchgod restart sleeper
2015/10/03 09:46:31 [INFO] [watchgod] sleeper: PID 48983 exit code -1
ZZZZZ
2015/10/03 09:46:32 [INFO] [watchgod] sleeper: started with PID 49081
sleeper: started with PID 49081

$ watchgod terminate
2015/10/03 09:46:56 [INFO] [watchgod] sleeper: PID 49081 exit code -1
2015/10/03 09:46:56 [INFO] [watchgod] terminated
2015/10/03 09:46:56 [INFO] [watchgod] WatchGOd: terminated
WatchGOd: terminated
```

## Configuration file
```json
{
    "StartTimeoutInSeconds": 1,
    "StopTimeoutInSeconds": 10,
    "IPCServerUrl" : "127.0.0.1:7099",
    "LogLevel" : "DEBUG",
    "Processes" : [
      {
        "Name" : "hello",
        "Command" : ["/bin/sh", "-c", "echo Hello World && sleep 10"]
      },
      {
        "Name" : "sleeper",
        "Command" : ["sleep", "30"]
      }
    ]
}
```

## Build

To build this project on your computer if you already have go-lang installed

```bash
$ make
```

Or if you want to build a linux version using docker:

```bash
$ make GOCC="docker run --rm -t -v ${GOPATH}:/go hbouvier/go-lang:1.5"
```

