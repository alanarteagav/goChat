// Main file for goServer program.
package main

import (
    "fmt"
    "github.com/alanarteagav/server"
    "os"
    "os/signal"
    "strconv"
    "syscall"
)

func main()  {
    osArguments := os.Args
    if len(osArguments) != 2 {
        fmt.Println("[ Use : ./goServer port ]")
        os.Exit(0)
    }
    portString := osArguments[1]
    port, convertionError := strconv.Atoi(portString)
    if convertionError != nil{
        fmt.Println("[ Sorry, invalid port ]")
        os.Exit(0)
    }
    server := server.NewServer(port)

    signals := make(chan os.Signal, 1)
    quitSignal := make(chan bool, 1)

    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        signal := <- signals
        fmt.Println()
        fmt.Println(signal)
        quitSignal <- true
    }()

    go server.Serve()
    fmt.Println("[Launching server ...........]")
    fmt.Println("[Listening in the port : " + portString + "]")

    <- quitSignal
    fmt.Println("[Signal received, closing server ... ]")
    os.Exit(0)

}
