package client

import (
    "testing"
    "time"
    "net"
    "log"
    "bufio"
    "strconv"
    "strings"
)

//A small server to test our client
type testServer struct { port int }

func (ts testServer) deliverMessage(connection net.Conn) {
    message := "beep"
    connection.Write([]byte(message + "\n"))
}

func (ts testServer) serve(stringChannel chan string){
    listener, err := net.Listen("tcp", ":" + strconv.Itoa(ts.port))
    if err != nil {
        log.Println("Error: can't listen through port")
    }
    for {
        connection, err := listener.Accept()
        if err != nil {
            log.Fatalln(err)
        }
        username, _ := bufio.NewReader(connection).ReadString('\n')

        go ts.deliverMessage(connection)

        stringChannel <- username
    }
}

server := testServer{ 3000 }

func (ts testServer) sendMessage()  {

}


func TestSetUsername(t *testing.T) {
    t.Error("TestSetUsername FAILED")
}

//Tests if the client recieves messages from the server
//The test fails if it takes more than two seconds
func TestListen(t *testing.T) {

    server := testServer{ 3000 }
    client := Client{ "clientName" }
    stringChannel := make(chan string)
    messageChannel := make(chan string)

    go server.serve(stringChannel)

    select {
    case <- stringChannel:
        client.Connect(strconv.Itoa(3000))
        client.Listen(messageChannel)
        message := <- messageChannel
        message = strings.Trim(message, "\n")
        if message != "beep" {
            t.Error("TestListen FAILED")
        }
    case <-time.After(2 * time.Second):
        t.Fatal("timeout")
    }
}

func TestRead(t *testing.T) {
    t.Error("TestRead FAILED")
}

func TestWrite(t *testing.T) {
    t.Error("TestWrite FAILED")
}

//Tests when the client connects to the server
//The test fails if it takes more than two seconds
func TestLogIn(t *testing.T)  {

    server := testServer{ 3000 }
    client := Client{ "clientName" }
    stringChannel := make(chan string)
    go server.serve(stringChannel)
    select {
    case <- stringChannel:
        client.Connect(strconv.Itoa(3000))
        id := <- stringChannel
        id = strings.Trim(id, "\n")
        if id != "clientName" {
            t.Error("TestLogIn FAILED")
        }
    case <-time.After(2 * time.Second):
        t.Fatal("timeout")
    }
}

func TestLogOut(t *testing.T) {
    t.Error("TestLogOut FAILED")
}

func TestAskForChatroom(t *testing.T) {
    t.Error("TestAskForChatroom FAILED")
}
