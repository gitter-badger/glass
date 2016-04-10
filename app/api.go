package app

import (
    "net"
    "github.com/acondolu/glassbox"
)

type PacketHandlerFunc func(p glassbox.Packet) bool
type ConnHandlerFunc func(c net.Conn)


type App struct {
    packet_handler PacketHandlerFunc
    connection_handler ConnHandlerFunc
}

func (app App) Init(auth string) {}

func (app App) OnNewPacket(handler PacketHandlerFunc) {}
func (app App) OnNewConnection(handler ConnHandlerFunc) {}


func (app App) Dial(glassbox.Entity) (net.Conn, error) {
    return nil, nil
}
func (app App) Send(glassbox.Packet) error {
    return nil;
}

func (app App) NewAuthorization() (string, error) {
    return "", nil
}

//func (app App) EntityFromAddress(addr net.Addr) glassbox.Entity {
//}
