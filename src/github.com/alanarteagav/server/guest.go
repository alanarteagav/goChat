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
    status UserStatus
}

// User status constants.
type UserStatus string
const (
    BUSY UserStatus = "BUSY"
    AWAY  UserStatus = "AWAY"
    ACTIVE  UserStatus = "ACTIVE"
    UNDEFINED  UserStatus = "UNDEFINED"
)

// Function that checks if a string is a valid status constant.
func toUserStatus(str string) UserStatus {
    switch str {
    case "BUSY":
        return BUSY
    case "AWAY":
        return AWAY
    case "ACTIVE":
        return ACTIVE
    default:
        return UNDEFINED
    }
}


// Guest constructor.
func NewGuest(serial int, connection net.Conn) *Guest {
    guest := new(Guest)
    guest.serial = serial
    guest.username = ""
    guest.connection = connection
    guest.chatRooms = make(map[string]ChatRoom)
    guest.status = ACTIVE
    return guest
}

// Returns guest's connection.
func (guest Guest) getConnection() net.Conn {
    return guest.connection
}

// Returns guest's serial.
func (guest Guest) getSerial() int {
    return guest.serial
}

// Returns guest's username.
func (guest Guest) getUsername() string {
    return guest.username
}

// Checks if the guest is identified.
func (guest Guest) isIdentified() bool {
    return guest.username != ""
}

// Sets a new username for the guest.
func (guest *Guest) setUsername(username string) {
    guest.username = username
}

// Sets a new status for the guest.
func (guest *Guest) setStatus(status UserStatus) {
    guest.status = status
}

//Gets the guest status.
func (guest Guest) getStatus() UserStatus {
    return guest.status
}

// Attaches the guest to a chatroom.
func (guest Guest) joinChatRoom(chatRoom ChatRoom) {
    _, isInChatRooms := guest.chatRooms[chatRoom.getName()]
    if !isInChatRooms {
        guest.chatRooms[chatRoom.getName()] = chatRoom
    }
}

// Removes the guest from a chatroom.
func (guest Guest) leaveChatRoom(chatRoom ChatRoom) {
    guestChatRoom, isInChatRooms := guest.chatRooms[chatRoom.getName()]
    if isInChatRooms {
        delete(guest.chatRooms, guestChatRoom.getName())
    }
}

// Checks if a guest is an specific chatroom.
func (guest Guest) isInChatRoom(chatRoom ChatRoom) bool {
    _, isInChatRooms := guest.chatRooms[chatRoom.getName()]
    return isInChatRooms
}

// Checks if two guests are equal.
func (guest Guest) equals(g *Guest) bool {
    if guest.serial != g.serial {
        return false
    }
    if guest.username != g.username {
        return false
    }
    return true
}
