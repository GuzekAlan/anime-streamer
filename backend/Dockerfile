# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o anime-streamer .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies (FFmpeg for video conversion)
RUN apk add --no-cache \
    ffmpeg \
    ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/anime-streamer .

# Create storage directories
RUN mkdir -p /storage/downloads /storage/hls

# Expose port
EXPOSE 8080

# Run the application
CMD ["./anime-streamer"]
