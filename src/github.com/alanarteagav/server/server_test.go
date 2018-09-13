package server

import (
    "bufio"
    "flag"
    "github.com/alanarteagav/events"
    "log"
    "math/rand"
    "net"
    "strings"
    "os"
    "testing"
)

var chatServer Server
var tClient testClient

//Function that returns an aleatoryPort.
func aleatoryPort() int {
    rand.Seed(113)
    return 4000
}

//A client type to be used in the server's test cases.
type testClient struct {
    connection net.Conn
    reader     bufio.Reader
    writer     bufio.Writer
    quit       chan bool
}

//The testClient constructor.
func newClient(port int) *testClient {
    client := new(testClient)
    connection, err := net.Dial("tcp", "localhost:4000")
    if err != nil {
        log.Fatalln(err)
    }
    client.connection = connection
    client.reader = *bufio.NewReader(connection)
    client.writer = *bufio.NewWriter(connection)
    client.quit = make(chan bool)
    return client
}

//testClient method to send messages to the server.
func (client testClient) sendMessage(message string)  {
    client.connection.Write([]byte(message + "\n"))
}

//testClient method to send events to the server.
func (client testClient) sendEvent(event events.ChatEvent)  {
    client.writer.WriteString(string(event) + "\n")
}

//testClient method to get messages from the server.
func (client testClient) getMessage() string {
    message, err := client.reader.ReadString('\n')
    if err != nil {
        log.Fatalln(err)
    }
    return message
}

//testClient method to get events from the server.
func (client testClient) getEvent() events.ChatEvent {
    message, err := client.reader.ReadString('\n')
    if err != nil {
        log.Fatalln(err)
    }
    return events.ToChatEvent(message)
}

//Code to execute before running the unit tests.
func TestMain(m *testing.M) {
    chatServer = *NewServer(aleatoryPort())
    go chatServer.Serve()
    tClient = *newClient(aleatoryPort())
    flag.Parse()
    runTests := m.Run()
    os.Exit(runTests)
}

// Guest Constructor test.
func TestNewGuest(t *testing.T) {
    testSerial := randomSerial()
    guest := NewGuest(testSerial, nil)
    if guest.GetSerial() != testSerial {
        t.Error("TestNewGuest FAILED")
    }
    if guest.GetUsername() != "" {
        t.Error("TestNewGuest FAILED")
    }
}

// Tests if the guest's username can be modified.
func TestGuestSetGetUsername(t *testing.T) {
    guest := NewGuest(randomSerial(), nil)
    username := "TEST_USERNAME"
    guest.SetUsername(username)
    if guest.GetUsername() != username {
        t.Error("TestGuestSetGetUsername FAILED")
    }
}

// Tests if the guest is in a ChatRoom.
func TestGuestIsInChatroom(t *testing.T) {
    guest := NewGuest(randomSerial(), nil)
    chatRoom := NewChatRoom("TEST_CHATROOM")
    if guest.IsInChatRoom(*chatRoom) {
        t.Error("TestGuestIsInChatroom FAILED")
    }
    guest.JoinChatRoom(*chatRoom)
    if !guest.IsInChatRoom(*chatRoom) {
        t.Error("TestGuestIsInChatroom FAILED")
    }
}

// Tests if the guest can join a chatRoom.
func TestGuestJoinChatroom(t *testing.T) {
    guest := NewGuest(randomSerial(), nil)
    chatRoom := NewChatRoom("TEST_CHATROOM")
    guest.JoinChatRoom(*chatRoom)
    if !guest.IsInChatRoom(*chatRoom) {
        t.Error("TestGuestJoinChatroom FAILED")
    }
}

// Tests if the guest can leave a chatRoom.
func TestGuestLeaveChatroom(t *testing.T) {
    guest := NewGuest(randomSerial(), nil)
    chatRoom := NewChatRoom("TEST_CHATROOM")
    guest.JoinChatRoom(*chatRoom)
    guest.LeaveChatRoom(*chatRoom)
    if guest.IsInChatRoom(*chatRoom) {
        t.Error("TestGuestLeaveChatroom FAILED")
    }
}

func TestGuestEquals(t *testing.T) {
    testSerial := randomSerial()
    guest := NewGuest(testSerial, nil)
    guest.SetUsername("A")
    guestA := NewGuest(testSerial, nil)
    guestA.SetUsername("A")
    guestB := NewGuest(randomSerial(), nil)
    guestB.SetUsername("B")
    if guest.Equals(guestB) {
        t.Error("TestGuestEquals FAILED")
    } else if !guest.Equals(guestA) {
        t.Error("TestGuestEquals FAILED")
    }
}

// ChatRoom Constructor test.
func TestNewChatRoom(t *testing.T) {
    name := "TEST_NAME"
    chatRoom := NewChatRoom(name)
    if chatRoom.GetName() != name {
        t.Error("TestNewChatRoom FAILED")
    }
    if chatRoom.GetConnectionCount() != 0 {
        t.Error("TestNewChatRoom FAILED")
    }
}

// Tests if the ChatRoom's username can be modified.
func TestChatRoomSetGetName(t *testing.T) {
    nameA := "TEST_NAME_A"
    chatRoom := NewChatRoom(nameA)
    nameB := "TEST_NAME_B"
    chatRoom.SetName(nameB)
    if chatRoom.GetName() != nameB {
        t.Error("TestChatRoomSetGetUsername FAILED")
    }
}

// Tests if the ChatRoom can accept a guest.
func TestChatRoomAddGuest(t *testing.T) {
    chatRoom := NewChatRoom("TEST_CHATROOM")
    guest := NewGuest(randomSerial(), nil)
    chatRoom.AddGuest(*guest)
    if chatRoom.GetConnectionCount() != 1 {
        t.Error("TestChatRoomAddGuest FAILED")
    }
}

// Tests the ChatRoom connection counter.
func TestChatRoomGetConnectionCount(t *testing.T) {
    chatRoom := NewChatRoom("TEST_CHATROOM")
    guestA := NewGuest(randomSerial(), nil)
    guestB := NewGuest(randomSerial(), nil)
    chatRoom.AddGuest(*guestA)
    chatRoom.AddGuest(*guestB)
    if chatRoom.GetConnectionCount() != 2 {
        t.Error("TestChatRoomGetConnectionCount FAILED")
    }
}

// Tests if the ChatRoom can remove a guest.
func TestChatRoomRemoveGuest(t *testing.T) {
    chatRoom := NewChatRoom("TEST_CHATROOM")
    guestA := NewGuest(randomSerial(), nil)
    guestB := NewGuest(randomSerial(), nil)
    chatRoom.AddGuest(*guestA)
    chatRoom.AddGuest(*guestB)
    chatRoom.RemoveGuest(*guestA)
    if chatRoom.GetConnectionCount() != 1 {
        t.Error("TestChatRoomRemoveGuest FAILED")
    }
    chatRoom.RemoveGuest(*guestB)
    if chatRoom.GetConnectionCount() != 0 {
        t.Error("TestChatRoomRemoveGuest FAILED")
    }
}

func TestChatRoomEquals(t *testing.T) {
    chatRoom := NewChatRoom("TEST_CHATROOM")
    chatRoomA := NewChatRoom("TEST_CHATROOM")
    chatRoomB := NewChatRoom("TEST_CHATROOM_B")
    chatRoomC := NewChatRoom("TEST_CHATROOM_C")
    guest1 := NewGuest(randomSerial(), nil)
    guest2 := NewGuest(randomSerial(), nil)
    chatRoom.AddGuest(*guest1)
    chatRoomA.AddGuest(*guest1)
    chatRoomB.AddGuest(*guest2)
    chatRoomC.AddGuest(*guest1)
    if chatRoom.Equals(chatRoomB) {
        t.Error("TestChatRoomEquals FAILED")
    }else if chatRoom.Equals(chatRoomC) {
        t.Error("TestChatRoomEquals FAILED")
    } else if !chatRoom.Equals(chatRoomA) {
        t.Error("TestChatRoomEquals FAILED")
    }
}

// Tests if the server sends a message to a specific client.
func TestSendMessage(t *testing.T) {
    guestConn, client := net.Pipe()
    message := "TEST_MESSAGE"
    guest := *NewGuest(randomSerial(), guestConn)
    go func ()  {
        sendMessage(message, guest)
        guestConn.Close()
    }()
    receivedMessage, _ := bufio.NewReader(client).ReadString('\n')
    receivedMessage = strings.Trim(receivedMessage, "\n")
    client.Close()
    if receivedMessage != message {
        t.Error("TestSendMessage FAILED")
    }
}

// Tests if the server sends an event to a specific client.
func TestSendEvent(t *testing.T) {
    guestConn, client := net.Pipe()
    event := events.UNDEFINED
    guest := *NewGuest(randomSerial(), guestConn)
    go func ()  {
        sendEvent(event, guest)
        guestConn.Close()
    }()
    receivedMessage, _ := bufio.NewReader(client).ReadString('\n')
    receivedEvent := events.ToChatEvent(strings.Trim(receivedMessage, "\n"))
    client.Close()
    if receivedEvent != event {
        t.Error("TestSendEvent FAILED")
    }
}
