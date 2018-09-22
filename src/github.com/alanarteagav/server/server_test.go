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
    "fmt"
)


var chatServer Server
var clientA testClient
var clientB testClient
var ownerGuest Guest

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
    ownerGuest = *NewGuest(randomSerial(), nil)
    clientB = *newClient(aleatoryPort())
    clientB.sendMessage("IDENTIFY Kenobi")
    _ = strings.Trim(clientB.getMessage(), "\n")
    flag.Parse()
    runTests := m.Run()
    os.Exit(runTests)
    fmt.Println("hi")
}


/* TESTS FOR THE GUEST STRUCT */

// Guest Constructor test.
func TestNewGuest(t *testing.T) {
    testSerial := randomSerial()
    guest := NewGuest(testSerial, nil)
    if guest.getSerial() != testSerial {
        t.Error("TestNewGuest FAILED")
    }
    if guest.getUsername() != "" {
        t.Error("TestNewGuest FAILED")
    }
}

// Tests if the guest's username can be modified.
func TestGuestSetgetUsername(t *testing.T) {
    guest := NewGuest(randomSerial(), nil)
    username := "TEST_USERNAME"
    guest.setUsername(username)
    if guest.getUsername() != username {
        t.Error("TestGuestSetgetUsername FAILED")
    }
}

func TestGuestSetgetStatus(t *testing.T) {
    guest := NewGuest(randomSerial(), nil)
    guest.setStatus(BUSY)
    if guest.getStatus() != BUSY {
        t.Error("TestGuestSetgetStatus FAILED")
    }
    guest.setStatus(AWAY)
    if guest.getStatus() != AWAY {
        t.Error("TestGuestSetgetStatus FAILED")
    }
    guest.setStatus(ACTIVE)
    if guest.getStatus() != ACTIVE {
        t.Error("TestGuestSetgetStatus FAILED")
    }
}

// Tests if the guest is in a ChatRoom.
func TestGuestisInChatRoom(t *testing.T) {
    guest := NewGuest(randomSerial(), nil)
    chatRoom := NewChatRoom(ownerGuest, "TEST_CHATROOM")
    if guest.isInChatRoom(*chatRoom) {
        t.Error("TestGuestisInChatRoom FAILED")
    }
    guest.joinChatRoom(*chatRoom)
    if !guest.isInChatRoom(*chatRoom) {
        t.Error("TestGuestisInChatRoom FAILED")
    }
}

// Tests if the guest can join a chatRoom.
func TestGuestjoinChatRoom(t *testing.T) {
    guest := NewGuest(randomSerial(), nil)
    chatRoom := NewChatRoom(ownerGuest, "TEST_CHATROOM")
    guest.joinChatRoom(*chatRoom)
    if !guest.isInChatRoom(*chatRoom) {
        t.Error("TestGuestjoinChatRoom FAILED")
    }
}

// Tests if the guest can leave a chatRoom.
func TestGuestleaveChatRoom(t *testing.T) {
    guest := NewGuest(randomSerial(), nil)
    chatRoom := NewChatRoom(ownerGuest, "TEST_CHATROOM")
    guest.joinChatRoom(*chatRoom)
    guest.leaveChatRoom(*chatRoom)
    if guest.isInChatRoom(*chatRoom) {
        t.Error("TestGuestleaveChatRoom FAILED")
    }
}

func TestGuestequals(t *testing.T) {
    testSerial := randomSerial()
    guest := NewGuest(testSerial, nil)
    guest.setUsername("A")
    guestA := NewGuest(testSerial, nil)
    guestA.setUsername("A")
    guestB := NewGuest(randomSerial(), nil)
    guestB.setUsername("B")
    if guest.equals(guestB) {
        t.Error("TestGuestequals FAILED")
    } else if !guest.equals(guestA) {
        t.Error("TestGuestequals FAILED")
    }
}

/* TESTS FOR THE CHATROOM STRUCT */

// ChatRoom Constructor test.
func TestNewChatRoom(t *testing.T) {
    name := "TEST_NAME"
    guest := NewGuest(randomSerial(), nil)
    chatRoom := NewChatRoom(*guest, name)
    if !guest.equals(chatRoom.getOwner()) {
        t.Error("TestNewChatRoom FAILED")
    } else if chatRoom.getName() != name {
        t.Error("TestNewChatRoom FAILED")
    } else if !chatRoom.hosts(guest){
        t.Error("TestNewChatRoom FAILED")
    }else if chatRoom.getConnectionCount() != 1 {
        t.Error("TestNewChatRoom FAILED")
    }
}

// Tests if the ChatRoom's username can be modified.
func TestChatRoomSetgetName(t *testing.T) {
    nameA := "TEST_NAME_A"
    chatRoom := NewChatRoom(ownerGuest, nameA)
    nameB := "TEST_NAME_B"
    chatRoom.setName(nameB)
    if chatRoom.getName() != nameB {
        t.Error("TestChatRoomSetgetUsername FAILED")
    }
}

// Tests if the ChatRoom can accept a guest.
func TestChatRoomaddInvitedGuest(t *testing.T) {
    chatRoom := NewChatRoom(ownerGuest, "TEST_CHATROOM")
    invitedGuest := NewGuest(randomSerial(), nil)
    notOwnerGuest := NewGuest(randomSerial(), nil)
    if chatRoom.wasInvited(invitedGuest) {
        t.Error("TestChatRoomaddInvitedGuest FAILED")
    }
    chatRoom.addInvitedGuest(*notOwnerGuest, invitedGuest)
    if chatRoom.wasInvited(invitedGuest) {
        t.Error("TestChatRoomaddInvitedGuest FAILED")
    }
    chatRoom.addInvitedGuest(ownerGuest, invitedGuest)
    if !chatRoom.wasInvited(invitedGuest) {
        t.Error("TestChatRoomaddInvitedGuest FAILED")
    }
}

// Tests if the ChatRoom can accept a guest.
func TestChatRoomaddGuest(t *testing.T) {
    chatRoom := NewChatRoom(ownerGuest, "TEST_CHATROOM")
    invitedGuest := NewGuest(randomSerial(), nil)
    chatRoom.addGuest(*invitedGuest)
    if chatRoom.hosts(invitedGuest) || chatRoom.getConnectionCount() != 1 {
        t.Error("TestChatRoomaddGuest FAILED")
    }
    chatRoom.addInvitedGuest(ownerGuest, invitedGuest)
    chatRoom.addGuest(*invitedGuest)
    if !chatRoom.hosts(invitedGuest) || chatRoom.getConnectionCount() != 2 {
        t.Error("TestChatRoomaddGuest FAILED")
    }
}

// Tests the ChatRoom connection counter.
func TestChatRoomgetConnectionCount(t *testing.T) {
    chatRoom := NewChatRoom(ownerGuest, "TEST_CHATROOM")
    guestA := NewGuest(randomSerial(), nil)
    guestB := NewGuest(randomSerial(), nil)
    chatRoom.addInvitedGuest(ownerGuest, guestA)
    chatRoom.addInvitedGuest(ownerGuest, guestB)
    chatRoom.addGuest(*guestA)
    chatRoom.addGuest(*guestB)
    if chatRoom.getConnectionCount() != 3 {
        t.Error("TestChatRoomgetConnectionCount FAILED")
    }
}

// Tests if the ChatRoom can remove a guest.
func TestChatRoomremoveGuest(t *testing.T) {
    chatRoom := NewChatRoom(ownerGuest, "TEST_CHATROOM")
    guestA := NewGuest(randomSerial(), nil)
    guestB := NewGuest(randomSerial(), nil)
    chatRoom.addInvitedGuest(ownerGuest, guestA)
    chatRoom.addInvitedGuest(ownerGuest, guestB)
    chatRoom.addGuest(*guestA)
    chatRoom.addGuest(*guestB)
    chatRoom.removeGuest(*guestA)
    if chatRoom.getConnectionCount() != 2 {
        t.Error("TestChatRoomremoveGuest FAILED")
    }
    chatRoom.removeGuest(*guestB)
    if chatRoom.getConnectionCount() != 1 {
        t.Error("TestChatRoomremoveGuest FAILED")
    }
}


/* TESTS FOR THE SERVER STRUCT */

// Tests if the server sends a message to a specific client.
func TestSend(t *testing.T) {
    guestConn, client := net.Pipe()
    message := "TEST_MESSAGE"
    guest := *NewGuest(randomSerial(), guestConn)
    go func ()  {
        send(message, guest)
        guestConn.Close()
    }()
    receivedMessage, _ := bufio.NewReader(client).ReadString('\n')
    receivedMessage = strings.Trim(receivedMessage, "\n")
    client.Close()
    if receivedMessage != message {
        t.Error("TestSend FAILED")
    }
}

// Tests if the server can accept a client.
func TestConnect(t *testing.T) {
    clientA = *newClient(aleatoryPort())
    clientA.sendMessage("TEST_MESSAGE")
    receivedMessage := strings.Trim(clientA.getMessage(), "\n")
    if receivedMessage != "...INVALID MESSAGE" {
        t.Error("TestConnect FAILED")
    }
}

func TestIDENTIFY(t *testing.T) {
    clientA.sendMessage("IDENTIFY")
    receivedMessage := strings.Trim(clientA.getMessage(), "\n")
    if receivedMessage != "...INVALID MESSAGE" {
        t.Error("TestIDENTIFY FAILED")
    }
    clientA.sendMessage("IDENTIFY Skywalker")
    receivedMessage = strings.Trim(clientA.getMessage(), "\n")
    if receivedMessage != "...SUCCESFUL IDENTIFICATION" {
        t.Error("TestIDENTIFY FAILED")
    }
    clientA.sendMessage("IDENTIFY Kenobi")
    receivedMessage = strings.Trim(clientA.getMessage(), "\n")
    if receivedMessage != "...USERNAME NOT AVAILABLE" {
        t.Error("TestIDENTIFY FAILED")
    }
}

func TestPUBLICMESSAGE(t *testing.T) {
    client := *newClient(aleatoryPort())
    client.sendMessage("PUBLICMESSAGE")
    receivedMessage := strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...INVALID MESSAGE" {
        t.Error("TestPUBLICMESSAGE FAILED")
    }
    client.sendMessage("PUBLICMESSAGE message")
    receivedMessage = strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...MUST IDENTIFY FIRST" {
        t.Error("TestPUBLICMESSAGE FAILED")
    }
    clientA.sendMessage("PUBLICMESSAGE message")
    receivedMessage = strings.Trim(clientB.getMessage(), "\n")
    if receivedMessage != "...PUBLIC-Skywalker: message" {
        t.Error("TestPUBLICMESSAGE FAILED")
    }
}

func TestMESSAGE(t *testing.T) {
    client := *newClient(aleatoryPort())
    client.sendMessage("MESSAGE")
    receivedMessage := strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...INVALID MESSAGE" {
        t.Error("TestMESSAGE FAILED")
    }
    client.sendMessage("MESSAGE Skywalker message")
    receivedMessage = strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...MUST IDENTIFY FIRST" {
        t.Error("TestMESSAGE FAILED")
    }
    clientA.sendMessage("MESSAGE QuiGon message")
    receivedMessage = strings.Trim(clientA.getMessage(), "\n")
    if receivedMessage != "...USER QuiGon NOT FOUND" {
        t.Error("TestMESSAGE FAILED")
    }
    clientB.sendMessage("MESSAGE Skywalker message")
    receivedMessageA := strings.Trim(clientA.getMessage(), "\n")
    receivedMessageB := strings.Trim(clientB.getMessage(), "\n")
    if receivedMessageB != "...MESSAGE SENT" {
        t.Error("TestMESSAGE FAILED")
    } else if receivedMessageA != "Kenobi: message" {
        t.Error("TestMESSAGE FAILED")
    }
}

func TestUSERS(t *testing.T) {
    client := *newClient(aleatoryPort())
    client.sendMessage("USERS")
    receivedMessage := strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...MUST IDENTIFY FIRST" {
        t.Error("TestUSERS FAILED")
    }
    clientA.sendMessage("USERS")
    receivedMessage = strings.Trim(clientA.getMessage(), "\n")
    receivedMessageArray := strings.Split(receivedMessage, " ")
    if len(receivedMessageArray) != 2 {
        t.Error("TestUSERS FAILED")
    }
    entryA := receivedMessageArray[0]
    entryB := receivedMessageArray[1]
    if entryA != "Kenobi" {
        if entryA != "Skywalker" {
            t.Error("TestUSERS FAILED")
        }
    } else if entryB != "Kenobi" {
        if entryB != "Skywalker" {
            t.Error("TestUSERS FAILED")
        }
    }
}

func TestSTATUS(t *testing.T) {
    client := *newClient(aleatoryPort())
    client.sendMessage("STATUS")
    receivedMessage := strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...INVALID MESSAGE" {
        t.Error("TestSTATUS FAILED")
    }

    client.sendMessage("STATUS BUSY")
    receivedMessage = strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...MUST IDENTIFY FIRST" {
        t.Error("TestSTATUS FAILED")
    }

    clientA.sendMessage("STATUS EVENT")
    receivedMessage = strings.Trim(clientA.getMessage(), "\n")
    if receivedMessage != "...INVALID STATUS" {
        t.Error("TestSTATUS FAILED")
    }

    clientA.sendMessage("STATUS ACTIVE")
    receivedMessage = strings.Trim(clientA.getMessage(), "\n")
    if receivedMessage != "Skywalker ACTIVE" {
        t.Error("TestSTATUS FAILED")
    }
    clientA.sendMessage("STATUS BUSY")
    receivedMessage = strings.Trim(clientA.getMessage(), "\n")
    if receivedMessage != "Skywalker BUSY" {
        t.Error("TestSTATUS FAILED")
    }
    clientA.sendMessage("STATUS AWAY")
    receivedMessage = strings.Trim(clientA.getMessage(), "\n")
    if receivedMessage != "Skywalker AWAY" {
        t.Error("TestSTATUS FAILED")
    }
}

func TestCREATEROOM(t *testing.T) {
    client := *newClient(aleatoryPort())
    client.sendMessage("CREATEROOM")
    receivedMessage := strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...INVALID MESSAGE" {
        t.Error("TestCREATEROOM FAILED")
    }
    client.sendMessage("CREATEROOM MosEisley")
    receivedMessage = strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...MUST IDENTIFY FIRST" {
        t.Error("TestCREATEROOM FAILED")
    }
    clientA.sendMessage("CREATEROOM MosEisley")
    receivedMessage = strings.Trim(clientA.getMessage(), "\n")
    if receivedMessage != "...ROOM CREATED" {
        t.Error("TestCREATEROOM FAILED")
    }
}

func TestINVITE(t *testing.T) {
    client := *newClient(aleatoryPort())
    client.sendMessage("INVITE")
    receivedMessage := strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...INVALID MESSAGE" {
        t.Error("TestINVITE FAILED")
    }
    client.sendMessage("INVITE MosEisley Skywalker")
    receivedMessage = strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...MUST IDENTIFY FIRST" {
        t.Error("TestINVITE FAILED")
    }
    client.sendMessage("IDENTIFY Yoda")
    _ = strings.Trim(client.getMessage(), "\n")
    client.sendMessage("INVITE MosEisley Skywalker")
    receivedMessage = strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...YOU ARE NOT THE OWNER OF THE ROOM" {
        t.Error("TestINVITE FAILED")
    }
    clientA.sendMessage("INVITE MosEisley Yoda")
    receivedMessageA := strings.Trim(clientA.getMessage(), "\n")
    receivedMessage = strings.Trim(client.getMessage(), "\n")
    if receivedMessageA != "...INVITATION SENT TO Yoda" {
        t.Error("TestINVITE FAILED")
    } else if receivedMessage != "...INVITATION TO JOIN MosEisley ROOM BY Skywalker" {
        t.Error("TestINVITE FAILED")
    }
}

func TestJOINROOM(t *testing.T) {
    client := *newClient(aleatoryPort())
    client.sendMessage("JOINROOM")
    receivedMessage := strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...INVALID MESSAGE" {
        t.Error("TestJOINROOM FAILED")
    }

    client.sendMessage("JOINROOM MosEisley")
    receivedMessage = strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...MUST IDENTIFY FIRST" {
        t.Error("TestJOINROOM FAILED")
    }

    clientB.sendMessage("JOINROOM MosEisley")
    receivedMessage = strings.Trim(clientB.getMessage(), "\n")
    if receivedMessage != "...YOU ARE NOT INVITED TO ROOM MosEisley" {
        t.Error("TestJOINROOM FAILED")
    }

    clientA.sendMessage("INVITE MosEisley Kenobi")
    _ = strings.Trim(clientA.getMessage(), "\n")
    _ = strings.Trim(clientB.getMessage(), "\n")
    clientB.sendMessage("JOINROOM MosEisley")
    receivedMessage = strings.Trim(clientB.getMessage(), "\n")
    if receivedMessage != "...SUCCESFULLY JOINED TO ROOM" {
        t.Error("TestJOINROOM FAILED")
    }
}

func TestROOMESSAGE(t *testing.T) {
    client := *newClient(aleatoryPort())
    client.sendMessage("ROOMESSAGE")
    receivedMessage := strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...INVALID MESSAGE" {
        t.Error("TestROOMESSAGE FAILED")
    }
    client.sendMessage("ROOMESSAGE MosEisley message")
    receivedMessage = strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...MUST IDENTIFY FIRST" {
        t.Error("TestROOMESSAGE FAILED")
    }
    clientA.sendMessage("ROOMESSAGE DeathStar message")
    receivedMessage = strings.Trim(clientA.getMessage(), "\n")
    if receivedMessage != "...ROOM NOT EXIST" {
        t.Error("TestROOMESSAGE FAILED")
    }
    client.sendMessage("IDENTIFY Palpatine")
    _ = strings.Trim(client.getMessage(), "\n")
    client.sendMessage("ROOMESSAGE MosEisley message")
    receivedMessage = strings.Trim(client.getMessage(), "\n")
    if receivedMessage != "...YOU ARE NOT PART OF THE ROOM" {
        t.Error("TestROOMESSAGE FAILED")
    }
    clientA.sendMessage("ROOMESSAGE MosEisley message")
    receivedMessageA := strings.Trim(clientA.getMessage(), "\n")
    receivedMessageB := strings.Trim(clientB.getMessage(), "\n")
    if receivedMessageA != "...MESSAGE SENT" {
        t.Error("TestROOMESSAGE FAILED")
    } else if receivedMessageB != "...MosEisley-Skywalker: message" {
        t.Error("TestROOMESSAGE FAILED")
    }
}
