FROM golang:1.21-alpine AS builder

# Install FFmpeg and other dependencies
RUN apk add --no-cache ffmpeg git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

FROM alpine:latest

# Install FFmpeg in the final image
RUN apk add --no-cache ffmpeg

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Create storage directories
RUN mkdir -p /app/storage/downloads /app/storage/hls

EXPOSE 8080

CMD ["./main"]