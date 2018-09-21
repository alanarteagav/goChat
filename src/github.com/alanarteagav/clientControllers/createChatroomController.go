package clientControllers

import "github.com/gotk3/gotk3/gtk"
import "log"

const (
    NAME_ENTRY_JOIN = "nameEntry"
    CREATE_BUTTON = "createButton"
)

type createChatroomController struct{
    builder gtk.Builder
    nameEntry gtk.Entry
    createButton gtk.Button
}

func NewCreateChatroomController(builder gtk.Builder) *createChatroomController {
    controller := new(createChatroomController)
    controller.builder = builder
    buttonObject, err := controller.builder.GetObject(CREATE_BUTTON)
        if err != nil {  log.Fatal(err.Error())  }
    createButton, ok := buttonObject.(*gtk.Button)
        if !ok { log.Fatal(err.Error()) }
    entryObject, err := controller.builder.GetObject(NAME_ENTRY_JOIN)
        if err != nil {  log.Fatal(err.Error())  }
    nameEntry, ok := entryObject.(*gtk.Entry)
        if !ok {  log.Fatal(err.Error())  }

    controller.createButton = *createButton
    controller.nameEntry = *nameEntry

    return controller
}

func (controller createChatroomController) GetButton() gtk.Button {
    return controller.createButton
}

func (controller createChatroomController) GetNameEntry() gtk.Entry {
    return controller.nameEntry
}

func (controller createChatroomController) GetName() (string, bool) {
    username, inputError := controller.nameEntry.GetText()
    if inputError != nil {
        return "", true
    } else if username == ""{
        return "", true
    } else {
        return username, false
    }
}
