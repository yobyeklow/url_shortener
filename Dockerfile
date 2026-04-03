FROM golang:1.26.1-alpine3.23 AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
CMD ["./main"]
