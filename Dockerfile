# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install swag and generate swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN swag init -g cmd/http/main.go -o api/swagger
RUN sed -i '/LeftDelim:/d' api/swagger/docs.go && sed -i '/RightDelim:/d' api/swagger/docs.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/goclean-http cmd/http/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/goclean-grpc cmd/grpc/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh goclean

# Set working directory
WORKDIR /app

# Copy built binaries from builder stage
COPY --from=builder /app/bin/goclean-http /app/bin/goclean-grpc ./

# Copy any additional files if needed
COPY --from=builder /app/api ./api

# Change ownership to non-root user
RUN chown -R goclean:goclean /app

# Switch to non-root user
USER goclean

# Expose ports
EXPOSE 8080 9090

# Default command (can be overridden in docker-compose)
CMD ["./goclean-http"]