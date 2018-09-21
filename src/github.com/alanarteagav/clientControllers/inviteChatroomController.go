package clientControllers

import "github.com/gotk3/gotk3/gtk"
import "log"

const (
    INVITE_CHAT_ENTRY = "nameEntry"
    INVITE_USER_ENTRY = "usernameEntry"
    INVITE_BUTTON = "inviteButton"
)

type inviteChatroomController struct{
    builder gtk.Builder
    nameEntry gtk.Entry
    usernameEntry gtk.Entry
    inviteButton gtk.Button
}

func NewInviteChatroomController(builder gtk.Builder) *inviteChatroomController {
    controller := new(inviteChatroomController)
    controller.builder = builder
    buttonObject, err := controller.builder.GetObject(INVITE_BUTTON)
        if err != nil {  log.Fatal(err.Error())  }
    inviteButton, ok := buttonObject.(*gtk.Button)
        if !ok { log.Fatal(err.Error()) }
    entryNameObject, err := controller.builder.GetObject(INVITE_CHAT_ENTRY)
        if err != nil {  log.Fatal(err.Error())  }
    nameEntry, ok := entryNameObject.(*gtk.Entry)
        if !ok {  log.Fatal(err.Error())  }
    entryUsernameObject, err := controller.builder.GetObject(INVITE_USER_ENTRY)
        if err != nil {  log.Fatal(err.Error())  }
    usernameEntry, ok := entryUsernameObject.(*gtk.Entry)
        if !ok {  log.Fatal(err.Error())  }

    controller.inviteButton = *inviteButton
    controller.nameEntry = *nameEntry
    controller.usernameEntry = *usernameEntry

    return controller
}

func (controller inviteChatroomController) GetButton() gtk.Button {
    return controller.inviteButton
}

func (controller inviteChatroomController) GetNameEntry() gtk.Entry {
    return controller.nameEntry
}
func (controller inviteChatroomController) GetUsernameEntry() gtk.Entry {
    return controller.usernameEntry
}

func (controller inviteChatroomController) GetName() (string, bool) {
    username, inputError := controller.nameEntry.GetText()
    if inputError != nil {
        return "", true
    } else if username == ""{
        return "", true
    } else {
        return username, false
    }
}

func (controller inviteChatroomController) GetUsername() (string, bool) {
    username, inputError := controller.usernameEntry.GetText()
    if inputError != nil {
        return "", true
    } else if username == ""{
        return "", true
    } else {
        return username, false
    }
}
