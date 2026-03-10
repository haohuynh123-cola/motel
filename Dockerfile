# Build stage
FROM golang:1.24-alpine AS builder

# Install SSL ca certificates
RUN apk add --no-cache ca-certificates

ENV GO111MODULE=on

WORKDIR /app

# Check out dependencies first to cache them
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/api .
# Copy the SSL certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy thư mục migrations vào image cuối cùng
COPY --from=builder /app/db/migrations ./db/migrations

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./api"]
