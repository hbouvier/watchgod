# WatchGOd is watchdog Daemon in Go

## How to use it

```bash
$ watchgod boot &
$ watchgod version
$ watchgod add sleeper sleep 60
$ watchgod start sleeper
$ watchgod add hello /bin/sh -c 'echo Hello World && sleep 15'
$ watchgod start hello
$ watchgod list
$ watchgod stop hello
$ watchgod stop sleeper
$ watchgod start hello
$ watchgod restart hello
$ watchgod terminate
```


