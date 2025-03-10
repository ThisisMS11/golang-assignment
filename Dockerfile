FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN go build -o kiranaclub ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kiranaclub .
COPY store_master.json .

EXPOSE 8080

CMD ["./kiranaclub"]