FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o /dashboard-service main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /dashboard-service /dashboard-service

EXPOSE 8084

CMD ["/dashboard-service"]
