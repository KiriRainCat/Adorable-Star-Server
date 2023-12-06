FROM golang:alpine AS builder

# Building the application
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build ./main.go

# Running the application
FROM alpine

ENV TZ=Asia/Shanghai
ENV GIN_MODE=release

RUN apk add --no-cache chromium

WORKDIR /app

COPY --from=builder /build/main /app/main

COPY ./config/config.yaml /app/config.yaml

EXPOSE 8005

CMD ["/app/main"]