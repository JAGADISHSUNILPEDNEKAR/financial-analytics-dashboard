FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /notification-service main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /notification-service /notification-service

EXPOSE 8085

CMD ["/notification-service"]
