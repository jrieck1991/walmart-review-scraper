FROM golang:latest

# tools
RUN apt update
RUN apt install -y tree lsof

# deps
RUN apt install -y xvfb firefox-esr default-jre

# go
RUN mkdir /go/src/web_driver
WORKDIR /go/src/web_driver
COPY web_driver .
RUN go get ./...
RUN go build -o driver .