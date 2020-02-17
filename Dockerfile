FROM golang:latest

LABEL maintainer="atomicfruitcake"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build ./srv/main.go   

EXPOSE 9090

CMD ["./main"]