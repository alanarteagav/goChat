package client

import (
    "net"
    "fmt"
    "bufio"
    "strconv"
    "strings"
    "os"
    "log"
)

// Client struct.
// Defines an username, the server's ip address and port, and the client
// connection (golang's equivalent for sockets).
type Client struct {
    username    string
    ipAddress    string
    port        int
    connection  net.Conn
}

// Client constructor.
// Receives an username and the server's ip address and port.
// It automatically dials to the server using the exported function Dial
// from the net package.
func NewClient(username string, ipAddress string, port int) *Client {
    client := new(Client)
    client.username = username
    connection, err := net.Dial("tcp", ipAddress +  ":" + strconv.Itoa(port))
    if err != nil {
        log.Fatalln(err)
        fmt.Println("Unable to connect to server")
    }
    client.ipAddress = ipAddress
    client.port = port
    client.connection = connection
    return client
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
func Listen(connection net.Conn) string {
    for {
        message, err := bufio.NewReader(connection).ReadString('\n')
        message = strings.Trim(message, "\n")
        if err != nil{
            fmt.Println("The server is off")
            os.Exit(1)
        }
        return message
    }
}

// SendMessage method.
// Receives a string and sends it through the client's connection.
func (client Client) SendMessage(message string)  {
    client.connection.Write([]byte(message))
}
