FROM golang:latest

RUN apt-get update

WORKDIR /go/src/server/mafia-rpc

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

EXPOSE 9090

CMD ["/bin/bash",  "-c", "go run ./server/main.go"]
