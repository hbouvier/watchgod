package watchgod

import (
	"fmt"
)

func ExecuteArgument(arguments []string, configuration Configuration, usage func()) {
	nbArgs := len(arguments)

	switch nbArgs {
	case 0:
		usage()
	case 1:
		switch arguments[0] {
		case "boot":
			Boot(IpcServerUrl(configuration.IPCServerURL), configuration)
		case "list":
			client := NewClient(IpcServerUrl(DefaultConfiguration().IPCServerURL))
			client.List()
		case "terminate":
			client := NewClient(IpcServerUrl(DefaultConfiguration().IPCServerURL))
			client.Terminate()
		case "version":
			fmt.Printf("Client version %s\n", VERSION)
			client := NewClient(IpcServerUrl(DefaultConfiguration().IPCServerURL))
			client.Version()
		default:
			usage()
		}
	case 2:
		switch arguments[0] {
		case "start":
			client := NewClient(IpcServerUrl(DefaultConfiguration().IPCServerURL))
			client.Start(arguments[1])
		case "stop":
			client := NewClient(IpcServerUrl(DefaultConfiguration().IPCServerURL))
			client.Stop(arguments[1])
		case "restart":
			client := NewClient(IpcServerUrl(DefaultConfiguration().IPCServerURL))
			client.Restart(arguments[1])
		default:
			usage()
		}
	default:
		switch arguments[0] {
		case "add":
			client := NewClient(IpcServerUrl(DefaultConfiguration().IPCServerURL))
			client.Add(arguments[1:])
		default:
			usage()
		}
	}
}
