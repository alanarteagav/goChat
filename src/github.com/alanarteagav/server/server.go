package server

import (
    "bufio"
    "errors"
    "github.com/alanarteagav/events"
    "fmt"
    "log"
    "net"
    "strings"
    "strconv"
)

type Server struct {
    port                int
    guestsDictionary    map[string]Guest
}

func NewServer(port int) *Server {
    server := new(Server)
    server.port = port
    server.guestsDictionary = make(map[string]Guest)
    return server
}

func (server Server) GetPort() int {
    return server.port
}

// Auxiliar function which listens strings from a connection.
func listen(connection net.Conn) (string, error) {
    message, err := bufio.NewReader(connection).ReadString('\n')
    if err != nil {
        return "",  errors.New("Can't listen through connection")
    }
    message = strings.Trim(message, "\n")
    return message, nil
}

//Server method, sends a message to a guest.
func sendEvent(event events.ChatEvent, guest Guest) {
    guest.GetConnection().Write([]byte(event + "\n"))
}

// Server method, sends a message to a guest.
func sendMessage(message string, guest Guest) {
    guest.GetConnection().Write([]byte(message + "\n"))
}

// Auxiliar function which send messages to all the guests in the
// guest dictionary.
func (server Server) deliverMessage(message string) {
    for _ , guest := range server.guestsDictionary {
        if &guest != nil {
            sendMessage(message, guest)
        }
    }
}

func (server Server) listen(guest *Guest) (string, error) {
    message, err := bufio.NewReader(guest.GetConnection()).ReadString('\n')
    if err != nil {
        fmt.Println("[Client : " + guest.GetUsername() + " disconnected]")
        delete(server.guestsDictionary, guest.GetUsername())
        return "", errors.New("Client out")
    }
    message = strings.Trim(message, "\n")
    return message, nil
}

func (server Server) handleConnection(guest *Guest)  {
    for {
        message, err := server.listen(guest)
        if err != nil{
            return
        }
        event := events.ToChatEvent(message)
        switch event {
            case events.MESSAGE:
                messageIn, err := server.listen(guest)
                if err != nil{
                    return
                }
                server.deliverMessage(messageIn)
                fmt.Println("[ MESSAGE " + messageIn + " ]")
            case events.UNDEFINED:
                fmt.Println("[UNDEFINED EVENT]")
            case events.ERROR:
                fmt.Println("[ERROR!]")
        }
    }
}

func (server Server) Serve()  {
    listener, err := net.Listen("tcp",
                                "localhost:" + strconv.Itoa(server.port))
    if err != nil {
        log.Fatalln(err)
    }
    for {
        connection, err := listener.Accept()
        if err != nil {
            log.Fatalln(err)
        }
        connection.Write([]byte(events.LOG_IN + "\n"))
        username, err := listen(connection)
        if err != nil {
            log.Fatalln(err)
        }
        fmt.Println("[ " + username + " has logged in ]")
        guest := NewGuest(username, connection)
        server.guestsDictionary[username] = *guest
        go server.handleConnection(guest)
    }
}
