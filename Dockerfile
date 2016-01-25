FROM golang:1.5
MAINTAINER Miquel Sabaté Solà <mikisabate@gmail.com>

COPY . /go/src/github.com/mssola/todo
WORKDIR /go/src/github.com/mssola/todo

RUN go get github.com/tools/godep && godep go build

RUN apt-get update && apt-get install -y ruby && gem install sass --no-ri --no-rdoc
ENV TODO_DEPLOY 1
RUN ./script/sass

RUN gem uninstall -x sass && apt-get remove -y ruby && apt-get clean && rm -rf /var/cache/apt

ENTRYPOINT ./todo
EXPOSE 3000
