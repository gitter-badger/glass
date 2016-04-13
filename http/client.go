package http

import (
    "github.com/acondolu/glass"
    "net/http"
    "time"
    "bufio"
    "errors"
    "bytes"
)

var app *glass.App

func Init(auth glass.AuthToken, ready func()) {
    app.Init(auth)
}

type Transport struct {
    peer glass.Peer
    streams [](chan<- []byte)

    Timeout time.Duration
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
    // Get stream count
    streamNumber := len(t.streams)
    // Opean a new channel
    stream := make(chan []byte, 0)
    // Store it into the streams list
    t.streams = append(t.streams, stream)
    // send frame with streamNumber set
    // TODO
    streamNumber = streamNumber
    // Set up timeout
    go func(){
        time.Sleep(t.Timeout * time.Second)
        stream <- nil
    }()
    // Wait for reply
    reply := <-stream
    // Close stream
    close(stream)
    // Parse and return reply
    if reply == nil {
        return nil, errors.New("glass/http: timeout awaiting reply")
    }
    if string(reply[:4]) != "HTTP" {
        // TODO Check for gzip compression
        // TODO return error
        return nil, errors.New("glass/http: HTTP/1 transport connection broken")
    }
    r := bufio.NewReader(bytes.NewReader(reply))
    return http.ReadResponse(r, req)
}
