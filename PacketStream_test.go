package glassbox

import (
    "net"
    "io/ioutil"
    "testing"
    "fmt"
)

func Test(t *testing.T) {
    // Start in a go rouftine
    go func() {
        conn, err := net.Dial("tcp", "localhost:3000")
        if err != nil {
            t.Fatal(err)
        }
        //defer conn.Close()

        ps := new(PacketStream)
        ps.conn = conn
        p := new(TestPacket)
        ps.Write(p)
        //ps.init(STREAM_)
        //        t.Fatal(g, e)
        conn.Close()
    }()

    l, err := net.Listen("tcp", "localhost:3000")
    if err != nil {
        t.Fatal(err)
    }
    defer l.Close()
    for {
        conn, err := l.Accept()
        fmt.Println("Connection accepted")
        if err != nil {
            return
        }
        defer conn.Close()

        buf, err := ioutil.ReadAll(conn)
        if err != nil {
            t.Fatal(err)
        }

        t.Logf("Received %d bytes: %s.", len(buf), string(buf[:]))
        //if msg := string(buf[:]); msg != message {
        //    t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, message)
        //}
        break
    }
    return // Done
}
