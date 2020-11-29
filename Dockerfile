FROM golang:latest as go

RUN mkdir /martian
COPY go.mod /martian
COPY go.sum /martian
WORKDIR /martian
RUN go mod download
COPY . /martian
RUN go build .

FROM ubuntu:latest
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y git gettext zip
COPY --from=go /martian/martian /usr/bin/martian
ENTRYPOINT ["/usr/bin/martian"]
