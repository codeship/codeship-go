FROM golang:1.8
LABEL maintainer=dev@codeship.com

RUN mkdir -p /go/src/github.com/codeship/codeship-go

ADD . /go/src/github.com/codeship/codeship-go/
WORKDIR /go/src/github.com/codeship/codeship-go
RUN make setup
