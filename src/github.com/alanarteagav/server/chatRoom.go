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
    chatRoom.guests[owner.GetUsername()] = owner
    return chatRoom
}

func (chatRoom ChatRoom) GetOwner() *Guest {
    return &chatRoom.owner
}

func (chatRoom ChatRoom) GetName() string {
    return chatRoom.name
}

func (chatRoom *ChatRoom) SetName(name string) {
    chatRoom.name = name
}

func (chatRoom ChatRoom) GetConnectionCount() int {
    return chatRoom.connectionCount
}

func (chatRoom ChatRoom) GetGuests() map[string]Guest {
    return chatRoom.guests
}

func (chatRoom ChatRoom) Hosts(guest *Guest) bool {
    if _, ok := chatRoom.guests[guest.GetUsername()]; ok{
        return true
    }
    return false
}

func (chatRoom *ChatRoom) AddGuest(guest Guest) {
    if chatRoom.WasInvited(&guest) {
        chatRoom.connectionCount++
        chatRoom.guests[guest.GetUsername()] = guest
    }
}

func (chatRoom *ChatRoom) AddInvitedGuest(owner Guest, guest *Guest) bool {
    if !chatRoom.owner.Equals(&owner) {
        return false
    } else {
        chatRoom.invitedGuests[guest.GetUsername()] = *guest
        return true
    }
}

func (chatRoom ChatRoom) WasInvited(guest *Guest) bool {
    if _, ok := chatRoom.invitedGuests[guest.GetUsername()]; ok{
        return true
    }
    return false
}

func (chatRoom *ChatRoom) RemoveGuest(guest Guest) {
    chatRoom.connectionCount--
    delete(chatRoom.guests, guest.GetUsername())
}

func (chatRoom *ChatRoom) Equals(cr *ChatRoom) bool {
    if chatRoom.name != cr.name {
        return false
    } else if chatRoom.connectionCount != cr.connectionCount {
        return false
    }
    for guestKey , guestsValue := range chatRoom.guests {
        guest := cr.guests[guestKey]
        if !guest.Equals(&guestsValue) {
            return false
        }
    }
    return true
}
