package main

import (
	"github.com/alanarteagav/clientControllers"
	//"github.com/alanarteagav/events"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
)

const (
	UI_GLADE = "./resources/UI.glade"
	MAIN_WINDOW_ID = "mainWindow"
	GOCHAT_ICON = "./resources/goChat-icon.png"
)


func start(){


	builder, err := gtk.BuilderNewFromFile(UI_GLADE)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	mainController := clientControllers.NewMainWindowController(*builder)
	windowObject, err := builder.GetObject(MAIN_WINDOW_ID)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	mainWindow, ok := windowObject.(*gtk.Window)
	if !ok {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	mainWindow.Connect("destroy", mainController.Exit)
	mainWindow.ShowAll()
}

func main(){
	gtk.Init(&os.Args)
	start()
	gtk.Main()
}
