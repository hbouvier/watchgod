package watchgod

import (
	//  "encoding/json"
	"fmt"
	"net/rpc"
)

type Client struct {
	url string
}

func NewClient(url string) Client {
	client := Client{url: url}
	return client
}

func (c *Client) Version() {
	c.call("Version", "*")
}

func (c *Client) List() {
	c.call("List", "*")
}

func (c *Client) Terminate() {
	c.call("Terminate", "now")
}

func (c *Client) Start(name string) {
	c.call("Start", name)
}

func (c *Client) Stop(name string) {
	c.call("Stop", name)
}

func (c *Client) Restart(name string) {
	c.call("Restart", name)
}

func (c *Client) Add(args []string) {
	var reply string
	err := c.rpc().Call("IPCServer.Add", args, &reply)
	if err != nil {
		FatalCli("Error: %s", err)
	}
	fmt.Printf("%s\n", reply)
}

func (c *Client) call(method string, argument string) {
	var reply string
	err := c.rpc().Call("IPCServer."+method, argument, &reply)
	if err != nil {
		FatalCli("Error: %s", err)
	}
	fmt.Printf("%s\n", reply)
}

func (c *Client) rpc() *rpc.Client {
	client, err := rpc.DialHTTP("tcp", c.url)
	if err != nil {
		FatalCli("Error connecting: %s", err)
	}
	return client
}
