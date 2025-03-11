# Use Golang Alpine as the builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o kiranaclub ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kiranaclub .

COPY storeMaster.csv .

EXPOSE 8080

# Run the application
CMD ["./kiranaclub"]

