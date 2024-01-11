# syntax=docker/dockerfile:1

## Build
FROM golang:1.21.4-alpine AS build

WORKDIR /opt

COPY go.mod ./
COPY go.sum ./
COPY config/config.yaml /opt/app/

RUN go mod download

COPY . ./

RUN go build -o ./kube-ns-cleaner

## Deploy
FROM alpine:3.17.2

WORKDIR /opt

COPY --from=build /opt/kube-ns-cleaner ./kube-ns-cleaner
COPY --from=build /opt/app/config.yaml /opt/app/

ENTRYPOINT ["/opt/kube-ns-cleaner"]
