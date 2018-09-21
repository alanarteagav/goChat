package server

type ChatRoom struct {
    owner            Guest
    name             string
    connectionCount  int
    guests           map[string]Guest
    invitedGuests    map[string]Guest
}

// ChatRoom constructor.
func NewChatRoom(owner Guest, name string) *ChatRoom {
    chatRoom := new(ChatRoom)
    chatRoom.owner = owner
    chatRoom.connectionCount = 1
    chatRoom.name = name
    chatRoom.guests = make(map[string]Guest)
    chatRoom.invitedGuests = make(map[string]Guest)
    chatRoom.guests[owner.getUsername()] = owner
    return chatRoom
}

// Returns the chatRoom owner.
func (chatRoom ChatRoom) getOwner() *Guest {
    return &chatRoom.owner
}

// Sets a new name to the chatroom.
func (chatRoom ChatRoom) getName() string {
    return chatRoom.name
}

// Gets the chatroom's name.
func (chatRoom *ChatRoom) setName(name string) {
    chatRoom.name = name
}

// Gets the chatroom's connection count.
func (chatRoom ChatRoom) getConnectionCount() int {
    return chatRoom.connectionCount
}

// Gets the chatroom's guests hash table.
func (chatRoom ChatRoom) getGuests() map[string]Guest {
    return chatRoom.guests
}

// Checks if the chatroom hosts an specific guest.
func (chatRoom ChatRoom) hosts(guest *Guest) bool {
    if _, ok := chatRoom.guests[guest.getUsername()]; ok{
        return true
    }
    return false
}

// Adds a new guest to the chatRoom.
func (chatRoom *ChatRoom) addGuest(guest Guest) {
    if chatRoom.wasInvited(&guest) {
        chatRoom.connectionCount++
        chatRoom.guests[guest.getUsername()] = guest
    }
}

// Adds a new guest to the invitedGuests hash table.
func (chatRoom *ChatRoom) addInvitedGuest(owner Guest, guest *Guest) bool {
    if !chatRoom.owner.equals(&owner) {
        return false
    } else {
        chatRoom.invitedGuests[guest.getUsername()] = *guest
        return true
    }
}

// Checks if a guest was invited to the chatRoom.
// (if it is in the invitedGuests hash table).
func (chatRoom ChatRoom) wasInvited(guest *Guest) bool {
    if _, ok := chatRoom.invitedGuests[guest.getUsername()]; ok{
        return true
    }
    return false
}

// Removes a guest from the chatRoom.
func (chatRoom *ChatRoom) removeGuest(guest Guest) {
    chatRoom.connectionCount--
    delete(chatRoom.guests, guest.getUsername())
}

// Checks if two chatRooms are equal.
func (chatRoom *ChatRoom) equals(cr *ChatRoom) bool {
    if chatRoom.name != cr.name {
        return false
    } else if chatRoom.connectionCount != cr.connectionCount {
        return false
    }
    for guestKey , guestsValue := range chatRoom.guests {
        guest := cr.guests[guestKey]
        if !guest.equals(&guestsValue) {
            return false
        }
    }
    return true
}
