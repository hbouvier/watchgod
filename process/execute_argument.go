package process

import (
	"fmt"
)

// CommandFlags ...
type CommandFlags struct {
	ListPtr      *bool
	TerminatePtr *bool
	VersionPtr   *bool
	StartPtr     *string
	StopPtr      *string
	RestartPtr   *string
	UserPtr      *string
	AddPtr       *string
}

// ExecuteArgument ...
func ExecuteArgument(commandFlags CommandFlags, arguments []string, configuration Configuration, version string, usage func()) {
	if *commandFlags.UserPtr != "" {
		if err := Su(*commandFlags.UserPtr); err != nil {
			FatalCli("Unable to change user to '%s' => %+v", *commandFlags.UserPtr, err)
		}
	}

	if *commandFlags.VersionPtr {
		fmt.Printf("Client version %s\n", version)
		client := NewClient(IpcServerURL(DefaultConfiguration().IPCServerURL))
		client.Version()
	} else if *commandFlags.ListPtr {
		client := NewClient(IpcServerURL(DefaultConfiguration().IPCServerURL))
		client.List()
	} else if *commandFlags.TerminatePtr {
		client := NewClient(IpcServerURL(DefaultConfiguration().IPCServerURL))
		client.Terminate()
	} else if *commandFlags.StartPtr != "" {
		client := NewClient(IpcServerURL(DefaultConfiguration().IPCServerURL))
		client.Start(*commandFlags.StartPtr)
	} else if *commandFlags.StopPtr != "" {
		client := NewClient(IpcServerURL(DefaultConfiguration().IPCServerURL))
		client.Stop(*commandFlags.StopPtr)
	} else if *commandFlags.RestartPtr != "" {
		client := NewClient(IpcServerURL(DefaultConfiguration().IPCServerURL))
		client.Restart(*commandFlags.RestartPtr)
	} else if *commandFlags.AddPtr != "" {
		client := NewClient(IpcServerURL(DefaultConfiguration().IPCServerURL))
		client.Add(arguments[1:])
	} else {
		Boot(IpcServerURL(configuration.IPCServerURL), arguments, configuration, version)
	}
}
