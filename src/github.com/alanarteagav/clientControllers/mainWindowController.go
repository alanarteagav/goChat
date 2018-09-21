package clientControllers

import "github.com/alanarteagav/client"
import "github.com/gotk3/gotk3/gtk"
import "github.com/gotk3/gotk3/glib"
import "os"
import "log"
import "fmt"
import "strings"
import "strconv"

var clientExists bool = false
var clientGlobal client.Client


const (
    CONNECT_WINDOW_GLADE = "./resources/loginWindow.glade"
    CREATE_CHATROOM_GLADE = "./resources/createChatroom.glade"
    JOIN_CHATROOM_GLADE = "./resources/joinChatroom.glade"
    INVITE_CHATROOM_GLADE = "./resources/inviteChatroom.glade"
    USERS_GLADE = "./resources/users.glade"
    DIALOGS_GLADE = "./resources/dialogs.glade"

    GLOBAL_TEXT_VIEW = "globalTextView"
    NOTEBOOK = "notebook"

    INFO_DIALOG = "infoDialog"
    ERROR_DIALOG = "errorDialog"

    LOG_IN_WINDOW = "loginWindow"
    CREATE_CHATROOM_WINDOW = "createChatroomWindow"
    JOIN_CHATROOM_WINDOW = "joinChatroomWindow"
    INVITE_CHATROOM_WINDOW = "inviteChatroomWindow"
    USERS_WINDOW = "usersWindow"

    CONNECT_MENU = "connectMenu"
    DISCONNECT_MENU = "disconnectMenu"
    CREATE_CHATROOM_MENU = "createChatroomMenu"
    JOIN_CHATROOM_MENU = "joinChatroomMenu"
    INVITE_CHATROOM_MENU = "inviteChatroomMenu"
    USERS_MENU = "usersMenu"
)

type MainWindowController interface {
    Show()
    Exit()
}

func setHandlers(controller *mainWindowController) map[string]interface{} {
     handlers := make(map[string]interface{})
     handlers["exit"] = controller.Exit
     handlers["connect"] = controller.triggerConnectWindow
     handlers["disconnect"] = controller.disconnect
     handlers["create"] = controller.triggerCreateWindow
     handlers["join"] = controller.triggerJoinWindow
     handlers["invite"] = controller.triggerInviteWindow
     handlers["users"] = controller.triggerUsersWindow
     handlers["sendMessage"] = controller.sendMessage

     return handlers
}

type mainWindowController struct{
    builder gtk.Builder

    listenChannel chan string

    notebook gtk.Notebook

    usersTextViews map[string]gtk.TextView
    chatroomsTextViews map[string]gtk.TextView
    globalTextView gtk.TextView

    connectMenu gtk.MenuItem
    disconnectMenu gtk.MenuItem
    createChatroomMenu gtk.MenuItem
    joinChatroomMenu gtk.MenuItem
    inviteChatroomMenu gtk.MenuItem
    usersMenu gtk.MenuItem

    listenControl listenController
}

func NewMainWindowController(builder gtk.Builder) *mainWindowController {
    controller := new(mainWindowController)
    controller.builder = builder
    controller.builder.ConnectSignals(setHandlers(controller))
    controller.listenChannel = make(chan string)

    //The chat's notebook
    notebookObject, err := controller.builder.GetObject(NOTEBOOK)
        if err != nil {  log.Fatal(err.Error())  }
    controller.notebook = *notebookObject.(*gtk.Notebook)

    // A textView hash table.
    controller.usersTextViews = make(map[string]gtk.TextView)
    controller.chatroomsTextViews = make(map[string]gtk.TextView)

    //The global chat textView
    textViewObject, err := builder.GetObject(GLOBAL_TEXT_VIEW)
        if err != nil {  log.Fatal(err.Error())  }
    globalTextView, ok := textViewObject.(*gtk.TextView)
        if !ok { log.Fatal(err.Error()) }
    controller.globalTextView = *globalTextView

    //The connect menu from the hamburger Button
    connectMenu, err := controller.builder.GetObject(CONNECT_MENU)
        if err != nil {  log.Fatal(err.Error())  }
    controller.connectMenu = *connectMenu.(*gtk.MenuItem)

    //The disconnect menu from the hamburger Button
    disconnectMenu, err := controller.builder.GetObject(DISCONNECT_MENU)
        if err != nil {  log.Fatal(err.Error())  }
    controller.disconnectMenu = *disconnectMenu.(*gtk.MenuItem)

    //The create Chatroom menu from the hamburger Button
    createChatroomMenu, err := controller.builder.GetObject(CREATE_CHATROOM_MENU)
        if err != nil {  log.Fatal(err.Error())  }
    controller.createChatroomMenu = *createChatroomMenu.(*gtk.MenuItem)

    //The join Chatroom menu from the hamburger Button
    joinChatroomMenu, err := controller.builder.GetObject(JOIN_CHATROOM_MENU)
        if err != nil {  log.Fatal(err.Error())  }
    controller.joinChatroomMenu = *joinChatroomMenu.(*gtk.MenuItem)

    //The invite to Chatroom menu from the hamburger Button
    inviteChatroomMenu, err := controller.builder.GetObject(INVITE_CHATROOM_MENU)
        if err != nil {  log.Fatal(err.Error())  }
    controller.inviteChatroomMenu = *inviteChatroomMenu.(*gtk.MenuItem)

    //The users menu from the hamburger Button
    usersMenu, err := controller.builder.GetObject(USERS_MENU)
        if err != nil {  log.Fatal(err.Error())  }
    controller.usersMenu = *usersMenu.(*gtk.MenuItem)

    // Sets the start state for the chat's menus.
    controller.disconnectMenu.SetSensitive(false)
    controller.createChatroomMenu.SetSensitive(false)
    controller.joinChatroomMenu.SetSensitive(false)
    controller.inviteChatroomMenu.SetSensitive(false)
    controller.usersMenu.SetSensitive(false)
    controller.globalTextView.SetSizeRequest(500, 300)

    return controller
}


func (controller *mainWindowController) triggerConnectWindow()  {
    connectBuilder, err := gtk.BuilderNewFromFile(CONNECT_WINDOW_GLADE)
        if err != nil {  log.Fatal(err.Error())  }
    windowObject, err := connectBuilder.GetObject(LOG_IN_WINDOW)
        if err != nil {  log.Fatal(err.Error())  }
    connectWindow, ok := windowObject.(*gtk.Window)
        if !ok { log.Fatal(err.Error()) }
    connectController := NewLoginWindowController(*connectBuilder)

    addressEntry := connectController.GetAddressEntry()
    portEntry := connectController.GetPortEntry()
    usernameEntry := connectController.GetUsernameEntry()
    button := connectController.GetButton()

    addressEntry.Connect("activate",
        func(){ controller.connect(*connectController, *connectWindow) })
    portEntry.Connect("activate",
        func(){ controller.connect(*connectController, *connectWindow) })
    usernameEntry.Connect("activate",
        func(){ controller.connect(*connectController, *connectWindow) })
    button.Connect("clicked",
        func(){ controller.connect(*connectController, *connectWindow) })
    connectWindow.Show()
}

func (controller *mainWindowController) connect(
    loginController loginWindowController, loginWindow gtk.Window) {

    dialogBuilder, err := gtk.BuilderNewFromFile(DIALOGS_GLADE)
        if err != nil {  log.Fatal(err.Error())  }

    address, addressError := loginController.GetAddress()
    port, portError := loginController.GetPort()
    username, _ := loginController.GetUsername()

    if addressError && portError {
        dialogObject, err := dialogBuilder.GetObject(INFO_DIALOG)
            if err != nil {  log.Fatal(err.Error())  }
        infoDialog, ok := dialogObject.(*gtk.MessageDialog)
            if !ok { log.Fatal(err.Error()) }
        infoDialog.SetMarkup("Enter valid server.\nPress ESC to close.")
        infoDialog.Connect("close", infoDialog.Close)
        infoDialog.Connect("response", infoDialog.Close)
        infoDialog.Show()

    } else if addressError {
        dialogObject, err := dialogBuilder.GetObject(INFO_DIALOG)
            if err != nil {  log.Fatal(err.Error())  }
        infoDialog, ok := dialogObject.(*gtk.MessageDialog)
            if !ok { log.Fatal(err.Error()) }
        infoDialog.SetMarkup("Enter valid address.\nPress ESC to close.")
        infoDialog.Connect("close", infoDialog.Close)
        infoDialog.Connect("response", infoDialog.Close)
        infoDialog.Show()

    } else if portError {
        infoDialog :=
            gtk.MessageDialogNew(nil, gtk.DIALOG_MODAL,
                                 gtk.MESSAGE_INFO, gtk.BUTTONS_NONE,
                                "Enter valid port (bigger than 1024)\n" +
                                "Press ESC to close.")
        infoDialog.Show()
    } else {
        newClient := client.NewClient(username)
        connectionError := newClient.Connect(address, port)
        if connectionError != nil {
            dialogObject, err := dialogBuilder.GetObject(ERROR_DIALOG)
                if err != nil {  log.Fatal(err.Error())  }
            errorDialog, ok := dialogObject.(*gtk.MessageDialog)
                if !ok { log.Fatal(err.Error()) }
            errorDialog.SetMarkup("Cannot connect to Server.\n" +
                                  "Press ESC to close." )
            errorDialog.Connect("close", errorDialog.Close)
            errorDialog.Connect("response", errorDialog.Close)
            errorDialog.Show()
        } else {
            clientGlobal = *newClient
            clientGlobal.SendMessage("IDENTIFY " + username)
            response, _ := newClient.Listen()
            if response == "...USERNAME NOT AVAILABLE" {
                dialogObject, err := dialogBuilder.GetObject(ERROR_DIALOG)
                    if err != nil {  log.Fatal(err.Error())  }
                errorDialog, ok := dialogObject.(*gtk.MessageDialog)
                    if !ok { log.Fatal(err.Error()) }
                errorDialog.SetMarkup("Username not available.\n" +
                                      "Press ESC to close.")
                errorDialog.Connect("close", errorDialog.Close)
                errorDialog.Connect("response", errorDialog.Close)
                errorDialog.Show()

                portEntry := loginController.GetPortEntry()
                portEntry.SetText(strconv.Itoa(port))
                portEntry.SetSensitive(false)

                addressEntry := loginController.GetAddressEntry()
                addressEntry.SetText(address)
                addressEntry.SetSensitive(false)
            } else {
                clientExists = true

                controller.listenControl = *NewListenController(
                                                clientGlobal,
                                                controller.usersTextViews,
                                                controller.chatroomsTextViews,
                                                controller.globalTextView)

                go controller.listenControl.listen(controller.listenChannel,
                                                   controller.notebook)

                controller.connectMenu.SetSensitive(false)
                controller.disconnectMenu.SetSensitive(true)
                controller.createChatroomMenu.SetSensitive(true)
                controller.joinChatroomMenu.SetSensitive(true)
                controller.inviteChatroomMenu.SetSensitive(true)
                controller.usersMenu.SetSensitive(true)
                loginWindow.Close()
            }
        }
    }
}

func (controller *mainWindowController) disconnect()  {
    clientGlobal.SendMessage("DISCONNECT")

    controller.connectMenu.SetSensitive(true)
    controller.disconnectMenu.SetSensitive(false)
    controller.createChatroomMenu.SetSensitive(false)
    controller.joinChatroomMenu.SetSensitive(false)
    controller.inviteChatroomMenu.SetSensitive(false)
    controller.usersMenu.SetSensitive(false)

    buffer, _ := controller.globalTextView.GetBuffer()
    buffer.SetText("")
}

func (controller *mainWindowController) triggerCreateWindow()  {
    createBuilder, err := gtk.BuilderNewFromFile(CREATE_CHATROOM_GLADE)
        if err != nil {  log.Fatal(err.Error())  }
    windowObject, err := createBuilder.GetObject(CREATE_CHATROOM_WINDOW)
        if err != nil {  log.Fatal(err.Error())  }
    createWindow, ok := windowObject.(*gtk.Window)
        if !ok { log.Fatal(err.Error()) }
    createController := NewCreateChatroomController(*createBuilder)

    nameEntry := createController.GetNameEntry()
    button := createController.GetButton()

    nameEntry.Connect("activate",
        func(){ controller.create(*createController, *createWindow) })
    button.Connect("clicked",
        func(){ controller.create(*createController, *createWindow) })
    createWindow.Show()
}

func (controller *mainWindowController) create(
    createController createChatroomController,
    createWindow gtk.Window) {

    dialogBuilder, err := gtk.BuilderNewFromFile(DIALOGS_GLADE)
        if err != nil { log.Fatal(err.Error()) }

    name, _ := createController.GetName()

    clientGlobal.SendMessage("CREATEROOM " + name)
    response := <- controller.listenChannel
    if response == "...ROOM NAME ALREADY IN USE" {
        dialogObject, err := dialogBuilder.GetObject(ERROR_DIALOG)
            if err != nil {  log.Fatal(err.Error())  }
        errorDialog, ok := dialogObject.(*gtk.MessageDialog)
            if !ok { log.Fatal(err.Error()) }
        errorDialog.SetMarkup("Room name already in use.\n Press ESC to close.")
        errorDialog.Connect("close", errorDialog.Close)
        errorDialog.Connect("response", errorDialog.Close)
        errorDialog.Show()
    } else {
        createWindow.Close()
    }
}


func (controller *mainWindowController) triggerJoinWindow()  {
    joinBuilder, err := gtk.BuilderNewFromFile(JOIN_CHATROOM_GLADE)
        if err != nil {  log.Fatal(err.Error())  }
    windowObject, err := joinBuilder.GetObject(JOIN_CHATROOM_WINDOW)
        if err != nil {  log.Fatal(err.Error())  }
    joinWindow, ok := windowObject.(*gtk.Window)
        if !ok { log.Fatal(err.Error()) }
    joinController := NewJoinChatroomController(*joinBuilder)

    nameEntry := joinController.GetNameEntry()
    button := joinController.GetButton()

    nameEntry.Connect("activate",
        func(){ controller.join(*joinController, *joinWindow) })
    button.Connect("clicked",
        func(){ controller.join(*joinController, *joinWindow) })
    joinWindow.Show()
}

func (controller *mainWindowController) join(
    joinController joinChatroomController,
    joinWindow gtk.Window) {

    dialogBuilder, err := gtk.BuilderNewFromFile(DIALOGS_GLADE)
        if err != nil {  log.Fatal(err.Error())  }

    name, _ := joinController.GetName()
    clientGlobal.SendMessage("JOINROOM " + name)
    response := <- controller.listenChannel
    if response == "...YOU ARE NOT INVITED TO ROOM " + name {
        dialogObject, err := dialogBuilder.GetObject(INFO_DIALOG)
            if err != nil {  log.Fatal(err.Error())  }
        errorDialog, ok := dialogObject.(*gtk.MessageDialog)
            if !ok { log.Fatal(err.Error()) }
        errorDialog.SetMarkup("You weren't invited to the room.\n" +
                              "Press ESC to close.")
        errorDialog.Connect("close", errorDialog.Close)
        errorDialog.Connect("response", errorDialog.Close)
        errorDialog.Show()
    } else if response == "...ROOM NOT EXISTS" {
        dialogObject, err := dialogBuilder.GetObject(ERROR_DIALOG)
            if err != nil {  log.Fatal(err.Error())  }
        errorDialog, ok := dialogObject.(*gtk.MessageDialog)
            if !ok { log.Fatal(err.Error()) }
        errorDialog.SetMarkup("Room not exists.\n" +
                              "Press ESC to close.")
        errorDialog.Connect("close", errorDialog.Close)
        errorDialog.Connect("response", errorDialog.Close)
        errorDialog.Show()
    } else {
        joinWindow.Close()
    }
}

func (controller *mainWindowController) triggerInviteWindow()  {
    inviteBuilder, err := gtk.BuilderNewFromFile(INVITE_CHATROOM_GLADE)
        if err != nil {  log.Fatal(err.Error())  }
    windowObject, err := inviteBuilder.GetObject(INVITE_CHATROOM_WINDOW)
        if err != nil {  log.Fatal(err.Error())  }
    inviteWindow, ok := windowObject.(*gtk.Window)
        if !ok { log.Fatal(err.Error()) }
    inviteController := NewInviteChatroomController(*inviteBuilder)

    nameEntry := inviteController.GetNameEntry()
    usernameEntry := inviteController.GetUsernameEntry()
    button := inviteController.GetButton()

    nameEntry.Connect("activate",
        func(){ controller.invite(*inviteController, *inviteWindow) })
    usernameEntry.Connect("activate",
        func(){ controller.invite(*inviteController, *inviteWindow) })
    button.Connect("clicked",
        func(){ controller.invite(*inviteController, *inviteWindow) })
    inviteWindow.Show()
}


func (controller *mainWindowController) invite(
    inviteController inviteChatroomController,
    inviteWindow gtk.Window) {

    dialogBuilder, err := gtk.BuilderNewFromFile(DIALOGS_GLADE)
        if err != nil {  log.Fatal(err.Error())  }

    name, _ := inviteController.GetName()
    username, _ := inviteController.GetUsername()
    clientGlobal.SendMessage("INVITE " + name + " " + username)
    fmt.Println("INVITE " + name + " " + username)
    response := <- controller.listenChannel
    fmt.Println(response)
    if response == "...YOU ARE NOT THE OWNER OF THE ROOM"{
        dialogObject, err := dialogBuilder.GetObject(INFO_DIALOG)
            if err != nil {  log.Fatal(err.Error())  }
        errorDialog, ok := dialogObject.(*gtk.MessageDialog)
            if !ok { log.Fatal(err.Error()) }
        errorDialog.SetMarkup("You are not the owner of the room.\n" +
                              "Press ESC to close.")
        errorDialog.Connect("close", errorDialog.Close)
        errorDialog.Connect("response", errorDialog.Close)
        errorDialog.Show()
    } else if response == "...ROOM NOT EXIST" {
        dialogObject, err := dialogBuilder.GetObject(ERROR_DIALOG)
            if err != nil {  log.Fatal(err.Error())  }
        errorDialog, ok := dialogObject.(*gtk.MessageDialog)
            if !ok { log.Fatal(err.Error()) }
        errorDialog.SetMarkup("Room not exists.\n" +
                              "Press ESC to close.")
        errorDialog.Connect("close", errorDialog.Close)
        errorDialog.Connect("response", errorDialog.Close)
        errorDialog.Show()
    } else if response == "...USER " + username + " NOT FOUND" {
        dialogObject, err := dialogBuilder.GetObject(ERROR_DIALOG)
            if err != nil {  log.Fatal(err.Error())  }
        errorDialog, ok := dialogObject.(*gtk.MessageDialog)
            if !ok { log.Fatal(err.Error()) }
        errorDialog.SetMarkup("User not found.\n" +
                              "Press ESC to close.")
        errorDialog.Connect("close", errorDialog.Close)
        errorDialog.Connect("response", errorDialog.Close)
        errorDialog.Show()
    } else {
        inviteWindow.Close()
    }
}

func (controller *mainWindowController) triggerUsersWindow()  {
    usersBuilder, err := gtk.BuilderNewFromFile(USERS_GLADE)
        if err != nil {  log.Fatal(err.Error())  }
    windowObject, err := usersBuilder.GetObject(USERS_WINDOW)
        if err != nil {  log.Fatal(err.Error())  }
    usersWindow, ok := windowObject.(*gtk.Window)
        if !ok { log.Fatal(err.Error()) }
    usersController := NewUsersController(*usersBuilder)

    clientGlobal.SendMessage("USERS")
    users := <- controller.listenChannel
    usersController.DisplayUsers(users)

    usersWindow.Show()
}

func (controller *mainWindowController) sendMessage()  {
    globalEntryObject, err := controller.builder.GetObject("globalTextEntry")
        if err != nil {  log.Fatal(err.Error())  }
    globalEntry, ok := globalEntryObject.(*gtk.Entry)
        if !ok {  log.Fatal(err.Error())  }

    message, _ := globalEntry.GetText()
    messageSplit := strings.SplitAfterN(message, " ", 2)
    if len(messageSplit) == 2 && clientExists {
        id := messageSplit[0]
        if strings.HasPrefix(id, "@") {
            id = strings.TrimLeft(id, "@")
            id = strings.TrimRight(id, " ")

            clientGlobal.SendMessage("MESSAGE " + id + " " + messageSplit[1])
            response := <- controller.listenChannel

            if response == "...USER " + id + " NOT FOUND" {
                glib.IdleAdd(func ()  {
                    buffer, _ := controller.globalTextView.GetBuffer()
                    iter := buffer.GetEndIter()
                    buffer.Insert(iter, "goChat: The requested user (" + id +
                                  ") doesn't exist.\n\n")
                    globalEntry.SetText("")
                })
            } else {
                if view, ok := controller.usersTextViews[id]; ok {
                    glib.IdleAdd(func ()  {
                        buffer, _ := view.GetBuffer()
                        iter := buffer.GetEndIter()
                        buffer.Insert(iter, "[YOU]" + "\n")
                        buffer.Insert(iter, messageSplit[1] + "\n\n")
                        view.ShowAll()

                        globalEntry.SetText("")
                    })
                } else {
                    glib.IdleAdd(func ()  {
                        privateView, _ := gtk.TextViewNew()
                        label, _ := gtk.LabelNew(id)
                        closeButton, _ :=  gtk.ButtonNew()
                        closeButton.SetLabel("×")
                        hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 3)
                        hBox.PackStart(label, true, true, 0)
                        hBox.PackEnd(closeButton, false, false, 0)
                        tabNumber := controller.notebook.AppendPage(privateView, hBox)
                        closeButton.Connect("clicked",
                            func(){
                                page, _ := controller.notebook.GetNthPage(tabNumber)
                                page.Hide()
                            })
                        hBox.ShowAll()
                        privateView.ShowAll()

                        controller.usersTextViews[id] = *privateView

                        privateBuffer, _ := privateView.GetBuffer()
                        iter := privateBuffer.GetEndIter()
                        privateBuffer.Insert(iter, "[YOU]" + "\n")
                        privateBuffer.Insert(iter, messageSplit[1] + "\n\n")

                        globalEntry.SetText("")
                    })
                }
            }

        } else if strings.HasPrefix(id, "#") {
            id = strings.TrimLeft(id, "#")
            id = strings.TrimRight(id, " ")

            clientGlobal.SendMessage("ROOMESSAGE " + id + " " + messageSplit[1])

            response := <- controller.listenChannel

            if response == "...ROOM NOT EXISTS" {
                glib.IdleAdd(func ()  {
                    buffer, _ := controller.globalTextView.GetBuffer()
                    iter := buffer.GetEndIter()
                    buffer.Insert(iter, "goChat: The requested " +
                                  "chatroom doesn't exist.\n\n")
                    globalEntry.SetText("")
                })
            } else if response == "...YOU ARE NOT PART OF THE ROOM" {
                glib.IdleAdd(func ()  {
                    buffer, _ := controller.globalTextView.GetBuffer()
                    iter := buffer.GetEndIter()
                    buffer.Insert(iter, "goChat: You are not part of the room.\n\n")
                    globalEntry.SetText("")
                })
            } else {
                if view, ok := controller.chatroomsTextViews[id]; ok {
                    glib.IdleAdd(func ()  {
                        buffer, _ := view.GetBuffer()
                        iter := buffer.GetEndIter()
                        buffer.Insert(iter, "[YOU]" + "\n")
                        buffer.Insert(iter, messageSplit[1] + "\n\n")
                        view.ShowAll()
                    })
                } else {
                    glib.IdleAdd(func ()  {
                        roomView, _ := gtk.TextViewNew()
                        label, _ := gtk.LabelNew(id)
                        closeButton, _ :=  gtk.ButtonNew()
                        closeButton.SetLabel("×")
                        hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 3)
                        hBox.PackStart(label, true, true, 0)
                        hBox.PackEnd(closeButton, false, false, 0)
                        tabNumber := controller.notebook.AppendPage(roomView, hBox)


                        closeButton.Connect("clicked",
                            func(){
                                page, _ := controller.notebook.GetNthPage(tabNumber)
                                page.Hide()
                            })
                        hBox.ShowAll()
                        roomView.ShowAll()

                        controller.chatroomsTextViews[id] = *roomView

                        buffer, _ := roomView.GetBuffer()
                        iter := buffer.GetEndIter()
                        buffer.Insert(iter, "[YOU]" + "\n")
                        buffer.Insert(iter, messageSplit[1] + "\n\n")
                    })
                }
            }
        } else {
            glib.IdleAdd(func ()  {
                clientGlobal.SendMessage("PUBLICMESSAGE " + message)
                buffer, _ := controller.globalTextView.GetBuffer()
                iter := buffer.GetEndIter()
                buffer.Insert(iter, "[YOU]" + "\n")
                buffer.Insert(iter, message + "\n\n")
                globalEntry.SetText("")
            })
        }
    } else if clientExists {
        glib.IdleAdd(func ()  {
            clientGlobal.SendMessage("PUBLICMESSAGE " + message)
            buffer, _ := controller.globalTextView.GetBuffer()
            iter := buffer.GetEndIter()
            buffer.Insert(iter, "[YOU]" + "\n")
            buffer.Insert(iter, message + "\n\n")
            globalEntry.SetText("")
        })
    } else {
        glib.IdleAdd(func ()  {
            globalEntry.SetText("")
        })
    }
}

func (controller mainWindowController) Exit()  {
    os.Exit(0)
}
