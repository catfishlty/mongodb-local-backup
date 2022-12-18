# build
FROM golang:alpine as builder
WORKDIR /app
ADD . ./
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go build -v -o mlb cmd/mlb/main.go

FROM alpine
MAINTAINER "catfishlty"
WORKDIR /app
RUN apk add --no-cache mongodb-tools
COPY --from=builder /app/mlb /app/mlb
ENV MLB_MONGOEXPORT = "/usr/bin/mongoexport"
ENV MLB_HOST = ""
ENV MLB_PORT = ""
ENV MLB_USERNAME = ""
ENV MLB_PASSWORD = ""
ENV MLB_CRON = ""
ENV MLB_TARGET = "json"
ENV MLB_OUTPUT = "/data"
ENV MLB_LOG = "info"
ENTRYPOINT ["/app/mlb", "start", "-d", "--log", "${MLB_LOG}"]
