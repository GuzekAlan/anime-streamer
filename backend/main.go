package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create storage directories
	os.MkdirAll("../storage/downloads", 0755)
	os.MkdirAll("../storage/hls", 0755)

	// Scan existing files and restore anime list
	scanExistingAnime()

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Routes
	api := r.Group("/api")
	{
		api.GET("/anime", getAnimeList)
		api.POST("/anime", addAnime)
		api.GET("/anime/:id", getAnime)
		api.DELETE("/anime/:id", deleteAnime)
		api.GET("/anime/:id/progress", getDownloadProgress)
		api.POST("/anime/:id/convert", convertToHLS)
	}

	// Serve HLS files
	r.Static("/hls", "../storage/hls")

	log.Println("Server starting on :8080")
	r.Run(":8080")
}

type Anime struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	TorrentURL  string            `json:"torrent_url"`
	Status      string            `json:"status"` // downloading, converting, ready, error
	Progress    int               `json:"progress"`
	HLSPath     string            `json:"hls_path,omitempty"`     // Master playlist path
	HLSUrls     map[string]string `json:"hls_urls,omitempty"`    // Quality -> URL mapping
	VideoPath   string            `json:"video_path,omitempty"`  // Original video file path
	Qualities   []string          `json:"qualities,omitempty"`
	CreatedAt   string            `json:"created_at"`
}

var animeList = make(map[string]*Anime)

func getAnimeList(c *gin.Context) {
	var list []*Anime
	for _, anime := range animeList {
		list = append(list, anime)
	}
	c.JSON(http.StatusOK, gin.H{"anime": list})
}

func getAnime(c *gin.Context) {
	id := c.Param("id")
	anime, exists := animeList[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}
	c.JSON(http.StatusOK, anime)
}

func addAnime(c *gin.Context) {
	var req struct {
		Name       string   `json:"name" binding:"required"`
		TorrentURL string   `json:"torrent_url" binding:"required"`
		Qualities  []string `json:"qualities"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default to all qualities if none specified
	selectedQualities := req.Qualities
	if len(selectedQualities) == 0 {
		selectedQualities = []string{"720p", "480p", "360p"}
	}

	anime := &Anime{
		ID:         generateID(),
		Name:       req.Name,
		TorrentURL: req.TorrentURL,
		Status:     "downloading",
		Progress:   0,
		Qualities:  selectedQualities, // Store selected qualities
		CreatedAt:  getCurrentTime(),
	}

	animeList[anime.ID] = anime

	// Start download in background
	go downloadTorrent(anime)

	c.JSON(http.StatusCreated, anime)
}

func deleteAnime(c *gin.Context) {
	id := c.Param("id")
	anime, exists := animeList[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}
	
	// Clean up active torrent if it exists
	cleanupTorrent(id)
	
	// Remove from anime list
	delete(animeList, id)
	
	log.Printf("Deleted anime: %s (%s)", anime.Name, id)
	c.JSON(http.StatusOK, gin.H{"message": "Anime deleted"})
}

func getDownloadProgress(c *gin.Context) {
	id := c.Param("id")
	anime, exists := animeList[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   anime.Status,
		"progress": anime.Progress,
	})
}

func convertToHLS(c *gin.Context) {
	id := c.Param("id")
	anime, exists := animeList[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}

	if anime.Status != "ready" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anime not ready for conversion"})
		return
	}

	go convertVideoToHLS(anime)
	c.JSON(http.StatusOK, gin.H{"message": "Conversion started"})
}

func scanExistingAnime() {
	log.Println("Scanning existing anime files...")
	
	// Scan HLS directory first (these are ready to stream)
	scanHLSDirectory()
	
	// Scan downloads directory for completed downloads not yet converted
	scanDownloadsDirectory()
	
	log.Printf("Found %d existing anime", len(animeList))
}

func scanHLSDirectory() {
	hlsDir := "../storage/hls"
	entries, err := os.ReadDir(hlsDir)
	if err != nil {
		log.Printf("Could not read HLS directory: %v", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			animeID := entry.Name()
			hlsPath := filepath.Join(hlsDir, animeID)
			
			// Check if master.m3u8 exists
			masterPlaylist := filepath.Join(hlsPath, "master.m3u8")
			if _, err := os.Stat(masterPlaylist); err == nil {
				// Find available qualities
				qualities := scanAvailableQualities(hlsPath)
				
				// Create HLS URLs map
				hlsUrls := make(map[string]string)
				for _, quality := range qualities {
					hlsUrls[quality] = fmt.Sprintf("/hls/%s/%s.m3u8", animeID, quality)
				}
				
				// Try to get anime name from directory structure or use ID
				animeName := getAnimeNameFromFiles(hlsPath, animeID)
				
				anime := &Anime{
					ID:        animeID,
					Name:      animeName,
					Status:    "ready",
					Progress:  100,
					HLSPath:   fmt.Sprintf("/hls/%s/master.m3u8", animeID),
					HLSUrls:   hlsUrls,
					Qualities: qualities,
					CreatedAt: getCurrentTime(),
				}
				
				animeList[animeID] = anime
				log.Printf("Restored HLS anime: %s (%s)", animeName, animeID)
			}
		}
	}
}

func scanDownloadsDirectory() {
	downloadsDir := "../storage/downloads"
	entries, err := os.ReadDir(downloadsDir)
	if err != nil {
		log.Printf("Could not read downloads directory: %v", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Check if this download has a corresponding HLS version
			dirName := entry.Name()
			
			// Skip if we already have this as HLS
			found := false
			for _, anime := range animeList {
				if anime.Status == "ready" && (anime.Name == dirName || anime.VideoPath == filepath.Join(downloadsDir, dirName)) {
					found = true
					break
				}
			}
			
			if !found {
				// Look for video files in this directory
				videoFile := findVideoFileInDir(filepath.Join(downloadsDir, dirName))
				if videoFile != "" {
					animeID := generateID()
					
					anime := &Anime{
						ID:        animeID,
						Name:      dirName,
						Status:    "ready", // Downloaded but not converted
						Progress:  100,
						VideoPath: videoFile,
						Qualities: []string{"720p", "480p", "360p"}, // Default qualities
						CreatedAt: getCurrentTime(),
					}
					
					animeList[animeID] = anime
					log.Printf("Restored downloaded anime: %s (%s)", dirName, animeID)
				}
			}
		}
	}
}

func scanAvailableQualities(hlsPath string) []string {
	var qualities []string
	qualityFiles := []string{"720p.m3u8", "480p.m3u8", "360p.m3u8"}
	
	for _, qualityFile := range qualityFiles {
		if _, err := os.Stat(filepath.Join(hlsPath, qualityFile)); err == nil {
			quality := strings.TrimSuffix(qualityFile, ".m3u8")
			qualities = append(qualities, quality)
		}
	}
	
	return qualities
}

func getAnimeNameFromFiles(hlsPath, fallbackID string) string {
	// Try to read a metadata file if it exists
	metadataFile := filepath.Join(hlsPath, "metadata.txt")
	if data, err := os.ReadFile(metadataFile); err == nil {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		if len(lines) > 0 && lines[0] != "" {
			return lines[0] // First line is the anime name
		}
	}
	
	// Fallback to using the directory name or ID
	return fmt.Sprintf("Anime_%s", fallbackID[:8])
}

func findVideoFileInDir(dirPath string) string {
	videoExts := []string{".mp4", ".mkv", ".avi", ".mov", ".wmv"}
	var videoFile string
	
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		ext := strings.ToLower(filepath.Ext(path))
		for _, validExt := range videoExts {
			if ext == validExt {
				videoFile = path
				return filepath.SkipDir
			}
		}
		return nil
	})
	
	return videoFile
}

// Utility functions
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}