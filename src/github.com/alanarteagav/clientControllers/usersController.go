package clientControllers

import "github.com/gotk3/gotk3/gtk"
import "github.com/gotk3/gotk3/glib"
import "log"
import "strings"

const (
    USERS_VIEW = "usersView"
)

type usersController struct{
    builder gtk.Builder

    usersView gtk.TextView
}

func NewUsersController(builder gtk.Builder) *usersController {
    controller := new(usersController)
    controller.builder = builder

    textViewObject, err := builder.GetObject(USERS_VIEW)
        if err != nil {  log.Fatal(err.Error())  }
    usersView, ok := textViewObject.(*gtk.TextView)
            if !ok { log.Fatal(err.Error()) }
    controller.usersView = *usersView

    return controller
}

func (controller usersController) DisplayUsers(users string) {
    usersList := strings.Split(users, " ")
    glib.IdleAdd(func ()  {
        buffer, _ := controller.usersView.GetBuffer()
        for _, user := range usersList {
            if user != "" {
                iter := buffer.GetEndIter()
                buffer.Insert(iter, "• " + user + "\n")
            }
        }
    })
}
