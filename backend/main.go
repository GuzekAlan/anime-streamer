package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create storage directories
	os.MkdirAll("../storage/downloads", 0755)
	os.MkdirAll("../storage/hls", 0755)

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
		Name       string `json:"name" binding:"required"`
		TorrentURL string `json:"torrent_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	anime := &Anime{
		ID:         generateID(),
		Name:       req.Name,
		TorrentURL: req.TorrentURL,
		Status:     "downloading",
		Progress:   0,
		CreatedAt:  getCurrentTime(),
	}

	animeList[anime.ID] = anime

	// Start download in background
	go downloadTorrent(anime)

	c.JSON(http.StatusCreated, anime)
}

func deleteAnime(c *gin.Context) {
	id := c.Param("id")
	if _, exists := animeList[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}
	delete(animeList, id)
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

// Utility functions
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}