FROM golang:1.17 as build
MAINTAINER Christian Heusel <christian@heusel.eu>

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download
ENV CGO_ENABLED=0

COPY server/graph ./server/graph
COPY server/lti_utils ./server/lti_utils
COPY server/server.go ./server/
COPY utils ./utils

RUN go build -a -o graphQLServer.runnable server/server.go

FROM scratch AS server
COPY --from=build /go/src/graphQLServer.runnable  ./
ENTRYPOINT ["./graphQLServer.runnable"]
