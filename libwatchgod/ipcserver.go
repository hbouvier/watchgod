package watchgod

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type IPCServer struct {
	watchgod *Watchgod
}

type RPCResponse struct {
	err error
	msg string
}

func StartIPCServer(url string, watchgod *Watchgod) {
	log.Printf("[INFO] [ipcserver] Starting IPC Server on %s\n", url)
	go func() {
		ipcServer := new(IPCServer)
		ipcServer.setwatchgod(watchgod)
		rpc.Register(ipcServer)
		rpc.HandleHTTP()
		listen, err := net.Listen("tcp", url)
		if err != nil {
			Fatal("StartIPCServer() listen error >>> %s", err)
		}
		http.Serve(listen, nil)
	}()
}

func (t *IPCServer) Add(args []string, reply *string) error {
	output := make(chan RPCResponse, 1)
	t.watchgod.Add(args[0], args[1:], output)
	return handlewatchgodResult(fmt.Sprintf("add %s %v", args[0], args), output, reply)
}

func (t *IPCServer) Restart(name *string, reply *string) error {
	output := make(chan RPCResponse, 1)
	t.watchgod.Restart(*name, output)
	return handlewatchgodResult(fmt.Sprintf("restart %s", *name), output, reply)
}

func (t *IPCServer) Start(name *string, reply *string) error {
	output := make(chan RPCResponse, 1)
	t.watchgod.Start(*name, output)
	return handlewatchgodResult(fmt.Sprintf("start %s", *name), output, reply)
}

func (t *IPCServer) Stop(name *string, reply *string) error {
	output := make(chan RPCResponse, 1)
	t.watchgod.Stop(*name, output)
	return handlewatchgodResult(fmt.Sprintf("stop %s", *name), output, reply)
}

func (t *IPCServer) List(filter *string, reply *string) error {
	output := make(chan RPCResponse, 1)
	t.watchgod.List(*filter, output)
	return handlewatchgodResult("list", output, reply)
}

func (t *IPCServer) Version(dummy *string, reply *string) error {
	output := make(chan RPCResponse, 1)
	t.watchgod.Version(output)
	return handlewatchgodResult("version", output, reply)
}

func (t *IPCServer) Terminate(reason *string, reply *string) error {
	output := make(chan RPCResponse, 1)
	t.watchgod.Terminate(*reason, output)
	return handlewatchgodResult("terminate", output, reply)
}

func (t *IPCServer) setwatchgod(watchgod *Watchgod) {
	t.watchgod = watchgod
}

func handlewatchgodResult(msg string, output chan RPCResponse, reply *string) error {
	var result error = nil

	response := <-output
	if response.err != nil {
		*reply = fmt.Sprintf("%s >>> %s", msg, response.err)
		result = response.err
	} else {
		*reply = response.msg
	}
	return result
}
