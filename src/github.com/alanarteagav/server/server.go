package server

type Server struct {
    port int
}

func NewServer(port int) *Server {
    server := new(Server)
    server.port = port
    return server
}

func (server Server) GetPort() int {
    return server.port
}

func (server Server) Serve()  {

}
