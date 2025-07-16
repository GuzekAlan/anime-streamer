#!/bin/bash

echo "🎌 Starting Anime Streaming Platform..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Create storage directories
mkdir -p storage/downloads storage/hls

echo "📦 Building and starting containers..."
docker-compose up --build

echo "✅ Application started!"
echo "🌐 Frontend: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"