package server

import (
    "bufio"
    "flag"
    "github.com/alanarteagav/events"
    "log"
    "math/rand"
    "net"
    "strings"
    "time"
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
    username := "TEST_USERNAME"
    guest := NewGuest(username, nil)

    if guest.GetUsername() != username {
        t.Error("TestNewGuest FAILED")
    }
}

// Tests if the guest's username can be modified.
func TestGuestSetGetUsername(t *testing.T) {
    usernameA := "TEST_USERNAME_A"
    guest := NewGuest(usernameA, nil)
    usernameB := "TEST_USERNAME_B"
    guest.SetUsername(usernameB)
    if guest.GetUsername() != usernameB {
        t.Error("TestGuestSetGetUsername FAILED")
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
    guestA := NewGuest("GUEST_A", nil)
    chatRoom.AddGuest(*guestA)
    if chatRoom.GetConnectionCount() != 1 {
        t.Error("TestChatRoomAddGuest FAILED")
    }
}

// Tests the ChatRoom connection counter.
func TestChatRoomGetConnectionCount(t *testing.T) {
    chatRoom := NewChatRoom("TEST_CHATROOM")
    guestA := NewGuest("GUEST_A", nil)
    guestB := NewGuest("GUEST_B", nil)
    chatRoom.AddGuest(*guestA)
    chatRoom.AddGuest(*guestB)
    if chatRoom.GetConnectionCount() != 2 {
        t.Error("TestChatRoomGetConnectionCount FAILED")
    }
}

// Tests if the ChatRoom can remove a guest.
func TestChatRoomRemoveGuest(t *testing.T) {
    chatRoom := NewChatRoom("TEST_CHATROOM")
    guestA := NewGuest("GUEST_A", nil)
    guestB := NewGuest("GUEST_B", nil)
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

// Tests if the server can accept a client.
func TestLogIn(t *testing.T) {
    logInSignal := strings.Trim(tClient.getMessage(), "\n")
    tClient.sendMessage("TEST_CLIENT" + "\n")
    if events.ToChatEvent(logInSignal) != events.LOG_IN {
        t.Error("TestLogIn FAILED")
    }
}

// Tests if the server sends a message to a specific client.
func TestSendMessage(t *testing.T) {
    guestConn, client := net.Pipe()
    message := "TEST_MESSAGE"
    guest := *NewGuest("TEST_GUEST", guestConn)
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
    guest := *NewGuest("TEST_GUEST", guestConn)
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

// Tests if the server delivers messages to the clients.
func TestDeliverMessage(t *testing.T) {
    tClient2 := *newClient(aleatoryPort())
    tClient2.sendMessage("TEST_CLIENT_2" + "\n")
    tClient2.getMessage()

    message := "TEST_DELIVER"
    <-time.After(1 * time.Second)
    tClient2.sendMessage("MESSAGE")
    <-time.After(1 * time.Second)
    tClient2.sendMessage(message)

    receivedMessage1 := strings.Trim(tClient.getMessage(), "\n")
    receivedMessage2 := strings.Trim(tClient2.getMessage(), "\n")

    if receivedMessage1 != receivedMessage2 {
        t.Error("TestDeliverMessage FAILED")
    } else if  receivedMessage1 != message {
        t.Error("TestDeliverMessage FAILED")
    }
}
