package client

import (
    "net"
    "github.com/acondolu/glassbox"
)

type PacketHandlerFunc func(p glassbox.Packet) bool
type ConnHandlerFunc func(c net.Conn)


type Client struct {
    packet_handler PacketHandlerFunc
    connection_handler ConnHandlerFunc
}

func (c Client) Init(auth string) {}

func (c Client) OnNewPacket(handler PacketHandlerFunc) {}
func (c Client) OnNewConnection(handler ConnHandlerFunc) {}


func (c Client) Dial(glassbox.Entity) (net.Conn, error) {
    return nil, nil
}
func (c Client) Send(glassbox.Packet) error {
    return nil;
}

func (c Client) NewAuthorization() (string, error) {
    return "", nil
}

//func (c Client) EntityFromAddress(addr net.Addr) glassbox.Entity {
//}
