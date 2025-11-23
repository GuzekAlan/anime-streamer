package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// scanExistingAnime scans for existing anime files in storage directories
func scanExistingAnime() {
	log.Println("Scanning existing anime files...")

	// Scan HLS directory first (these are ready to stream)
	scanHLSDirectory()

	// Scan downloads directory for completed downloads not yet converted
	scanDownloadsDirectory()

	log.Printf("Found %d existing anime", len(animeList))
}

// scanHLSDirectory scans the HLS directory for already converted anime
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

// scanDownloadsDirectory scans the downloads directory for completed downloads
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

// scanAvailableQualities checks which quality playlists are available
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

// getAnimeNameFromFiles tries to get the anime name from metadata files
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
