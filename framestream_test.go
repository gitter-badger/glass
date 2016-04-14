package glass

import (
//    "net"
    //"io/ioutil"
    "testing"
    //"crypto/rsa"
    //"math/big"
    "fmt"
    //"sync"
)

var app *App

func ProcessTestFrame(p *TestFrame) {
    fmt.Println("[S] Frame Received Correctly. Exiting...")
    app.Close()
}

func Test(t *testing.T) {
    fmt.Println("[-] Starting")
    go func() {
        var err error
        var stream *FrameStream
        fmt.Println("[C] Client goroutine started. Dialing...")
        token := AuthToken{
            Router: &Peer{ Addr: "localhost:3001" },
        }
        var app = &App{ Token: token }
        defer app.Close()
        if stream, err = app.Connect(); err != nil {
            t.Fatal(err)
        }
        fmt.Println("[C] Connected")
        defer app.Close()

        var p = new(TestFrame)
        err = stream.Send(p)
        if err != nil {
            fmt.Println(err.Error())
        }
    }()

    fmt.Println("[S] Starting app server")
    app = &App{
        Token: AuthToken{ Me: &Peer{Addr:"localhost:3001"} },
        ProcessTestFrame: ProcessTestFrame,
    }
    fmt.Println("[S] Starting Handshake")
    if err := app.ListenAndServe(); err != nil {
        t.Fatal(err)
        return
    }
    app.Block()
}
