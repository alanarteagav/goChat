/* Server package.

 It provides a single exported method "Serve", and a constructor, that
 allows to set up a server with a given port.

 The package consists in three files related with the following
 structs:

 - Server (main struct)

 - Guest (auxiliar struct)

 - ChatRoom (auxiliar struct)
*/
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


// Auxiliar function that returns an aleatory integer to temoparily
// identify a new guest.
func randomSerial() int {
    return rand.Int()
}

// Server struct.
// Defines the server's port, and hash tables for identified and
// unidentified guests.
type Server struct {
    port             int
    guestsById       map[int]Guest
    guestsByUsername map[string]Guest
    chatRooms        map[string]ChatRoom
}

// Server constructor.
func NewServer(port int) *Server {
    server := new(Server)
    server.port = port
    server.guestsById = make(map[int]Guest)
    server.guestsByUsername = make(map[string]Guest)
    server.chatRooms = make(map[string]ChatRoom)
    return server
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
    guest.getConnection().Write([]byte(message + "\n"))
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

// Auxiliar function which send messages to all identified users.
func (server Server) deliverToUsers(message string, sender Guest) {
    for _ , guest := range server.guestsByUsername {
        if &guest != nil && guest.getUsername() != sender.getUsername() {
            send(message, guest)
        }
    }
}

// Auxiliar method which listens strings from a guest.
func (server Server) listen(guest *Guest) (string, error) {
    message, err := bufio.NewReader(guest.getConnection()).ReadString('\n')
    if err != nil {
        serialString := strconv.Itoa(guest.getSerial())
        fmt.Println("[Client : " + serialString + " disconnected]")
        delete(server.guestsByUsername, guest.getUsername())
        delete(server.guestsById, guest.getSerial())
        return "", errors.New("Client out")
    }
    message = strings.Trim(message, "\n")
    return message, nil
}

// Handles a particular guest connection.
func (server Server) handleConnection(guest *Guest)  {
    for {
        message, err := server.listen(guest)
        if err != nil {
            return
        }
        stringArray := strings.Split(message, " ")
        var event string
        event = stringArray[0]
        switch event {
        case string(events.IDENTIFY):
            if len(stringArray) != 2 {
                send(string(events.INVALID), *guest)
            } else {
                username := stringArray[1]
                _, usernameAlreadyExists := server.guestsByUsername[username]
                if usernameAlreadyExists {
                    send("...USERNAME NOT AVAILABLE", *guest)
                } else {
                    delete(server.guestsByUsername, guest.getUsername())
                    guest.setUsername(username)
                    server.guestsByUsername[username] = *guest
                    send("...SUCCESFUL IDENTIFICATION", *guest)
                }
            }
        case string(events.PUBLICMESSAGE):
            if len(stringArray) < 2 {
                send(string(events.INVALID), *guest)
            } else if !guest.isIdentified() {
                send(string(events.IDENTIFY_ERROR), *guest)
            } else {
                stringToSend :=
                    strings.TrimPrefix(message, "PUBLICMESSAGE" + " ")
                stringToSend = "...PUBLIC-" + guest.getUsername() +
                               ": " + stringToSend
                server.deliverToUsers(stringToSend, *guest)
            }
        case string(events.MESSAGE):
            if len(stringArray) < 3 {
                send(string(events.INVALID), *guest)
            } else if !guest.isIdentified() {
                send(string(events.IDENTIFY_ERROR), *guest)
            } else {
                username := stringArray[1]
                if guestInHash, ok := server.guestsByUsername[username]; ok{
                    stringToSend :=
                    strings.TrimPrefix(
                        message, "MESSAGE" + " " + username + " ")
                    stringToSend = guest.getUsername() + ": " + stringToSend
                    send(stringToSend, guestInHash)
                    send("...MESSAGE SENT", *guest)
                } else {
                    send("...USER " + username + " NOT FOUND", *guest)
                }
            }
        case string(events.USERS):
            if len(stringArray) != 1 {
                send(string(events.INVALID), *guest)
            } else if !guest.isIdentified() {
                send(string(events.IDENTIFY_ERROR), *guest)
            } else {
                sendString := ""
                for username, _ := range server.guestsByUsername {
                    sendString += username + " "
                }
                sendString = strings.TrimSpace(sendString)
                send(sendString, *guest)
            }
        case string(events.STATUS):
            if len(stringArray) != 2 {
                send(string(events.INVALID), *guest)
            } else if !guest.isIdentified() {
                send(string(events.IDENTIFY_ERROR), *guest)
            } else {
                status := toUserStatus(stringArray[1])
                if status != UNDEFINED {
                    guest.setStatus(status)
                    send(guest.getUsername() +  " " + string(guest.getStatus()),
                        *guest)
                } else {
                    send("...INVALID STATUS", *guest)
                }
            }
        case string(events.CREATEROOM):
            if len(stringArray) != 2 {
                send(string(events.INVALID), *guest)
            } else if !guest.isIdentified() {
                send(string(events.IDENTIFY_ERROR), *guest)
            } else {
                chatRoomName := stringArray[1]
                chatRoom := *NewChatRoom(*guest, chatRoomName)
                server.chatRooms[chatRoomName] = chatRoom
                send("...ROOM CREATED", *guest)
            }
        case string(events.INVITE):
            if len(stringArray) < 3 {
                send(string(events.INVALID), *guest)
            } else if !guest.isIdentified() {
                send(string(events.IDENTIFY_ERROR), *guest)
            } else {
                chatRoomName := stringArray[1]
                if chatRoom, ok := server.chatRooms[chatRoomName]; ok{
                    stringArray = stringArray[2:]
                    for _, username := range stringArray {
                        if guestInHash, ok := server.guestsByUsername[username]; ok{
                            if chatRoom.addInvitedGuest(*guest, &guestInHash) {
                                send("...INVITATION SENT TO " + username, *guest)
                                send("...INVITATION TO JOIN " + chatRoomName +
                                    " ROOM BY " + guest.getUsername(),
                                    guestInHash)
                            } else {
                                send("...YOU ARE NOT THE OWNER OF THE ROOM",
                                    *guest)
                            }
                        }
                    }
                } else {
                    send("...ROOM NOT EXISTS",
                        *guest)
                }
            }
        case string(events.JOINROOM):
            if len(stringArray) != 2 {
                send(string(events.INVALID), *guest)
            } else if !guest.isIdentified() {
                send(string(events.IDENTIFY_ERROR), *guest)
            } else {
                chatRoomName := stringArray[1]
                if chatRoom, ok := server.chatRooms[chatRoomName]; ok{
                    if chatRoom.wasInvited(guest){
                       chatRoom.addGuest(*guest)
                       send("...SUCCESFULLY JOINED TO ROOM", *guest)
                    } else {
                        send("...YOU ARE NOT INVITED TO ROOM " + chatRoomName,
                            *guest)
                    }
                } else {
                    send("Sorry, " + chatRoomName + " doesn't exist",
                            *guest)
                }
            }
        case string(events.ROOMESSAGE):
            if len(stringArray) < 3 {
                send(string(events.INVALID), *guest)
            } else if !guest.isIdentified() {
                send(string(events.IDENTIFY_ERROR), *guest)
            } else {
                chatRoomName := stringArray[1]
                if chatRoom, ok := server.chatRooms[chatRoomName]; ok{
                    if chatRoom.hosts(guest){
                        for _ , chatRoomGuest := range chatRoom.getGuests() {
                            stringToSend :=
                            strings.TrimPrefix(message, "ROOMESSAGE" +
                                                " " + chatRoomName + " ")
                            stringToSend = "..." + chatRoomName + "-" +
                                            guest.getUsername() + ": " +
                                            stringToSend
                            if !chatRoomGuest.equals(guest) {
                                send(stringToSend, chatRoomGuest)
                            } else {
                                send("...MESSAGE SENT", *guest)
                            }
                        }
                    } else {
                        send("...YOU ARE NOT PART OF THE ROOM", *guest)
                    }
                } else {
                    send("...ROOM NOT EXISTS", *guest)
                }
            }
        case string(events.DISCONNECT):
            delete(server.guestsByUsername, guest.getUsername())
            delete(server.guestsById, guest.getSerial())
            guest.getConnection().Close()
            return
        default :
            send(string(events.INVALID), *guest)
        }
    }
}

// Method which sets up a server.
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
