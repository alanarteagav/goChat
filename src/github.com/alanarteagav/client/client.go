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

type Client struct {
    username    string
    ipAdress    string
    port        int
    connection  net.Conn
}

func NewClient(username string, ipAdress string, port int) *Client {
    client := new(Client)
    client.username = username
    connection, err := net.Dial("tcp", ipAdress +  ":" + strconv.Itoa(port))
    if err != nil {
        log.Fatalln(err)
        fmt.Println("Unable to connect to server")
    }
    client.ipAdress = ipAdress
    client.port = port
    client.connection = connection
    return client
}

func (client Client) GetUsername() string {
    return client.username
}

func (client *Client) SetUsername(username string) {
    client.username = username
}

func (client Client) GetConnection() net.Conn {
    return client.connection
}

func (client *Client) SetConnection(connection net.Conn) {
    client.connection = connection
}

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

func (client Client) SendMessage(message string)  {
    client.connection.Write([]byte(message + "\n"))
}
