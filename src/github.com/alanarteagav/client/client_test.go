package client

import (
    "testing"
    "os"
    "flag"
    "net"
    "log"
    "bufio"
    "strconv"
)

func aleatoryPort() int {
    return 4000
}

var chatServer testServer
var stringChannel chan string
var port int

// A small server to test our client
type testServer struct {
    port       int
    guests     []net.Conn
    connection net.Conn
    listener   net.Listener
}

// Creates a new testServer
func newTestServer(port int) *testServer {
    server := new(testServer)
    listener, err := net.Listen("tcp", "localhost:" + strconv.Itoa(port))
    if err != nil {
        log.Fatalln(err)
    }
    server.guests = make([]net.Conn, 1)
    server.listener = listener
    return server
}

// Serves
func (ts testServer) serve(stringChannel chan string){
    for {
        connection, err := ts.listener.Accept()
        if err != nil {
            log.Fatalln(err)
        }
        ts.guests = append(ts.guests, connection)
        go ts.handle(connection)
    }
}

// Delivers the messages
func (ts testServer) handle(connection net.Conn) {
    for {
        message, err := bufio.NewReader(connection).ReadString('\n')
        if err != nil {
            log.Fatalln(err)
        }
        for _, guest := range ts.guests {
            if guest != nil {
                guest.Write([]byte(message + "\n"))
            }
        }
    }
}

// Code to execute before running the unit tests.
func TestMain(m *testing.M) {
    port = aleatoryPort()
    chatServer := *newTestServer(port)
    stringChannel = make(chan string)
    go chatServer.serve(stringChannel)
    flag.Parse()
    runTests := m.Run()
    os.Exit(runTests)
}

// Tests client constructor.
func TestNewClient(t *testing.T) {
    username := "NAME"
    client := NewClient(username)
    if username != client.GetUsername(){
        t.Error("TestNewClient FAILED")
    }
}

// Tests if the client's username can be modified.
func TestSetGetUsername(t *testing.T) {
    username := "NAME"
    client := NewClient(username)
    newUsername := "NEW_NAME"
    client.SetUsername(newUsername)
    if newUsername != client.GetUsername(){
        t.Error("TestSetGetUsername FAILED")
    }
}
