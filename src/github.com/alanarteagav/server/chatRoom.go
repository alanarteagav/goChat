package server

type ChatRoom struct {
    owner            Guest
    name             string
    connectionCount  int
    guests           map[string]Guest
    invitedGuests    map[string]Guest
}

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

func (chatRoom ChatRoom) getOwner() *Guest {
    return &chatRoom.owner
}

func (chatRoom ChatRoom) getName() string {
    return chatRoom.name
}

func (chatRoom *ChatRoom) setName(name string) {
    chatRoom.name = name
}

func (chatRoom ChatRoom) getConnectionCount() int {
    return chatRoom.connectionCount
}

func (chatRoom ChatRoom) getGuests() map[string]Guest {
    return chatRoom.guests
}

func (chatRoom ChatRoom) hosts(guest *Guest) bool {
    if _, ok := chatRoom.guests[guest.getUsername()]; ok{
        return true
    }
    return false
}

func (chatRoom *ChatRoom) addGuest(guest Guest) {
    if chatRoom.wasInvited(&guest) {
        chatRoom.connectionCount++
        chatRoom.guests[guest.getUsername()] = guest
    }
}

func (chatRoom *ChatRoom) addInvitedGuest(owner Guest, guest *Guest) bool {
    if !chatRoom.owner.equals(&owner) {
        return false
    } else {
        chatRoom.invitedGuests[guest.getUsername()] = *guest
        return true
    }
}

func (chatRoom ChatRoom) wasInvited(guest *Guest) bool {
    if _, ok := chatRoom.invitedGuests[guest.getUsername()]; ok{
        return true
    }
    return false
}

func (chatRoom *ChatRoom) removeGuest(guest Guest) {
    chatRoom.connectionCount--
    delete(chatRoom.guests, guest.getUsername())
}

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
