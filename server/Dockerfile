FROM golang:1.20.3-bullseye AS build
MAINTAINER Christian Heusel <christian@heusel.eu>

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY server/graph ./server/graph
COPY server/lti_utils ./server/lti_utils
COPY server/prometheus ./server/prometheus
COPY server/server.go server/dummylogin.html ./server/
COPY utils ./utils

RUN CGO_ENABLED=0 go build -a -o graphql-server server/server.go

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/src/graphql-server /
ENTRYPOINT ["/graphql-server"]
