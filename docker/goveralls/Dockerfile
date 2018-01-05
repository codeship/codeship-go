FROM golang:1.9
LABEL maintainer=dev@codeship.com

RUN mkdir -p /go/src/github.com/codeship/codeship-go

RUN go get golang.org/x/tools/cmd/cover && \
    go get github.com/mattn/goveralls

ADD . /go/src/github.com/codeship/codeship-go/
WORKDIR /go/src/github.com/codeship/codeship-go
