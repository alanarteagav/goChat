package clientControllers

import "github.com/alanarteagav/client"
import "github.com/gotk3/gotk3/gtk"
import "github.com/gotk3/gotk3/glib"
import "os"
import "log"
import "fmt"
import "strconv"

var clientExists bool = false
var clientGlobal client.Client


const (
    LOGIN_WINDOW_GLADE = "./resources/loginWindow.glade"
    CREATE_CHATROOM_GLADE = "./resources/createChatroom.glade"
    JOIN_CHATROOM_GLADE = "./resources/joinChatroom.glade"
    DIALOGS_GLADE = "./resources/dialogs.glade"

    GLOBAL_TEXT_VIEW = "globalTextView"

    INFO_DIALOG = "infoDialog"
    ERROR_DIALOG = "errorDialog"

    LOG_IN_WINDOW = "loginWindow"
    CREATE_CHATROOM_WINDOW = "createChatroomWindow"
    JOIN_CHATROOM_WINDOW = "joinChatroomWindow"

    CONNECT_MENU = "connectMenu"
    DISCONNECT_MENU = "disconnectMenu"
    CREATE_CHATROOM_MENU = "createChatroomMenu"
    JOIN_CHATROOM_MENU = "joinChatroomMenu"
    INVITE_CHATROOM_MENU = "inviteChatroomMenu"
)

type MainWindowController interface {
    Show()
    Exit()
}

func setHandlers(controller *mainWindowController) map[string]interface{} {
     handlers := make(map[string]interface{})
     handlers["exit"] = controller.Exit
     handlers["connect"] = controller.connect
     handlers["disconnect"] = controller.disconnect
     handlers["newChatroom"] = controller.createChatroom
     handlers["joinChatroom"] = controller.joinChatroom
     handlers["inviteChatroom"] = controller.createChatroom
     handlers["sendMessage"] = controller.sendMessage

     return handlers
}

type mainWindowController struct{
    builder gtk.Builder

    globalTextView gtk.TextView

    connectMenu gtk.MenuItem
    disconnectMenu gtk.MenuItem
    createChatroomMenu gtk.MenuItem
    joinChatroomMenu gtk.MenuItem
    inviteChatroomMenu gtk.MenuItem

    lController listenController

    listenChannel chan string
}

func NewMainWindowController(builder gtk.Builder) *mainWindowController {
    controller := new(mainWindowController)
    controller.builder = builder
    controller.builder.ConnectSignals(setHandlers(controller))
    controller.listenChannel = make(chan string)

    textViewObject, err := builder.GetObject(GLOBAL_TEXT_VIEW)
        if err != nil {  log.Fatal(err.Error())  }
    globalTextView, ok := textViewObject.(*gtk.TextView)
        if !ok { log.Fatal(err.Error()) }
    controller.globalTextView = *globalTextView

    connectMenu, err := controller.builder.GetObject(CONNECT_MENU)
        if err != nil {  log.Fatal(err.Error())  }
    controller.connectMenu = *connectMenu.(*gtk.MenuItem)

    disconnectMenu, err := controller.builder.GetObject(DISCONNECT_MENU)
        if err != nil {  log.Fatal(err.Error())  }
    controller.disconnectMenu = *disconnectMenu.(*gtk.MenuItem)

    createChatroomMenu, err := controller.builder.GetObject(CREATE_CHATROOM_MENU)
        if err != nil {  log.Fatal(err.Error())  }
    controller.createChatroomMenu = *createChatroomMenu.(*gtk.MenuItem)

    joinChatroomMenu, err := controller.builder.GetObject(JOIN_CHATROOM_MENU)
        if err != nil {  log.Fatal(err.Error())  }
    controller.joinChatroomMenu = *joinChatroomMenu.(*gtk.MenuItem)

    inviteChatroomMenu, err := controller.builder.GetObject(INVITE_CHATROOM_MENU)
        if err != nil {  log.Fatal(err.Error())  }
    controller.inviteChatroomMenu = *inviteChatroomMenu.(*gtk.MenuItem)

    controller.disconnectMenu.SetSensitive(false)
    controller.createChatroomMenu.SetSensitive(false)
    controller.joinChatroomMenu.SetSensitive(false)
    controller.inviteChatroomMenu.SetSensitive(false)
    controller.globalTextView.SetSizeRequest(500, 300)
    return controller
}


func (controller *mainWindowController) connect()  {
    loginBuilder, err := gtk.BuilderNewFromFile(LOGIN_WINDOW_GLADE)
        if err != nil {  log.Fatal(err.Error())  }
    windowObject, err := loginBuilder.GetObject(LOG_IN_WINDOW)
        if err != nil {  log.Fatal(err.Error())  }
    loginWindow, ok := windowObject.(*gtk.Window)
        if !ok { log.Fatal(err.Error()) }
    loginController := NewLoginWindowController(*loginBuilder)

    connectButton := loginController.GetButton()
    addressEntry := loginController.GetAddressEntry()
    portEntry := loginController.GetPortEntry()
    usernameEntry := loginController.GetUsernameEntry()

    connectButton.Connect("clicked",
        func(){
            controller.connectServer(*loginController, *loginWindow)
        })
    addressEntry.Connect("activate",
        func(){
            controller.connectServer(*loginController, *loginWindow)
        })
    portEntry.Connect("activate",
        func(){
            controller.connectServer(*loginController, *loginWindow)
        })
    usernameEntry.Connect("activate",
        func(){
            controller.connectServer(*loginController, *loginWindow)
        })

    loginWindow.Show()
}


func (controller *mainWindowController) disconnect()  {
    clientGlobal.SendMessage("DISCONNECT")
    controller.connectMenu.SetSensitive(true)
    controller.disconnectMenu.SetSensitive(false)
    controller.createChatroomMenu.SetSensitive(false)
    controller.joinChatroomMenu.SetSensitive(false)
    buffer, _ := controller.globalTextView.GetBuffer()
    buffer.SetText("")
}

func (controller *mainWindowController) createChatroom()  {
    createChatroomBuilder, err := gtk.BuilderNewFromFile(CREATE_CHATROOM_GLADE)
        if err != nil {  log.Fatal(err.Error())  }
    windowObject, err := createChatroomBuilder.GetObject(CREATE_CHATROOM_WINDOW)
        if err != nil {  log.Fatal(err.Error())  }
    createChatroomWindow, ok := windowObject.(*gtk.Window)
        if !ok { log.Fatal(err.Error()) }
    createChatroomController := NewCreateChatroomController(*createChatroomBuilder)

    createButton := createChatroomController.GetButton()
    nameEntry := createChatroomController.GetNameEntry()

    createButton.Connect("clicked",
        func(){
            controller.newChatroom(*createChatroomController, *createChatroomWindow)
        })
    nameEntry.Connect("activate",
        func(){
            controller.newChatroom(*createChatroomController, *createChatroomWindow)
        })
    createChatroomWindow.Show()
}


func (controller *mainWindowController) joinChatroom()  {
    joinChatroomBuilder, err := gtk.BuilderNewFromFile(JOIN_CHATROOM_GLADE)
        if err != nil {  log.Fatal(err.Error())  }
    windowObject, err := joinChatroomBuilder.GetObject(JOIN_CHATROOM_WINDOW)
        if err != nil {  log.Fatal(err.Error())  }
    joinChatroomWindow, ok := windowObject.(*gtk.Window)
        if !ok { log.Fatal(err.Error()) }
    joinChatroomController := NewJoinChatroomController(*joinChatroomBuilder)

    joinButton := joinChatroomController.GetButton()
    nameEntry := joinChatroomController.GetNameEntry()

    joinButton.Connect("clicked",
        func(){
            controller.join(*joinChatroomController, *joinChatroomWindow)
        })
    nameEntry.Connect("activate",
        func(){
            controller.join(*joinChatroomController, *joinChatroomWindow)
        })
    joinChatroomWindow.Show()
}


func (controller *mainWindowController) newChatroom(
    createChatroomController createChatroomController,
    createChatroomWindow gtk.Window) {

    dialogBuilder, err := gtk.BuilderNewFromFile(DIALOGS_GLADE)
        if err != nil {  log.Fatal(err.Error())  }

    name, _ := createChatroomController.GetName()

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
        createChatroomWindow.Close()
    }
}

func (controller *mainWindowController) join(
    joinChatroomController joinChatroomController,
    joinChatroomWindow gtk.Window) {

    dialogBuilder, err := gtk.BuilderNewFromFile(DIALOGS_GLADE)
        if err != nil {  log.Fatal(err.Error())  }

    name, _ := joinChatroomController.GetName()
    fmt.Println(name)

    clientGlobal.SendMessage("JOINROOM " + name)
    response := <- controller.listenChannel
    fmt.Println(response)
    if response == "...YOU ARE NOT INVITED TO ROOM " + name {
        dialogObject, err := dialogBuilder.GetObject(INFO_DIALOG)
            if err != nil {  log.Fatal(err.Error())  }
        errorDialog, ok := dialogObject.(*gtk.MessageDialog)
            if !ok { log.Fatal(err.Error()) }
        errorDialog.SetMarkup("You weren't invited to the room.\n Press ESC to close.")
        errorDialog.Connect("close", errorDialog.Close)
        errorDialog.Connect("response", errorDialog.Close)
        errorDialog.Show()
    } else if response == "...ROOM NOT EXISTS" {
        dialogObject, err := dialogBuilder.GetObject(ERROR_DIALOG)
            if err != nil {  log.Fatal(err.Error())  }
        errorDialog, ok := dialogObject.(*gtk.MessageDialog)
            if !ok { log.Fatal(err.Error()) }
        errorDialog.SetMarkup("Room not exists.\n Press ESC to close.")
        errorDialog.Connect("close", errorDialog.Close)
        errorDialog.Connect("response", errorDialog.Close)
        errorDialog.Show()
    } else {
        joinChatroomWindow.Close()
    }
}

func (controller *mainWindowController) connectServer(
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
            infoDialog.SetMarkup("Enter valid server.\n Press ESC to close." )
            infoDialog.Connect("close", infoDialog.Close)
            infoDialog.Connect("response", infoDialog.Close)

            infoDialog.Show()

        } else if addressError {
            dialogObject, err := dialogBuilder.GetObject(INFO_DIALOG)
                if err != nil {  log.Fatal(err.Error())  }
            infoDialog, ok := dialogObject.(*gtk.MessageDialog)
                if !ok { log.Fatal(err.Error()) }
            infoDialog.SetMarkup("Enter valid address.\n Press ESC to close." )
            infoDialog.Connect("close", infoDialog.Close)
            infoDialog.Connect("response", infoDialog.Close)
            infoDialog.Show()


        } else if portError {
            infoDialog :=
                gtk.MessageDialogNew(nil, gtk.DIALOG_MODAL,
                                     gtk.MESSAGE_INFO, gtk.BUTTONS_NONE,
                                     "Entre una puerto válido (mayor a 1024)")
            infoDialog.Show()

        } else {
            newClient := client.NewClient(username)
            connectionError := newClient.Connect(address, port)
            if connectionError != nil {
                dialogObject, err := dialogBuilder.GetObject(ERROR_DIALOG)
                    if err != nil {  log.Fatal(err.Error())  }
                errorDialog, ok := dialogObject.(*gtk.MessageDialog)
                    if !ok { log.Fatal(err.Error()) }
                errorDialog.SetMarkup("Cannot connect to Server.\n Press ESC to close." )
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
                    errorDialog.SetMarkup("Username not available.\n Press ESC to close.")
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

                    controller.lController = *NewListenController(clientGlobal, controller.globalTextView)
                    go controller.lController.listen(controller.listenChannel)

                    controller.connectMenu.SetSensitive(false)
                    controller.disconnectMenu.SetSensitive(true)
                    controller.createChatroomMenu.SetSensitive(true)
                    controller.joinChatroomMenu.SetSensitive(true)
                    loginWindow.Close()
                }
            }
        }
}

func (controller mainWindowController) sendMessage()  {

    globalEntryObject, err := controller.builder.GetObject("globalTextEntry")
        if err != nil {  log.Fatal(err.Error())  }
    globalEntry, ok := globalEntryObject.(*gtk.Entry)
        if !ok {  log.Fatal(err.Error())  }

    if clientExists {
        message, _ := globalEntry.GetText()
        glib.IdleAdd(func ()  {
            clientGlobal.SendMessage("PUBLICMESSAGE " + message)
            globalEntry.SetText("")
        })
    }
}

func (controller mainWindowController) Exit()  {
    os.Exit(0)
}
