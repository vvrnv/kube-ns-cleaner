FROM golang:1.19.3-alpine3.16 as build
WORKDIR /app
COPY go.* ./
RUN go mod download
RUN export CGO_ENABLED=0
COPY . ./
RUN go build -v -o kube-ns-cleaner

FROM alpine:3.16.3
RUN apk add --no-cache --update bash curl ca-certificates
COPY --from=build /app/kube-ns-cleaner /app/kube-ns-cleaner
ENTRYPOINT  [ "/app/kube-ns-cleaner" ]
