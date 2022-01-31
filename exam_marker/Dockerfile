FROM golang:1.17.6-bullseye
# Specify pdfium version
ARG PdfiumVersion=4861
RUN DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt-get -yqq update
RUN DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt-get -y install tzdata
# Install pkg-config, etc.
RUN  apt-get clean && apt-get install -yqq apt-utils pkg-config  wget

# Create .pc file for pkg-config
RUN echo "\n" \
     "prefix=/home\n" \
     "Name: pdfium\n" \
     "Description: pdfium\n" \
     "Version: $PdfiumVersion\n" \
     "Requires:\n" \
     "Libs: -L/home/lib -lpdfium\n" \
     "Cflags: -I/home/include\n" >  /home/pdfium.pc

# Download and extract pdfium binary
RUN cd /home && wget --quiet  https://github.com/bblanchon/pdfium-binaries/releases/download/chromium%2F4861/pdfium-linux-x64.tgz \
&& tar -xf pdfium-linux-x64.tgz && rm pdfium-linux-x64.tgz

# Setting up paths for pkg-config
ENV LD_LIBRARY_PATH=/home/lib
ENV PKG_CONFIG_PATH=/home/

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY server/graph ./server/graph
COPY server/lti_utils ./server/lti_utils

COPY exam_marker/exam_marker.go ./exam_marker/exam_marker.go
COPY utils ./utils
RUN go get github.com/brunsgaard/go-pdfium-render
RUN go build -a -o exam_marker.runnable ./exam_marker/exam_marker.go

ENTRYPOINT ["./exam_marker.runnable"]
