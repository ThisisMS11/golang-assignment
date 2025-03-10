FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
# If you have a go.sum file
# COPY go.sum ./
# RUN go mod download

COPY . .

RUN go build -o image-processor ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/image-processor .
COPY store_master.json .

EXPOSE 8080

CMD ["./image-processor"]