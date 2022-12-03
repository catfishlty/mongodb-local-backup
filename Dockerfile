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
ENTRYPOINT ["/app/mlb", "start", "-d"]
