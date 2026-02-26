FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /user-service main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /user-service /user-service

EXPOSE 8083

CMD ["/user-service"]
