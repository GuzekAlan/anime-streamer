package main

import (
	"log"
	"os"

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