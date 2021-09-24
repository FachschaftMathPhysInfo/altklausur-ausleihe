FROM golang:1.17 as build
MAINTAINER Christian Heusel <christian@heusel.eu>

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download
ENV CGO_ENABLED=0

COPY server/graph ./server/graph
COPY server/lti_utils ./server/lti_utils

COPY exam_marker/exam_marker.go ./exam_marker/exam_marker.go
COPY utils ./utils

RUN go build -a -o exam_marker.runnable ./exam_marker/exam_marker.go

FROM scratch AS exam_marker
COPY --from=build /go/src/exam_marker.runnable ./
ENTRYPOINT ["./exam_marker.runnable"]
