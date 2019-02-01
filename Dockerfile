ARG GO_VERSION=1.11

FROM golang:${GO_VERSION}
LABEL maintainer=dev@codeship.com

WORKDIR /codeship
COPY . .

ENV GOBIN /codeship/bin
ENV PATH $GOBIN:$PATH

RUN make setup tools
