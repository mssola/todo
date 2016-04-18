FROM golang:1.6-alpine
MAINTAINER Miquel Sabaté Solà <mikisabate@gmail.com>

COPY . /go/src/github.com/mssola/todo
WORKDIR /go/src/github.com/mssola/todo

ENV TODO_DEPLOY 1

RUN go build -ldflags="-s -w" \
  && apk --no-cache add --update -t deps ruby && gem install sass --no-ri --no-rdoc \
  && apk --no-cache add --update bash openssl \
  && ./script/sass \
  && rm -r public/stylesheets/*.scss; rm -r public/stylesheets/include \
  && rm -r app lib vendor Godeps script \
  && apk del --purge deps; rm -rf /tmp/* /var/cache/apk/*

ENTRYPOINT ["./todo"]
EXPOSE 3000
