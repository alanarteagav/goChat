package client

import (
    "net"
    "errors"
    "bufio"
    "strconv"
    "strings"
    "time"
)

// Client struct.
// Defines an username, the server's ip address and port, and the client
// connection (golang's equivalent for sockets).
type Client struct {
    username    string
    connection  net.Conn
}

// Client constructor.
// Receives an username and the server's ip address and port.
// It automatically dials to the server using the exported function Dial
// from the net package.
func NewClient(username string) *Client {
    client := new(Client)
    client.username = username
    client.connection = nil

    return client
}


func (client *Client) Connect(ipAddress string, port int) error {

    success := make(chan net.Conn)
    fail := make(chan error)

    go func () {
        connection, err := net.Dial("tcp", ipAddress +  ":" + strconv.Itoa(port))
        if err != nil {
            fail <- errors.New("Cannot connect to server")
        }
        success <- connection
    }()

    select {
    case <-time.After(5 * time.Second):
        return errors.New("Cannot connect to server | Timed Out")
    case err := <- fail :
        return err
    case connection := <- success :
        client.connection = connection
        return nil
    }
}


// Returns client's username.
func (client Client) GetUsername() string {
    return client.username
}

// Sets a new username for the client.
func (client *Client) SetUsername(username string) {
    client.username = username
}

// Returns client's connection (socket).
func (client Client) GetConnection() net.Conn {
    return client.connection
}

// Sets a new connection (socket) for the client.
func (client *Client) SetConnection(connection net.Conn) {
    client.connection = connection
}

// Listen function.
// Receives a connection (socket) to listen from, and returns a string
// (if a string can be retrieved from the connection).
func (client Client) Listen() (string, error) {
        message, err := bufio.NewReader(client.connection).ReadString('\n')
        message = strings.Trim(message, "\n")
        if err != nil{
            return "", errors.New("Cannot hear through connection")
        }
        return message, nil
}

// SendMessage method.
// Receives a string and sends it through the client's connection.
func (client Client) SendMessage(message string)  {
    client.connection.Write([]byte(message + "\n"))
}
