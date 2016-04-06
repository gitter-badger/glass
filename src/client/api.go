package client

import (
    "glassbox"
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


func (c Client) Dial(Entity) (net.Conn, error) {}
func (c Client) Send(glassbox.Packet) error {}

func (c Client) NewAuthorization() (string, error) {}

func (c Client) EntityFromAddress(addr net.Addr) glassbox.Entity {
    return nil
}
