# Bash script for building goChat
export CGO_ENABLED="1"
export GOPATH=$(pwd)/
export GOBIN=$(pwd)/bin

# Remove BIN DIR
if [ -e ./bin/ ]
then
    echo "[ REMOVING BIN DIR ]"
    rm -r bin/
fi

# Remove BIN DIR
if [ -e ./pkg/ ]
then
    echo "[ REMOVING PKG DIR ]"
    rm -r pkg/
fi


# Remove BAD FILE 1
if [ -e ./src/github.com/gotk3/gotk3/gtk/gtk_since_3_20.go ]
then
    echo "removing gtk_since_3_20.go (BAD FILE 2)"
    rm ./src/github.com/gotk3/gotk3/gtk/gtk_since_3_20.go
else
    echo "not removing gtk_since_3_20.go (BAD FILE 2)"
fi

# Remove BAD FILE 2
if [ -e ./src/github.com/gotk3/gotk3/gtk/menu_since_3_22.go ]
then
    echo "removing menu_since_3_22.go (BAD FILE 2)"
    rm ./src/github.com/gotk3/gotk3/gtk/menu_since_3_22.go
else
    echo "not removing menu_since_3_22.go (BAD FILE 2)"
fi

# Remove BAD FILE 3
if [ -e ./src/github.com/gotk3/gotk3/gtk/shortcutswindow_since_3_22.go ]
then
    echo "removing shortcutswindow_since_3_22.go (BAD FILE 3)"
    rm ./src/github.com/gotk3/gotk3/gtk/shortcutswindow_since_3_22.go
else
    echo "not removing menu_since_3_22.go (BAD FILE 3)"
fi
#go build -o ./src/github.com/alanarteagav/goClient/client.so -buildmode=c-shared ./src/github.com/alanarteagav/client/client.go
go build ./...
go install ./...
#gcc -o bin/goClientC ./src/github.com/alanarteagav/goClient/gtkClient.c -Wall `pkg-config --cflags --libs gtk+-3.0` -export-dynamic
