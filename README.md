# 🎌 Anime Torrent Streaming Platform

A full-stack application for downloading torrents and streaming anime with HLS (HTTP Live Streaming) support.

## ✨ Features

- **Torrent Download**: Add anime via magnet links or torrent URLs
- **HLS Conversion**: Automatic conversion to multiple quality streams (720p, 480p, 360p)
- **Real-time Progress**: Live download and conversion progress tracking
- **Video Player**: Custom HLS video player with quality selection
- **Responsive UI**: Modern React interface with beautiful animations
- **Docker Support**: Easy deployment with Docker Compose

## 🛠 Tech Stack

- **Backend**: Go with Gin framework
- **Frontend**: React with TypeScript
- **Video Processing**: FFmpeg for HLS conversion
- **Torrent Client**: Anacrolix torrent library
- **Streaming**: HLS (HTTP Live Streaming)
- **Containerization**: Docker & Docker Compose

## 📁 Project Structure

```
anime-streaming/
├── backend/                 # Go server
│   ├── main.go             # Main server file
│   ├── torrent.go          # Torrent handling logic
│   ├── Dockerfile          # Backend container
│   ├── go.mod              # Go dependencies
│   └── go.sum              # Go dependency checksums
├── frontend/               # React application
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── hooks/          # Custom React hooks
│   │   ├── types/          # TypeScript types
│   │   ├── App.tsx         # Main App component
│   │   └── index.tsx       # React entry point
│   ├── public/             # Static files
│   ├── package.json        # Node dependencies
│   └── Dockerfile          # Frontend container
├── storage/                # Created at runtime
│   ├── downloads/          # Downloaded torrent files
│   └── hls/               # HLS stream segments
├── docker-compose.yml      # Multi-container setup
├── start.sh               # Quick start script
└── README.md              # This file
```

## 🚀 Quick Start

### Option 1: Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd anime-streaming
   ```

2. **Start the application**
   ```bash
   ./start.sh
   ```
   Or manually:
   ```bash
   docker-compose up --build
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

### Option 2: Manual Setup

#### Prerequisites
- Go 1.21+
- Node.js 18+
- FFmpeg
- Git

#### Backend Setup
```bash
cd backend
go mod tidy
go run main.go torrent.go
```

#### Frontend Setup
```bash
cd frontend
npm install
npm start
```

## 📖 Usage Guide

### Adding Anime

1. **Get a magnet link** or torrent URL for your anime
2. **Open the application** at http://localhost:3000
3. **Fill out the form**:
   - Enter the anime name
   - Paste the magnet link or torrent URL
4. **Click "Add Anime"** to start the download

### Streaming Process

The application follows this workflow:

1. **Download**: Torrent is downloaded to `storage/downloads/`
2. **Convert**: Video files are converted to HLS format with multiple qualities
3. **Stream**: HLS segments are served and can be played in the browser

### Video Player Controls

- **Play/Pause**: Click the video or use the play button
- **Seek**: Use the progress bar to jump to different parts
- **Quality**: Select from available qualities (720p, 480p, 360p)
- **Fullscreen**: Double-click the video (browser dependent)

## 🔧 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/anime` | Get all anime |
| POST | `/api/anime` | Add new anime |
| GET | `/api/anime/:id` | Get specific anime |
| DELETE | `/api/anime/:id` | Delete anime |
| GET | `/api/anime/:id/progress` | Get download progress |
| POST | `/api/anime/:id/convert` | Start HLS conversion |
| GET | `/hls/:id/*` | Serve HLS files |

## 🐳 Docker Configuration

The application uses multi-container Docker setup:

- **Backend**: Go server with torrent and FFmpeg capabilities
- **Frontend**: React development server
- **FFmpeg**: Dedicated container for video processing

### Environment Variables

- `GIN_MODE`: Gin framework mode (debug/release)
- `REACT_APP_API_URL`: Backend API URL for frontend

## 📝 Development

### Adding New Features

1. **Backend**: Add routes in `main.go`, implement logic in separate files
2. **Frontend**: Create components in `src/components/`, add hooks in `src/hooks/`
3. **Types**: Update TypeScript interfaces in `src/types/`

### File Structure Guidelines

- Keep components small and focused
- Use custom hooks for API calls
- Maintain type safety with TypeScript
- Follow Go best practices for backend

## ⚠️ Important Notes

### Legal Considerations
- Only use this application with content you have the legal right to download
- Respect copyright laws in your jurisdiction
- This tool is for educational and personal use only

### Performance Tips
- Ensure sufficient disk space for downloads and HLS segments
- FFmpeg conversion is CPU-intensive
- Consider using SSD storage for better performance

### Security
- The application runs on localhost by default
- Do not expose to public networks without proper security measures
- Consider implementing authentication for production use

## 🛠 Troubleshooting

### Common Issues

**Docker containers won't start**
- Ensure Docker and Docker Compose are installed
- Check if ports 3000 and 8080 are available
- Run `docker-compose logs` to see error messages

**Torrent downloads fail**
- Verify the magnet link is valid
- Check firewall settings
- Ensure sufficient disk space

**Video conversion fails**
- Verify FFmpeg is installed and accessible
- Check video file format compatibility
- Monitor system resources during conversion

**HLS playback issues**
- Ensure browser supports HLS (most modern browsers do)
- Check network connectivity
- Verify HLS files were generated correctly

### Logs and Debugging

- Backend logs: `docker-compose logs backend`
- Frontend logs: `docker-compose logs frontend`
- All logs: `docker-compose logs`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is for educational purposes. Please ensure you comply with all applicable laws and regulations when using this software.

---

**Happy streaming! 🍿**