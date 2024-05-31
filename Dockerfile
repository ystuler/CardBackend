FROM golang:1.22.1-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go

# Stage 2: Run the Go binary
FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/main /app/
COPY config /app/config
EXPOSE 8000
CMD ["./main"]
