package main

import (
    "bufio"
    "fmt"
    "github.com/alanarteagav/client"
    "log"
    "os"
    "strconv"
)

func main() {

    osArguments := os.Args
    if len(osArguments) != 3 {
        fmt.Println("[ Use : ./goClient address port ]")
        os.Exit(0)
    }

    ipAddress := osArguments[1]
    portString := osArguments[2]

    port, convertionError := strconv.Atoi(portString)
    if convertionError != nil{
        fmt.Println("[ Sorry, invalid port ]")
        os.Exit(0)
    }

    reader := bufio.NewReader(os.Stdin)
    fmt.Println()
    fmt.Print("Enter username: ")
    username, err := reader.ReadString('\n')
    if err != nil {
        log.Fatalln(err)
    }

    mainClient := client.NewClient(username, ipAddress, port)
    mainClient.SendMessage(username + "\n")

    go func() {
        for{
            message := client.Listen(mainClient.GetConnection())
            fmt.Println(":: " +message)
        }
    }()

    for {
        fmt.Println("[ Enter text: ]")
        text, err := reader.ReadString('\n')
        if err != nil {
            log.Fatalln(err)
        }
        mainClient.SendMessage(text)
    }

}
