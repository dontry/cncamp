# httpserver - A simple HTTP server

## Description
This is a simple HTTP server demo. It is used to demonstrate how to use Docker to build and run a Go application.  

Ping the server with `curl http://localhost:8080/healthz` to get a response of `OK`.

## Local usage 

```bash
go run main.go
```

## Docker usage

```bash
docker build -t httpserver .

docker run -p 8080:8080 httpserver
```

or you can pull the image from Docker Hub

```bash
docker pull dontry/httpserver
```