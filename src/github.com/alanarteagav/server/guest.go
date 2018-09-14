package server

import "net"

// Guest struct (auxiliar for server).
// Defines an username, and the guest connection
// (golang's equivalent for sockets).
type Guest struct {
    serial int
    username   string
    connection net.Conn
    chatRooms  map[string]ChatRoom
}

// Guest constructor.
func NewGuest(serial int, connection net.Conn) *Guest {
    guest := new(Guest)
    guest.serial = serial
    guest.username = ""
    guest.connection = connection
    guest.chatRooms = make(map[string]ChatRoom)
    return guest
}

// Returns guest's connection.
func (guest Guest) GetConnection() net.Conn {
    return guest.connection
}

// Returns guest's serial.
func (guest Guest) GetSerial() int {
    return guest.serial
}

// Returns guest's username.
func (guest Guest) GetUsername() string {
    return guest.username
}

// Sets a new username for the guest.
func (guest *Guest) SetUsername(username string) {
    guest.username = username
}

func (guest Guest) JoinChatRoom(chatRoom ChatRoom) {
    _, isInChatRooms := guest.chatRooms[chatRoom.GetName()]
    if !isInChatRooms {
        guest.chatRooms[chatRoom.GetName()] = chatRoom
    }
}

func (guest Guest) LeaveChatRoom(chatRoom ChatRoom) {
    guestChatRoom, isInChatRooms := guest.chatRooms[chatRoom.GetName()]
    if isInChatRooms {
        delete(guest.chatRooms, guestChatRoom.GetName())
    }
}

func (guest Guest) GetChatRooms() map[string]ChatRoom {
    return nil
}

func (guest Guest) IsInChatRoom(chatRoom ChatRoom) bool {
    _, isInChatRooms := guest.chatRooms[chatRoom.GetName()]
    return isInChatRooms
}

func (guest Guest) Equals(g *Guest) bool {
    if guest.serial != g.serial {
        return false
    }
    if guest.username != g.username {
        return false
    }
    return true
}
