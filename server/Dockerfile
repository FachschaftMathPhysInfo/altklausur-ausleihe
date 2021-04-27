FROM golang:1.16
MAINTAINER Christian Heusel <christian@heusel.eu>

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY graph ./graph
COPY server.go ./

RUN go build -a -o graphQLServer server.go
ENTRYPOINT ["./graphQLServer"]
