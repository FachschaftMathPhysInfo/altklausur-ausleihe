FROM golang:1.16
MAINTAINER Christian Heusel <christian@heusel.eu>

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY utils ./utils
COPY exam_marker.go ./

RUN go build -a -o exam_marker exam_marker.go
ENTRYPOINT ["./exam_marker"]
