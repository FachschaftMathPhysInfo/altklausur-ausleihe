FROM golang:1.20.4-bullseye
MAINTAINER Christian Heusel <christian@heusel.eu>

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

# Install pkg-config, etc.
RUN DEBIAN_FRONTEND=noninteractive apt-get -y update && \
    apt-get clean && \
    apt-get install -y apt-utils pkg-config wget

# Specify pdfium version
ARG PdfiumVersion=5619

# Create .pc file for pkg-config
RUN mkdir /go/src/extern && echo "\n" \
    "prefix=/go/src/extern\n" \
    "Name: pdfium\n" \
    "Description: pdfium\n" \
    "Version: $PdfiumVersion\n" \
    "Requires:\n" \
    "Libs: -L/go/src/extern/lib -lpdfium\n" \
    "Cflags: -I/go/src/extern/include\n" > /go/src/extern/pdfium.pc

# Download and extract pdfium binary
RUN cd /go/src/extern && \
    wget --quiet  https://github.com/bblanchon/pdfium-binaries/releases/download/chromium%2F$PdfiumVersion/pdfium-linux-x64.tgz && \
    tar -xf pdfium-linux-x64.tgz && \
    rm pdfium-linux-x64.tgz

# Setting up paths for pkg-config
ENV LD_LIBRARY_PATH=/go/src/extern/lib
ENV PKG_CONFIG_PATH=/go/src/extern/


COPY server/graph ./server/graph
COPY server/lti_utils ./server/lti_utils

COPY exam_marker/exam_marker.go ./exam_marker/exam_marker.go
COPY exam_marker/prometheus ./exam_marker/prometheus
COPY utils ./utils
RUN go build -a -o exam_marker.runnable ./exam_marker/exam_marker.go

ENTRYPOINT ["./exam_marker.runnable"]
