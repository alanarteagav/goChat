package server

import "net"

// Guest struct (auxiliar for server).
// Defines an username, and the guest connection
// (golang's equivalent for sockets).
type Guest struct {
    username   string
    connection net.Conn
    chatRooms  map[string]ChatRoom
}

// Guest constructor.
func NewGuest(username string, connection net.Conn) *Guest {
    guest := new(Guest)
    guest.username = username
    guest.connection = connection
    return guest
}

// Returns guest's connection.
func (guest Guest) GetConnection() net.Conn {
    return guest.connection
}

// Returns guest's username.
func (guest Guest) GetUsername() string {
    return guest.username
}

// Sets a new username for the guest.
func (guest *Guest) SetUsername(username string) {
    guest.username = username
}

func (guest Guest) JoinChatRoom(chatRoom ChatRoom) bool {
    return false
}

func (guest Guest) LeaveChatRoom(chatRoom ChatRoom) bool {
    return false
}

func (guest Guest) GetChatRooms() map[string]ChatRoom {
    return nil
}

func (guest Guest) IsInChatRoom(chatRoom ChatRoom) bool {
    return false
}

func (guest Guest) Equals(g *Guest) bool {
    return false
}
