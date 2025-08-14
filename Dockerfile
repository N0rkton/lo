FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o main cmd/main.go

FROM alpine

WORKDIR /build

ARG HOST="localhost:8080"
ENV HOST_PORT=${HOST}
EXPOSE 8080

COPY --from=builder /build/main /build/main

CMD ["/build/main"]