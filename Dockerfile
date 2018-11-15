FROM golang:1.11.2
LABEL maintainer=dev@codeship.com
WORKDIR /go/src/github.com/codeship/codeship-go
COPY . .
RUN make setup
