FROM golang:1.9 AS build
LABEL maintainer=dev@codeship.com
RUN mkdir -p /go/src/github.com/codeship/gov
ADD ./docker/gov /go/src/github.com/codeship/gov
WORKDIR /go/src/github.com/codeship/gov
RUN go get golang.org/x/build/version
RUN go install

# go 1.8.6
RUN go get golang.org/x/build/version/go1.8.6 && \
    go1.8.6 download

# go 1.10.rc1
RUN go get golang.org/x/build/version/go1.10rc1 && \
    go1.10rc1 download

FROM build
LABEL maintainer=dev@codeship.com
RUN mkdir -p /go/src/github.com/codeship/codeship-go
ADD . /go/src/github.com/codeship/codeship-go/
WORKDIR /go/src/github.com/codeship/codeship-go
