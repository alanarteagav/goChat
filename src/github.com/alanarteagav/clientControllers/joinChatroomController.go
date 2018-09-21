package clientControllers

import "github.com/gotk3/gotk3/gtk"
import "log"

const (
    JOIN_ENTRY = "nameEntry"
    JOIN_BUTTON = "joinButton"
)

type joinChatroomController struct{
    builder gtk.Builder
    nameEntry gtk.Entry
    joinButton gtk.Button
}

func NewJoinChatroomController(builder gtk.Builder) *joinChatroomController {
    controller := new(joinChatroomController)
    controller.builder = builder
    buttonObject, err := controller.builder.GetObject(JOIN_BUTTON)
        if err != nil {  log.Fatal(err.Error())  }
    joinButton, ok := buttonObject.(*gtk.Button)
        if !ok { log.Fatal(err.Error()) }
    entryObject, err := controller.builder.GetObject(JOIN_ENTRY)
        if err != nil {  log.Fatal(err.Error())  }
    nameEntry, ok := entryObject.(*gtk.Entry)
        if !ok {  log.Fatal(err.Error())  }

    controller.joinButton = *joinButton
    controller.nameEntry = *nameEntry

    return controller
}

func (controller joinChatroomController) GetButton() gtk.Button {
    return controller.joinButton
}

func (controller joinChatroomController) GetNameEntry() gtk.Entry {
    return controller.nameEntry
}

func (controller joinChatroomController) GetName() (string, bool) {
    username, inputError := controller.nameEntry.GetText()
    if inputError != nil {
        return "", true
    } else if username == ""{
        return "", true
    } else {
        return username, false
    }
}
