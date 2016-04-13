package http

import (
    "github.com/acondolu/glass"
    "strings"
    "net/http"
    "time"
    "bufio"
    "errors"
)

var app *glass.App

func Init(auth string) {
    app.Init(auth)
}

type Transport struct {
    peer glass.Peer
    streams [](chan<- glass.Frame)

    Timeout time.Duration
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
    // Get stream count
    streamNumber := len(t.streams)
    // Opean a new channel
    stream := make(chan glass.Frame, 0)
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
    var payload = reply.Content()
    if string(payload[:4]) != "HTTP" {
        // TODO Check for gzip compression
        // TODO return error
        return nil, errors.New("glass/http: HTTP/1 transport connection broken")
    }
    var str = string(payload)
    r := bufio.NewReader(strings.NewReader(str))
    return http.ReadResponse(r, req)
}
