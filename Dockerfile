FROM golang:1.26.1-alpine3.23 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o main ./cmd/main.go

# Stage 2: Run the Go binary
FROM alpine:3.23.3
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /build/main /app/
COPY config /app/config
EXPOSE 8000
CMD ["./main"]
