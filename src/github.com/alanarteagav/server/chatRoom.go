package server

type ChatRoom struct {
    name             string
    connectionCount  int
    guests           map[string]Guest
}

func NewChatRoom(name string) *ChatRoom {
    chatRoom := new(ChatRoom)
    return chatRoom
}

func (chatRoom ChatRoom) GetName() string {
    return ""
}

func (chatRoom *ChatRoom) SetName(name string) {

}

func (chatRoom ChatRoom) GetConnectionCount() int {
    return 0
}

func (chatRoom ChatRoom) GetGuests() map[string]Guest {
    return nil
}

func (chatRoom *ChatRoom) AddGuest(guest Guest) {

}

func (chatRoom *ChatRoom) RemoveGuest(guest Guest) {

}
