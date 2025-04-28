# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go-api

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/go-api .

EXPOSE 8080

CMD ["./go-api"]