package clientControllers

import "github.com/gotk3/gotk3/gtk"
import "github.com/gotk3/gotk3/glib"
import "github.com/alanarteagav/client"
import "log"
import "strings"

type ListenController interface {
    listen()
}

type listenController struct{
    chatClient client.Client
    usersTextViews map[string]gtk.TextView
    chatroomsTextViews map[string]gtk.TextView
    globalTextView gtk.TextView
}

func NewListenController(chatClient client.Client,
                         usersTextViews map[string]gtk.TextView,
                         chatroomsTextViews map[string]gtk.TextView,
                         globalTextView gtk.TextView) *listenController {
    controller := new(listenController)
    controller.chatClient = chatClient
    controller.usersTextViews = usersTextViews
    controller.chatroomsTextViews = chatroomsTextViews
    controller.globalTextView = globalTextView

    return controller
}


func containsState(state string) bool {
    if strings.Contains(state, "ACTIVE"){
        return true
    } else if strings.Contains(state, "BUSY"){
        return true
    } else if strings.Contains(state, "AWAY"){
        return true
    } else {
        return false
    }
}


func (controller listenController) listen(listenChannel chan string,
                                          notebook gtk.Notebook){
    globalBuffer, _ := controller.globalTextView.GetBuffer()
    for {
        message, err := controller.chatClient.Listen()
        if err != nil {
            return
        } else if strings.HasPrefix(message, "...INVITATION TO JOIN ") {
            message = strings.TrimPrefix(message, "...")

            dialogBuilder, err := gtk.BuilderNewFromFile(DIALOGS_GLADE)
                if err != nil {  log.Fatal(err.Error())  }

            glib.IdleAdd(func ()  {
                dialogObject, err := dialogBuilder.GetObject(INFO_DIALOG)
                    if err != nil {  log.Fatal(err.Error())  }
                errorDialog, ok := dialogObject.(*gtk.MessageDialog)
                    if !ok { log.Fatal(err.Error()) }
                errorDialog.SetMarkup(message + "\n" +
                                      "Press ESC to close.")
                errorDialog.Connect("close", errorDialog.Close)
                errorDialog.Connect("response", errorDialog.Close)
                errorDialog.Show()
            })
            continue
        } else if strings.HasPrefix(message, "...PUBLIC") {
            message = strings.TrimPrefix(message, "...PUBLIC")
            publicMessage := strings.SplitAfterN(message, ":", 2)
            username := strings.TrimLeft(publicMessage[0], "-")
            if len(publicMessage) == 2 {
                glib.IdleAdd(func ()  {
                    iter := globalBuffer.GetEndIter()
                    globalBuffer.Insert(iter, username + "\n")
                    globalBuffer.Insert(iter, publicMessage[1] + "\n\n")
                    controller.globalTextView.ScrollToIter(iter, 0.2, true,
                                                           0.0, 1.0)
                })
            } else {
                continue
            }
            continue
        } else if strings.HasPrefix(message, "...") &&
                  strings.Contains(message, ":") {
            message = strings.TrimPrefix(message, "...")


            chatRoomAndMessage := strings.SplitAfterN(message, "-", 2)
            chatroomName := "[C] " + chatRoomAndMessage[0]
            chatroomName = strings.TrimRight(chatroomName, "-")
            userAndMessage := strings.SplitAfterN(chatRoomAndMessage[1], ":", 2)
            username := strings.TrimLeft(userAndMessage[0], "-")
            username = strings.TrimRight(username, ":")

            if len(chatRoomAndMessage) + len(userAndMessage) == 4 {
                if view, ok := controller.chatroomsTextViews[chatroomName]; ok {
                    glib.IdleAdd(func ()  {
                        buffer, _ := view.GetBuffer()
                        iter := buffer.GetEndIter()
                        buffer.Insert(iter, userAndMessage[1] + "\n")
                        view.ScrollToIter(iter, 0.0, false, 0.0, 1.0)
                        view.ShowAll()
                    })
                } else {
                    glib.IdleAdd(func ()  {
                        roomView, _ := gtk.TextViewNew()
                        label, _ := gtk.LabelNew(chatroomName)
                        closeButton, _ :=  gtk.ButtonNew()
                        closeButton.SetLabel("×")
                        hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 3)
                        hBox.PackStart(label, true, true, 0)
                        hBox.PackEnd(closeButton, false, false, 0)
                        tabNumber := notebook.AppendPage(roomView, hBox)


                        closeButton.Connect("clicked",
                            func(){
                                page, _ := notebook.GetNthPage(tabNumber)
                                page.Hide()
                            })
                        hBox.ShowAll()
                        roomView.ShowAll()

                        controller.chatroomsTextViews[chatroomName] = *roomView

                        privateBuffer, _ := roomView.GetBuffer()
                        iter := privateBuffer.GetEndIter()
                        privateBuffer.Insert(iter, username + "\n")
                        privateBuffer.Insert(iter, userAndMessage[1] + "\n\n")
                        roomView.ScrollToIter(iter, 0.0, false, 0.0, 1.0)
                    })
                }
            } else {
                continue
            }
            continue
        } else if len(strings.Split(message, " ")) == 2 &&
                  containsState(message) {

            dialogBuilder, err := gtk.BuilderNewFromFile(DIALOGS_GLADE)
                if err != nil {  log.Fatal(err.Error())  }

            glib.IdleAdd(func ()  {
                dialogObject, err := dialogBuilder.GetObject(INFO_DIALOG)
                    if err != nil {  log.Fatal(err.Error())  }
                errorDialog, ok := dialogObject.(*gtk.MessageDialog)
                    if !ok { log.Fatal(err.Error()) }
                errorDialog.SetMarkup(message + "\n" +
                                      "Press ESC to close.")
                errorDialog.Connect("close", errorDialog.Close)
                errorDialog.Connect("response", errorDialog.Close)
                errorDialog.Show()
            })
            continue
        } else {
            go func() {
                listenChannel <- message
            }()
            if strings.Contains(message, ":"){
                privateMessage := strings.SplitAfterN(message, ":", 2)
                if len(privateMessage) == 2 {
                    username := privateMessage[0]
                    if view, ok := controller.usersTextViews[username]; ok {
                        glib.IdleAdd(func ()  {
                            buffer, _ := view.GetBuffer()
                            iter := buffer.GetEndIter()
                            buffer.Insert(iter, username + "\n")
                            buffer.Insert(iter, privateMessage[1] + "\n\n")
                            view.ScrollToIter(iter, 0.0, false, 0.0, 1.0)
                            view.ShowAll()
                        })
                    } else {
                        glib.IdleAdd(func ()  {
                            privateView, _ := gtk.TextViewNew()
                            label, _ := gtk.LabelNew(username)
                            closeButton, _ :=  gtk.ButtonNew()
                            closeButton.SetLabel("×")
                            hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 3)
                            hBox.PackStart(label, true, true, 0)
                            hBox.PackEnd(closeButton, false, false, 0)
                            tabNumber := notebook.AppendPage(privateView, hBox)
                            closeButton.Connect("clicked",
                                func(){
                                    page, _ := notebook.GetNthPage(tabNumber)
                                    page.Hide()
                                })
                            hBox.ShowAll()
                            privateView.ShowAll()

                            controller.usersTextViews[username] = *privateView

                            privateBuffer, _ := privateView.GetBuffer()
                            iter := privateBuffer.GetEndIter()
                            privateBuffer.Insert(iter, username + "\n")
                            privateBuffer.Insert(iter, privateMessage[1] + "\n\n")
                            privateView.ScrollToIter(iter, 0.0, false, 0.0, 1.0)
                        })
                    }
                } else {
                    continue
                }
            }
        }
    }
}
