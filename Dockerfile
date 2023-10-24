# syntax = docker/dockerfile:1.4
# Dockerfile syntax versions: https://hub.docker.com/r/docker/dockerfile
# Dockerfile References: https://docs.docker.com/engine/reference/builder/
FROM golang:1.19-alpine AS builder

LABEL maintainer="ming"

RUN apk update && apk --no-cache add build-base git

# Go debugging tools
RUN go install github.com/google/gops@latest

WORKDIR /root

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build cd /root/src && \
  go build -tags musl -o prom2dd

######## Start a new stage from scratch #######
FROM alpine

# RUN apk update
WORKDIR /root/bin
RUN mkdir /root/config/

# Copy debug tools
COPY --from=builder /go/bin/gops /usr/bin

# Copy the Pre-built binary file and default configurations from the previous stage
COPY --from=builder /root/src/prom2dd /root/bin
RUN ln -s /root/bin/prom2dd /bin/prom2dd

# Command to run the executable
ENTRYPOINT ["./prom2dd"]
