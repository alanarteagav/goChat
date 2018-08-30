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

func TestLogIn(t *testing.T) {
    logInSignal := strings.Trim(tClient.getMessage(), "\n")
    tClient.sendMessage("TEST_CLIENT" + "\n")
    if events.ToChatEvent(logInSignal) != events.LOG_IN {
        t.Error("TestLogIn FAILED")
    }
}

func TestSendMessage(t *testing.T) {
    message := "TEST_MESSAGE"
    <-time.After(1 * time.Second)
    tClient.sendMessage("MESSAGE")
    <-time.After(1 * time.Second)
    tClient.sendMessage(message)
    receivedMessage := strings.Trim(tClient.getMessage(), "\n")
    if receivedMessage != message {
        t.Error("TestSendMessage FAILED")
    }
}

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
