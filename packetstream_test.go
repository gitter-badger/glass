package glassbox

import (
    "net"
    "io/ioutil"
    "testing"
    //"crypto/rsa"
    //"math/big"
    "fmt"
)

func Test(t *testing.T) {
    // Key from the RSA module tests
    /*
    priv := &rsa.PrivateKey{
    PublicKey: rsa.PublicKey{
        N: fromBase10("290684273230919398108010081414538931343"),
        E: 65537,
    },
    D: fromBase10("31877380284581499213530787347443987241"),
    Primes: []*big.Int{
        fromBase10("16775196964030542637"),
        fromBase10("17328218193455850539"),
    },
    }
    t.Logf("%d", priv.D)
    */
    fmt.Println("Starting")
    go func() {
        fmt.Println("client goroutine started. Dialing...")
        conn, err := net.Dial("tcp", "localhost:3001")
        if err != nil {
            t.Fatal(err)
        }
        defer conn.Close()

        ps := new(PacketStream)
        fmt.Println("Starting handshake 1")
        ps.In(nil, conn, nil)
        p := new(TestPacket)
        ps.Write(p)
        //ps.init(STREAM_)
        //        t.Fatal(g, e)
        conn.Close()
    }()

    l, err := net.Listen("tcp", "localhost:3001")
    if err != nil {
        t.Fatal(err)
    }
    defer l.Close()
    fmt.Println("Listening to connections")
    for {
        conn, err := l.Accept()
        fmt.Println("Connection accepted")
        if err != nil {
            return
        }
        defer conn.Close()
        ps := new(PacketStream)
        fmt.Println("Starting handshake 2")
        ps.Out(nil, conn, nil)
        //_ = ps.In;
        buf, err := ioutil.ReadAll(conn)
        if err != nil {
            t.Fatal(err)
        }

        fmt.Println("Received %d bytes: %s.", len(buf), string(buf[:]))
        //if msg := string(buf[:]); msg != message {
        //    t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, message)
        //}
        break
    }
    return // Done
}
