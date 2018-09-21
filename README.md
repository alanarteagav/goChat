# goChat
A simple chat written in go (includes both client and server)

## Downloading GO

First of all, we need to download go in order to be able to
use goChat, here is the link to Golang's download page.

[Golang download page](https://golang.org/dl/)

## Setting GOPATH

In order to use goChat, you should configure your GOPATH, by running
these commands in your terminal. (suppose you've downloaded the goChat
project in your directory /Users/yourname/files/goChat).

```
export GOPATH=$(pwd)
export GOBIN=$(pwd)/bin
```

## Getting dependencies

To get goChat's necessary dependencies, you must go to the "goChat" directory:
```
cd /Users/yourname/files/goChat
```

and then run this command:
```
go get ./...
```

## Building goChat

Once you are in the goChat directory, you must run this command:
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
http://localhost:8080/pkg/github.com/alanarteagav/
```
(you can also click in the following link )
[goChat documentation](http://localhost:8080/pkg/github.com/alanarteagav/)
