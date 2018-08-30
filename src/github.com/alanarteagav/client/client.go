package client

import "net"

type Client struct {
    username    string
    ipAdress    string
    port        int
    connection  net.Conn
}

func NewClient(username string, ipAdress string, port int) *Client {
    client := new(Client)
    return client
}

func (client Client) GetUsername() string {
    return ""
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
    return ""
}

func (client Client) SendMessage(message string)  {

}
