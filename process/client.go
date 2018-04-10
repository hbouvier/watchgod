package process

import (
	"fmt"
	"net/rpc"
)

// Client ...
type Client struct {
	url string
}

// NewClient ...
func NewClient(url string) Client {
	client := Client{url: url}
	return client
}

// Version ...
func (c *Client) Version() {
	c.call("Version", "*")
}

// List ...
func (c *Client) List() {
	c.call("List", "*")
}

// Terminate ...
func (c *Client) Terminate() {
	c.call("Terminate", "now")
}

// Start ...
func (c *Client) Start(name string) {
	c.call("Start", name)
}

// Stop ...
func (c *Client) Stop(name string) {
	c.call("Stop", name)
}

// Restart ...
func (c *Client) Restart(name string) {
	c.call("Restart", name)
}

// Add ...
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
