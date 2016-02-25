package main

import (
	"log"
	"net/http"
//	"time"
	"crypto/rsa"
	"net/rpc"
	"net"
	"path/filepath"
	"os/exec"
	"fmt"
	"os"
)

type GlassClient struct {
	// events queue
	key *rsa.PrivateKey
	publicKeys map[string]*rsa.PublicKey
	rpcaddr string
	rpcconn net.Listener
	apps map[string]string
	addrs map[string]string
	outs map[string]*rpc.Client
}

type Payload struct {
	PublicHeader  map[string]string
	PrivateHeader map[string]string
	Sender string
	Content string
}

func (c *GlassClient) RegisterApp(authtoken string, reply *int) error {
	addr := c.addrs[authtoken]
	conn, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return err
	}
	c.outs[authtoken] = conn
	return nil
}

func (c *GlassClient) StartApp(appid string) error {
	authtoken := "!!!FIXME!!!"
	addr := "localhost:1235"
	c.addrs[authtoken] = addr
	c.apps[appid] = authtoken

	path := filepath.Join("apps", appid, "main.go")
	cmd := exec.Command("go", "run", path, authtoken, addr, c.rpcaddr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	defer cmd.Run()
	return nil
}

/*func (c *GlassClient) RegisterApp(appid string) (<-chan *Payload, chan<- *Payload) {
	in, out := make(chan *Resource), make(chan *Resource)
	c.ins[appid] = in
	c.outs[appid] = out
	for r := range in {
		out <- Payload{nil, nil, "testSender", "testContent"}
	}
	return in, out
}
*/

/*
func Chat(in <-chan *Payload, out chan<- *Payload) {
	for r := range in {
		out <- Payload{nil, nil, "testSender", "testContent"}
	}
	while true

}
*/

func main() {
	fmt.Println("Starting Glass")
	glass := new(GlassClient)
	glass.rpcaddr = "localhost:1234"
	rpc.Register(glass)
	rpc.HandleHTTP()
	conn, err := net.Listen("tcp", glass.rpcaddr)
	if err != nil {
		log.Fatal("Can't start rpc interface", err)
		os.Exit(1)
	}
	glass.rpcconn = conn
	glass.addrs = make(map[string]string)
	glass.apps = make(map[string]string)
	glass.outs = make(map[string]*rpc.Client)

	go func(){
		err = glass.StartApp("SimpleApp")
		if err != nil {
			log.Fatal("Error starting SimpleApp", err)
		}
	}()
	http.Serve(conn, nil)

}
