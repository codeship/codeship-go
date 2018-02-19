FROM golang:1.10
LABEL maintainer=dev@codeship.com
RUN go get -v github.com/alecthomas/gometalinter && \
    go get -v golang.org/x/tools/cmd/cover && \
    gometalinter --install
RUN mkdir -p /go/src/github.com/codeship/codeship-go
ADD . /go/src/github.com/codeship/codeship-go/
WORKDIR /go/src/github.com/codeship/codeship-go
