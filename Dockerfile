FROM golang:1.9 AS build
LABEL maintainer=dev@codeship.com
RUN go get -v github.com/alecthomas/gometalinter && \
	go get -v golang.org/x/tools/cmd/cover
RUN gometalinter --install

FROM build
LABEL maintainer=dev@codeship.com
RUN mkdir -p /go/src/github.com/codeship/codeship-go
ADD . /go/src/github.com/codeship/codeship-go/
WORKDIR /go/src/github.com/codeship/codeship-go
