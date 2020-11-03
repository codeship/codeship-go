ARG GO_VERSION=1.15

FROM golang:${GO_VERSION}
LABEL maintainer=dev@codeship.com

# go 1.13.15
RUN go get golang.org/dl/go1.13.15 && \
    go1.13.15 download

# go 1.14.10
RUN go get golang.org/dl/go1.14.10 && \
    go1.14.10 download

WORKDIR /go/src/github.com/codeship/codeship-go
COPY . .

RUN make setup
