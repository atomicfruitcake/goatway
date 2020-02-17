# Goatway

## What is this

Goatway is an API Gateway service written in Golang using the Gorilla/MUX framework with a Redis DB. Goatway maintains a register of all jobs received to maintain a log of the states of jobs being processed by downstream services

## How do I use this
You can run the application using

```bash
go run /path/to/goatway/srv/main.go
```
Goatway uses GOLANG version 1.13

Note you must have a Redis server running. Ensure you have redis installed and start by using:

```bash
redis-server
```

If the application start successfully, you will see the message `"Goatway HTTP Webserver is on port 9090"`.

You can conform this by going to running

```bash
curl http://localhost:9090/health
```

And you will receive the response `OK`

## How do I build this
You can compile Goatway down into a binary by building with `go build` using the following commands:

```bash
go build /path/to/goatway//main.go
```

This will compile the program into a binary called `main` in your local directory. You can then run that program with:

```bash
./main
```
