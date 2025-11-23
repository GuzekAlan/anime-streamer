package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getAnimeList returns all anime in the list
func getAnimeList(c *gin.Context) {
	var list []*Anime
	for _, anime := range animeList {
		list = append(list, anime)
	}
	c.JSON(http.StatusOK, gin.H{"anime": list})
}

// getAnime returns a single anime by ID
func getAnime(c *gin.Context) {
	id := c.Param("id")
	anime, exists := animeList[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}
	c.JSON(http.StatusOK, anime)
}

// addAnime adds a new anime and starts downloading
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

// deleteAnime removes an anime and cleans up resources
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

// getDownloadProgress returns the download/conversion progress of an anime
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

// convertToHLS starts HLS conversion for an anime
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
