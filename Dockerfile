FROM golang:1.17 as build-base
MAINTAINER Christian Heusel <christian@heusel.eu>

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY utils ./utils
ENV CGO_ENABLED=0

COPY server/graph ./server/graph
COPY server/lti_utils ./server/lti_utils

# SERVER
FROM build-base as server-build

COPY server/server.go ./server/
RUN go build -a -o graphQLServer.runnable server/server.go

FROM scratch AS server

COPY --from=server-build /go/src/graphQLServer.runnable ./
ENTRYPOINT ["./graphQLServer.runnable"]

# EXAM_MARKER
FROM build-base as exam_marker-build

COPY exam_marker/exam_marker.go ./exam_marker/exam_marker.go
RUN go build -a -o exam_marker.runnable ./exam_marker/exam_marker.go

FROM scratch AS exam_marker

COPY --from=exam_marker-build /go/src/exam_marker.runnable ./
ENTRYPOINT ["./exam_marker.runnable"]
