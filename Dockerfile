ARG GO_VERSION=1.20

FROM golang:${GO_VERSION}-alpine AS builder

RUN go env -w GO111MODULE=on
RUN apk add --no-cache git
RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY database database
COPY events events
COPY feed-service feed-service
COPY models models
COPY pusher-service pusher-service
COPY query-service query-service
COPY repository repository
COPY search search

RUN go install ./...

FROM alpine:latest
WORKDIR /usr/bin

COPY --from=builder /go/bin .

