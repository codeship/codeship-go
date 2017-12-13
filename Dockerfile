FROM golang:1.9-alpine
LABEL maintainer=dev@codeship.com

RUN mkdir -p /go/src/github.com/codeship/codeship-go
RUN apk update && apk add make && apk add git

ADD . /go/src/github.com/codeship/codeship-go/
WORKDIR /go/src/github.com/codeship/codeship-go
RUN make setup
