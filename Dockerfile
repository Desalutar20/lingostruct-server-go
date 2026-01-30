FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./api ./cmd/api

CMD ["./api"]
