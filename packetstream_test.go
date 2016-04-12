package glassbox

import (
    "net"
    //"io/ioutil"
    "testing"
    //"crypto/rsa"
    //"math/big"
    "fmt"
    //"sync"
)

type SimpleInstance struct {
    PS *PacketStream
}
func (i *SimpleInstance) ProcessSimplePacket(p *SimplePacket) {}
func (i *SimpleInstance) ProcessTestPacket(p *TestPacket) {
    fmt.Println("[--] Packet Received Correctly. Exiting...")
    i.PS.Shutdown(nil)
}

func Test(t *testing.T) {
    fmt.Println("[--] Starting")
    go func() {
        fmt.Println("[->] Client goroutine started. Dialing...")
        conn, err := net.Dial("tcp", "localhost:3001")
        if err != nil {
            t.Fatal(err)
        }
        fmt.Println("[->] Connected")
        defer conn.Close()

        ps := new(PacketStream)
        fmt.Println("[->] Starting Handshake")
        if err := ps.Out(nil, conn, nil); err != nil {
            t.Fatal(err.Error())
        }
        fmt.Println("[->] Handshake Over")
        p := new(TestPacket)
        err = ps.Send(p)
        if err != nil {
            fmt.Println(err.Error())
        }
        //ps.init(STREAM_)
        //        t.Fatal(g, e)
        conn.Close()
    }()

    l, err := net.Listen("tcp", "localhost:3001")
    if err != nil {
        t.Fatal(err)
    }
    defer l.Close()
    fmt.Println("[<-] Listening to connections")
    //for {
        conn, err := l.Accept()
        fmt.Println("[<-] Connection accepted")
        if err != nil {
            return
        }
        defer conn.Close()
        ps := new(PacketStream)
        inst := new(SimpleInstance)
        inst.PS = ps
        fmt.Println("[<-] Starting Handshake")
        if err = ps.In(inst, conn, nil); err != nil {
            t.Fatal(err)
            return
        }
        fmt.Println("[<-] Handshake Over")
        //_ = ps.In;
        //buf, err := ioutil.ReadAll(conn)
        //if err != nil {
        //    t.Fatal(err)
        //}

        //fmt.Println("Received %d bytes: %s.", len(buf), string(buf[:]))
        //if msg := string(buf[:]); msg != message {
        //    t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, message)
        //}
        ps.Serve()
    //    break
    //}
    //return // Done
}
