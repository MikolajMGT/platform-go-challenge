# syntax = docker/dockerfile:1.2

ARG NAME=assets

FROM golang:1.20.5-alpine3.18 as builder
ARG NAME

WORKDIR /go/src/app

COPY . ${NAME}

ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
RUN cd ${NAME} && go mod download && go mod tidy && /usr/local/go/bin/go build -o /go/src/app/${NAME}-service cmd/main.go

FROM alpine:3.18
ARG NAME

COPY config.yaml /etc/assets/
COPY --from=builder /go/src/app/${NAME}-service /usr/local/bin/${NAME}-service
WORKDIR /usr/local/bin/

ENV app_cmd="./${NAME}-service"
ENTRYPOINT $app_cmd
