package glass

import (
    "testing"
    "fmt"
)

// Global app variable
// (in order to stop the app when finished)
var app *App

func ProcessTestFrame(p *TestFrame) {
    fmt.Println("[S] Frame Received Correctly. Exiting...")
    app.Close()
}

func Test(t *testing.T) {
    fmt.Println("[-] Starting")
    go func() {
        fmt.Println("[C] Client GoRoutine")
        var err error
        var stream *FrameStream
        // Authorization token with a peer's address sepcified
        token := AuthToken{
            Router: &Peer{ Addr: "localhost:3001" },
        }
        var app = &App{ Token: token }
        if stream, err = app.Connect(); err != nil {
            t.Fatal(err)
        }
        fmt.Println("[C] Client App Connected")
        defer app.Close()

        var p = new(TestFrame)
        if stream.Send(p) != nil {
            t.Fatal(err)
        }
        fmt.Println("[C] Test Frame Sent")
    }()

    fmt.Println("[S] Starting Server App")
    app = &App{
        Token: AuthToken{ Me: &Peer{Addr:"localhost:3001"} },
        ProcessTestFrame: ProcessTestFrame,
    }
    fmt.Println("[S] Listen & Server")
    if err := app.ListenAndServe(); err != nil {
        t.Fatal(err)
        return
    }
    app.Block()
}
