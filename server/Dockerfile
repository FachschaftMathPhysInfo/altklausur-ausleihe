FROM golang:1.18-bullseye
MAINTAINER Christian Heusel <christian@heusel.eu>

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY server/graph ./server/graph
COPY server/lti_utils ./server/lti_utils
COPY server/server.go server/dummylogin.html ./server/
COPY utils ./utils

RUN go build -a -o graphQLServer.runnable server/server.go
ENTRYPOINT ["./graphQLServer.runnable"]
