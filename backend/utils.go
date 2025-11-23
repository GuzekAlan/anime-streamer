package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// generateID generates a unique ID based on current timestamp
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// getCurrentTime returns the current time formatted as a string
func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// findVideoFileInDir searches for video files in a directory
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

// findVideoFile finds a video file in the torrent download directory
func findVideoFile(torrentName string) string {
	downloadPath := filepath.Join("../storage/downloads", torrentName)
	videoExts := []string{".mp4", ".mkv", ".avi", ".mov", ".wmv"}

	var videoFile string
	filepath.Walk(downloadPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
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

// isFFmpegAvailable checks if FFmpeg is installed and available
func isFFmpegAvailable() bool {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	return err == nil
}

// executeCommand executes a shell command and logs the output
func executeCommand(cmdStr string) error {
	log.Printf("Executing: %s", cmdStr)

	// Parse the command string into parts
	parts := []string{"sh", "-c", cmdStr}
	cmd := exec.Command(parts[0], parts[1:]...)

	// Set up output capture
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command failed: %v, Output: %s", err, string(output))
		return err
	}

	log.Printf("Command completed successfully")
	return nil
}

// saveAnimeMetadata saves anime metadata to a file for future restoration
func saveAnimeMetadata(outputDir string, anime *Anime) {
	metadataFile := filepath.Join(outputDir, "metadata.txt")
	content := fmt.Sprintf("%s\n%s\n%s", anime.Name, anime.TorrentURL, anime.CreatedAt)

	if err := os.WriteFile(metadataFile, []byte(content), 0644); err != nil {
		log.Printf("Failed to save metadata for %s: %v", anime.Name, err)
	} else {
		log.Printf("Saved metadata for %s", anime.Name)
	}
}
