# goChat
A simple chat written in go (includes both client and server)

## Downloading GO

First of all, we need to download go in order to be able to
use goChat, here is the link to Golang's download page.

[Golang download page](https://golang.org/dl/)

## Setting GOPATH

In order to use goChat, you should configure your GOPATH.
First, move to the goChat directory, using :

```
cd /(the path to the project directory)/goChat
```

and then run these commands in your terminal.
```
export GOPATH=$(pwd)
export GOBIN=$(pwd)/bin
```

## Getting dependencies

To get goChat's necessary dependencies, you must go to the "goChat" directory:
```
cd /(the path to the project directory)/goChat
```

and then run this command:
```
go get ./...
```
### Warning
Some Gotk3 project files are defective and prevent goChat from complying properly, so it is necessary to remove them, that can be done removing the following files:

 - goChat/src/github.com/gotk3/gotk3/gtk/gtk_since_3_20.go
 - goChat//src/github.com/gotk3/gotk3/gtk/shortcutswindow_since_3_22.go
 - goChat/src/github.com/gotk3/gotk3/gtk/menu_since_3_22.go

Or running the bash script "build.sh" included in the project using the command
```
. build.sh
```


## Building goChat

Once you are in the goChat directory, you must run this command:
```
go install ./...
```

(you should also see the Gotk3 project Readme to see if you have
the necessary dependencies installed on your computer, some
of them are GTK+3, GDK 3, GLib 2 and Cairo).

[Gotk3 Project Page](https://github.com/gotk3/gotk3)

## Running unit tests

Once you are in the goChat directory, simply run this command:
```
go test ./...
```

alternatively, you can test the server and client files with the following commands (once you are in the goChat directory):
```
cd src/github.com/alanarteagav/server/
go test
```

and for the client:
```
cd src/github.com/alanarteagav/client/
go test
```


## Running client and server

If you are in the goChat directory, run:
```
./bin/goServer port
```
to run the server.

And:
```
./bin/goClient
```
to run the client.


## goClient use

Type any line in the text entry to send a public message to the global
chat, and type
```
@user (your message here)
```
to send a private message to a specific user.

And:
```
#chatroom (your message here)
```
to send a message to a specific chatroom.


## Generating Documentation

To generate documentation, run:
```
godoc -http=:8080
```
where '8080' is the port in which the computer will open a server to
show the documentation.

Then, open your web browser and type in your address bar the following
address:
```
http://localhost:8080/pkg/github.com/alanarteagav/
```
(you can also click in the following link )
[goChat documentation (unexported)](http://localhost:8080/pkg/github.com/alanarteagav/?m=all)
[goChat documentation (general)](http://localhost:8080/pkg/github.com/alanarteagav/)
