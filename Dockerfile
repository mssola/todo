FROM golang:1.5
MAINTAINER Miquel Sabaté Solà <mikisabate@gmail.com>

ADD . /go/src/github.com/mssola/todo
WORKDIR /go/src/github.com/mssola/todo

RUN go get github.com/tools/godep
RUN godep restore
RUN godep go build

ENTRYPOINT ./todo
EXPOSE 3000
