# goChat
A simple chat written in go (includes both client and server)

## Setting GOPATH

In order to use goChat, you should configure your GOPATH, by running
these commands in your terminal. (suppose you've downloaded the goChat
project in your directory /Users/yourname/files/goChat).

```
export GOPATH=/Users/yourname/files/goChat
export GOBIN=/Users/yourname/files/goChat/bin
```

## Building goChat

To build goChat, you must go to the "goChat" directory:
```
cd /Users/yourname/files/goChat
```

and then run this command:
```
go install ./...
```

## Running unit tests

Once you are in the goChat directory, simply run this command:
```
go test ./...
```

## Running client and server

If you are in the goChat directory, run:
```
./bin/goServer
```
to run the server.

And:
```
./bin/goClient
```
to run the client.
