package server

type ChatRoom struct {
    name             string
    connectionCount  int
    guests           map[string]Guest
}

func NewChatRoom(name string) *ChatRoom {
    chatRoom := new(ChatRoom)
    chatRoom.connectionCount = 0
    chatRoom.name = name
    chatRoom.guests = make(map[string]Guest)
    return chatRoom
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

func (chatRoom *ChatRoom) AddGuest(guest Guest) {
    chatRoom.connectionCount++
    chatRoom.guests[guest.GetUsername()] = guest
}

func (chatRoom *ChatRoom) RemoveGuest(guest Guest) {
    chatRoom.connectionCount--
    delete(chatRoom.guests, guest.GetUsername())
}

func (chatRoom *ChatRoom) Equals(cr *ChatRoom) bool {
    return false
}
