ARG GO_VERSION=1.16.6

FROM golang:${GO_VERSION}-alpine AS builder

RUN go env -w GO111MODULE=direct
RUN apk add --no-cache git
RUN apk add --no-cache add ca-certificates && update-ca-certificates

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download