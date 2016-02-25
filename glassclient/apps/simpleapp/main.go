package main

import (
	"os"
	"log"
	"fmt"
	"net"
	"net/rpc"
	"net/http"
)

type Payload struct {
	PublicHeader  map[string]string
	PrivateHeader map[string]string
	Sender string
	Content string
}

type SimpleApp struct {
	// events queue
	in net.Listener
	out *rpc.Client
}

type AppConfig struct {
	AuthToken string
	InAddr string
	OutAddr string
}

func (app *SimpleApp) Receive(payload *Payload, reply *int) error {
	return nil
}

func (app *SimpleApp) Init(config *AppConfig) error {
	// out
	oconn, err := rpc.DialHTTP("tcp", config.OutAddr)
	if err != nil {
		return err
	}
	app.out = oconn
	// in
	rpc.Register(app)
	rpc.HandleHTTP()
	iconn, err := net.Listen("tcp", config.InAddr)
	if err != nil {
		return err
	}
	app.in = iconn
	reply := new(int)
	defer app.out.Go("GlassClient.RegisterApp", config.AuthToken, reply, nil)
	return http.Serve(iconn, nil)
}

/*
func Chat(in <-chan *Payload, out chan<- *Payload) {
	for r := range in {
		out <- Payload{nil, nil, "testSender", "testContent"}
	}
	while true

}
*/

func main() {
	fmt.Println("Hello SimpleApp")
	config := &AppConfig{os.Args[1], os.Args[2], os.Args[3]}
	app := new(SimpleApp)
	err := app.Init(config)
	log.Fatal("SimpleApp terminated: ", err)
}
