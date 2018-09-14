package server

import (
    "bufio"
    "errors"
    "github.com/alanarteagav/events"
    "fmt"
    "log"
    "math/rand"
    "net"
    "strings"
    "strconv"
)

func randomSerial() int {
    return rand.Int()
}

// Server struct.
// Defines the server's port, and a dictionary of guests.
type Server struct {
    port        int
    guestsById  map[int]Guest
}

// Server constructor.
func NewServer(port int) *Server {
    server := new(Server)
    server.port = port
    server.guestsById = make(map[int]Guest)
    return server
}

// Returns server's port.
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

// Server method, sends a message to a guest.
func send(message string, guest Guest) {
    guest.GetConnection().Write([]byte(message + "\n"))
}

// Auxiliar function which send messages to all the guests in the
// guest dictionary.
func (server Server) deliver(message string) {
    for _ , guest := range server.guestsById {
        if &guest != nil {
            send(message, guest)
        }
    }
}

// Auxiliar method which listens strings from a guest.
func (server Server) listen(guest *Guest) (string, error) {
    message, err := bufio.NewReader(guest.GetConnection()).ReadString('\n')
    if err != nil {
        fmt.Println("[Client : " + guest.GetUsername() + " disconnected]")
        delete(server.guestsById, guest.GetSerial())
        return "", errors.New("Client out")
    }
    message = strings.Trim(message, "\n")
    return message, nil
}

// Handles a particular guest connection.
func (server Server) handleConnection(guest *Guest)  {
    for {
        message, err := server.listen(guest)
        if err != nil{
            return
        }
        event := events.ToChatEvent(message)
        switch event {
        default :
            send("NOT A VALID EVENT", *guest)
        }
    }
}

// Serves.
func (server Server) Serve()  {
    listener, err := net.Listen("tcp", ":" + strconv.Itoa(server.port))
    if err != nil {
        log.Fatalln(err)
    }
    for {
        connection, err := listener.Accept()
        if err != nil {
            log.Fatalln(err)
        }
        guestSerial := randomSerial()
        for {
            if _, ok := server.guestsById[guestSerial]; ok {
                guestSerial = randomSerial()
            } else {
                break
            }
        }
        guest := NewGuest(guestSerial, connection)
        server.guestsById[guestSerial] = *guest
        serialString := strconv.Itoa(guestSerial)
        fmt.Println("[ NEW CLIENT WITH ID " + serialString + " HAS LOGGED IN ]")
        go server.handleConnection(guest)
    }
}
