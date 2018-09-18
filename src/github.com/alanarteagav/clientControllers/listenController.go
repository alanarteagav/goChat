package clientControllers

import "github.com/gotk3/gotk3/gtk"
import "github.com/gotk3/gotk3/glib"
import "github.com/alanarteagav/client"
import "fmt"
import "strings"

type ListenController interface {
    listen()
}

type listenController struct{
    chatClient client.Client
    textView gtk.TextView
}

func NewListenController(chatClient client.Client,
                         textView gtk.TextView) *listenController {
    controller := new(listenController)
    controller.chatClient = chatClient
    controller.textView = textView
    return controller
}

func (controller listenController) listen(listenChannel chan string){
    textBuffer, _ := controller.textView.GetBuffer()
    for {
        message, err := controller.chatClient.Listen()
        if err != nil {
            return
        } else if strings.HasPrefix(message, "...INVITATION TO JOIN ROOM") {
            fmt.Println("INVITED")
            continue
        } else {
            go func() {
                listenChannel <- message
            }()
            glib.IdleAdd(func ()  {
                iter := textBuffer.GetEndIter()
                textBuffer.Insert(iter, message + "\n")
                fmt.Println(message)
            })
        }
    }
}
