# Anime Streamer

A full-stack anime streaming platform that downloads torrents and streams them via HLS (HTTP Live Streaming).

## Features

- 🎬 Download anime via torrent magnet links
- 🎥 Automatic HLS conversion with multiple quality options (720p, 480p, 360p)
- 📺 Adaptive streaming with quality selection
- 🔄 Real-time download and conversion progress tracking
- 💾 Persistent storage for downloads and converted files
- 🎨 Clean, responsive UI

## Tech Stack

### Backend
- **Go** - High-performance backend server
- **Gin** - Web framework
- **anacrolix/torrent** - Torrent client
- **FFmpeg** - Video conversion to HLS

### Frontend
- **React** - UI framework
- **TypeScript** - Type-safe JavaScript
- **HLS.js** - HLS video player

## Project Structure

```
anime-streamer/
├── backend/
│   ├── main.go           # Server setup and routes
│   ├── types.go          # Anime type definitions
│   ├── handlers.go       # HTTP request handlers
│   ├── torrent.go        # Torrent download logic
│   ├── hls.go            # HLS conversion logic
│   ├── scanner.go        # File scanning and restoration
│   ├── utils.go          # Utility functions
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── hooks/        # Custom React hooks
│   │   ├── types/        # TypeScript types
│   │   └── config.ts     # Configuration
│   ├── nginx.conf
│   └── Dockerfile
├── storage/              # Created automatically
│   ├── downloads/        # Torrent downloads
│   └── hls/              # HLS converted files
└── docker-compose.yml
```

## Getting Started

### Prerequisites

- Docker and Docker Compose installed on your system

### Running with Docker

1. Clone the repository:
```bash
git clone <your-repo-url>
cd anime-streamer
```

2. Start the application:
```bash
docker-compose up -d
```

3. Access the application:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

4. Stop the application:
```bash
docker-compose down
```

### Development Setup

#### Backend

1. Navigate to backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run .
```

The backend will start on `http://localhost:8080`

**Note:** FFmpeg must be installed on your system for HLS conversion.

#### Frontend

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm start
```

The frontend will start on `http://localhost:3000`

## Configuration

### Environment Variables

Create a `.env` file based on `.env.example`:

```bash
# Backend Configuration
BACKEND_PORT=8080

# Frontend Configuration
REACT_APP_API_URL=http://localhost:8080
```

### Docker Configuration

To change ports or other settings, edit `docker-compose.yml`:

- Frontend port: Change `3000:80` to `<your-port>:80`
- Backend port: Change `8080:8080` to `<your-port>:8080`

## Usage

1. **Add Anime**: Enter anime name and torrent magnet link, select desired qualities
2. **Download**: The backend automatically downloads the torrent
3. **Conversion**: After download completes, automatic HLS conversion begins
4. **Stream**: Once ready, click play to stream with adaptive quality selection

## Storage

All data is persisted in the `./storage` directory:
- `downloads/` - Original torrent downloads
- `hls/` - Converted HLS files with multiple quality variants

The application automatically restores anime from existing files on startup.

## API Endpoints

- `GET /api/anime` - List all anime
- `POST /api/anime` - Add new anime
- `GET /api/anime/:id` - Get anime details
- `DELETE /api/anime/:id` - Delete anime
- `GET /api/anime/:id/progress` - Get download/conversion progress
- `POST /api/anime/:id/convert` - Manually trigger HLS conversion
- `GET /hls/*` - Serve HLS streaming files

## License

MIT

## Notes

- Make sure you have sufficient disk space for downloads and conversions
- HLS conversion can be CPU-intensive
- Only use torrent links for content you have rights to download
