package glass

import (
    "net"
    //"io/ioutil"
    "testing"
    //"crypto/rsa"
    //"math/big"
    "fmt"
    //"sync"
)

var fs FrameStream

func ProcessTestFrame(p *TestFrame) {
    fmt.Println("[S] Frame Received Correctly. Exiting...")
    fs.Shutdown(nil)
}

func Test(t *testing.T) {
    fmt.Println("[-] Starting")
    go func() {
        fmt.Println("[C] Client goroutine started. Dialing...")
        conn, err := net.Dial("tcp", "localhost:3001")
        if err != nil {
            t.Fatal(err)
        }
        fmt.Println("[C] Connected")
        defer conn.Close()

        ps := new(FrameStream)
        fmt.Println("[C] Starting Handshake")
        if err := ps.Out(nil, conn); err != nil {
            t.Fatal(err.Error())
        }
        fmt.Println("[C] Handshake Over")
        p := new(TestFrame)
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
    fmt.Println("[S] Listening to connections")
    //for {
        conn, err := l.Accept()
        fmt.Println("[S] Connection accepted")
        if err != nil {
            return
        }
        defer conn.Close()
        defer conn.Close()
        app := &App{ProcessTestFrame: ProcessTestFrame}
        fmt.Println("[S] Starting Handshake")
        if err = fs.In(app, conn); err != nil {
            t.Fatal(err)
            return
        }
        fmt.Println("[S] Handshake Over")
        //_ = ps.In;
        //buf, err := ioutil.ReadAll(conn)
        //if err != nil {
        //    t.Fatal(err)
        //}

        //fmt.Println("Received %d bytes: %s.", len(buf), string(buf[:]))
        //if msg := string(buf[:]); msg != message {
        //    t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, message)
        //}
        fs.Serve()
    //    break
    //}
    //return // Done
}
