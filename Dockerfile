FROM opensuse/amd64:latest
MAINTAINER Miquel Sabaté Solà <mikisabate@gmail.com>

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

COPY . /go/src/github.com/mssola/todo
WORKDIR /go/src/github.com/mssola/todo

ENV TODO_DEPLOY 1

RUN zypper ref && zypper in -y go openssl ruby2.1-rubygem-sass \
  && go build -ldflags="-s -w" && ./script/sass \
  && rm -r public/stylesheets/*.scss && rm -r public/stylesheets/*.map \
  && rm -r public/stylesheets/include && rm public/images/snapshot.png \
  && rm -r app lib vendor Godeps script *.go *.yml LICENSE Dockerfile Makefile README.md\
  && zypper clean -a && zypper rm -u -y go ruby2.1-rubygem-sass \
  && rm -rf /tmp/*

ENTRYPOINT ["./todo"]
EXPOSE 3000
