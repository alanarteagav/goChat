package server

import "net"

type Guest struct {

}

func NewGuest(name string, connection net.Conn) *Guest {
    return new(Guest)
}

func (guest Guest) GetUsername() string {
    return ""
}

func (guest Guest) SetUsername(username string) {

}
