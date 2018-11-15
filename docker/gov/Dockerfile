FROM golang:1.11.2
LABEL maintainer=dev@codeship.com

# go 1.9.7
RUN go get golang.org/dl/go1.9.7 && \
    go1.9.7 download

# go 1.10.5
RUN go get golang.org/dl/go1.10.5 && \
    go1.10.5 download

WORKDIR /go/src/github.com/codeship/codeship-go
COPY . .
