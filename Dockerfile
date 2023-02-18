# build
FROM golang:alpine as builder
WORKDIR /app
ADD . ./
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go build -v -o mlb cmd/mlb/main.go
RUN apk add --no-cache mongodb-tools

FROM alpine
MAINTAINER "catfishlty"
WORKDIR /app
COPY --from=builder /app/mlb /app/mlb
COPY --from=builder /usr/bin/mongoexport /usr/bin/mongoexport
ENV MLB_MONGOEXPORT "/usr/bin/mongoexport"
ENV MLB_TYPE json
ENV MLB_OUTPUT /mlb/backup
ENV MLB_LOG_LEVEL info
RUN mkdir -p $MLB_OUTPUT
ENTRYPOINT /app/mlb start -d --log $MLB_LOG_LEVEL
