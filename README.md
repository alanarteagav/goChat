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
http://localhost:6060/pkg/github.com/alanarteagav/
```
