package clientControllers

import "github.com/gotk3/gotk3/gtk"
import "strconv"
import "log"

const (
    PORT_ENTRY = "portEntry"
    ADDRESS_ENTRY = "addressEntry"
    USERNAME_ENTRY = "usernameEntry"
    CONNECT_BUTTON = "connectButton"
)

type LoginWindowController interface {
    GetAddress()
    GetPort()
}

type loginWindowController struct{
    builder gtk.Builder
    connectButton gtk.Button
    addressEntry gtk.Entry
    portEntry gtk.Entry
    usernameEntry gtk.Entry
}

func NewLoginWindowController(builder gtk.Builder) *loginWindowController {
    controller := new(loginWindowController)
    controller.builder = builder
    buttonObject, err := controller.builder.GetObject(CONNECT_BUTTON)
        if err != nil {  log.Fatal(err.Error())  }
    connectButton, ok := buttonObject.(*gtk.Button)
        if !ok { log.Fatal(err.Error()) }
    addressEntryObject, err := controller.builder.GetObject(ADDRESS_ENTRY)
        if err != nil {  log.Fatal(err.Error())  }
    addressEntry, ok := addressEntryObject.(*gtk.Entry)
        if !ok {  log.Fatal(err.Error())  }
    portEntryObject, err := controller.builder.GetObject(PORT_ENTRY)
        if err != nil {  log.Fatal(err.Error())  }
    portEntry, ok := portEntryObject.(*gtk.Entry)
        if !ok {  log.Fatal(err.Error())  }
    usernameEntryObject, err := controller.builder.GetObject(USERNAME_ENTRY)
        if err != nil {  log.Fatal(err.Error())  }
    usernameEntry, ok := usernameEntryObject.(*gtk.Entry)
        if !ok {  log.Fatal(err.Error())  }
    controller.connectButton = *connectButton
    controller.addressEntry = *addressEntry
    controller.portEntry = *portEntry
    controller.usernameEntry = *usernameEntry
    return controller
}

func (controller loginWindowController) GetButton() gtk.Button {
    return controller.connectButton
}

func (controller loginWindowController) GetAddressEntry() gtk.Entry {
    return controller.addressEntry
}

func (controller loginWindowController) GetPortEntry() gtk.Entry {
    return controller.portEntry
}

func (controller loginWindowController) GetUsernameEntry() gtk.Entry {
    return controller.usernameEntry
}

func (controller loginWindowController) GetAddress() (string, bool) {
    address, inputError := controller.addressEntry.GetText()
    if inputError != nil {
        return "", true
    } else if address == ""{
        return "", true
    } else {
        return address, false
    }
}

func (controller loginWindowController) GetUsername() (string, bool) {
    username, inputError := controller.usernameEntry.GetText()
    if inputError != nil {
        return "", true
    } else if username == ""{
        return "", true
    } else {
        return username, false
    }
}

func (controller loginWindowController) GetPort() (int, bool) {
    portString, inputError := controller.portEntry.GetText()
    if inputError != nil {
        return 0, true
    } else {
        port, numberError := strconv.Atoi(portString)
        if numberError != nil {
            return port, true
        } else if port < 1025 {
            return port, true
        } else {
            return port, false
        }
    }
}
