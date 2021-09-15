FROM golang:1.17
MAINTAINER Christian Heusel <christian@heusel.eu>

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY server/graph ./server/graph
COPY server/lti_utils ./server/lti_utils

COPY exam_marker/exam_marker.go ./exam_marker/exam_marker.go
COPY utils ./utils

RUN go build -a -o exam_marker ./exam_marker/exam_marker.go
ENTRYPOINT ["./exam_marker"]
