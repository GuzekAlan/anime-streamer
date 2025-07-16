package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
)

var (
	sharedTorrentClient *torrent.Client
	clientMutex         sync.Mutex
	activeTorrents      = make(map[string]*torrent.Torrent)
)

func getSharedTorrentClient() (*torrent.Client, error) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if sharedTorrentClient == nil {
		cfg := torrent.NewDefaultClientConfig()
		cfg.DataDir = "../storage/downloads"
		
		client, err := torrent.NewClient(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create shared torrent client: %v", err)
		}
		
		sharedTorrentClient = client
		log.Printf("Created shared torrent client")
	}
	
	return sharedTorrentClient, nil
}

func downloadTorrent(anime *Anime) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error downloading torrent for %s: %v", anime.Name, r)
			anime.Status = "error"
			// Clean up from active torrents map
			clientMutex.Lock()
			delete(activeTorrents, anime.ID)
			clientMutex.Unlock()
		}
	}()

	// Get shared torrent client
	client, err := getSharedTorrentClient()
	if err != nil {
		log.Printf("Error getting torrent client: %v", err)
		anime.Status = "error"
		return
	}

	// Add torrent
	t, err := client.AddMagnet(anime.TorrentURL)
	if err != nil {
		log.Printf("Error adding torrent: %v", err)
		anime.Status = "error"
		return
	}

	// Store torrent reference for cleanup
	clientMutex.Lock()
	activeTorrents[anime.ID] = t
	clientMutex.Unlock()

	// Wait for torrent info
	<-t.GotInfo()
	log.Printf("Starting download: %s", t.Info().Name)

	// Start downloading
	t.DownloadAll()

	// Monitor progress
	for {
		stats := t.Stats()
		totalLength := t.Length()
		if totalLength > 0 {
			completed := stats.BytesReadData.Int64()
			progress := int((completed * 100) / totalLength)
			anime.Progress = progress
			
			if progress >= 100 {
				anime.Status = "converting"
				videoFile := findVideoFile(t.Info().Name)
				anime.HLSPath = videoFile
				log.Printf("Download completed: %s", anime.Name)
				
				// Clean up from active torrents map
				clientMutex.Lock()
				delete(activeTorrents, anime.ID)
				clientMutex.Unlock()
				
				// Start HLS conversion automatically
				go convertVideoToHLS(anime)
				break
			}
		}
		time.Sleep(2 * time.Second)
	}
}

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

func convertVideoToHLS(anime *Anime) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error converting video for %s: %v", anime.Name, r)
			anime.Status = "error"
		}
	}()

	// Check if FFmpeg is available
	if !isFFmpegAvailable() {
		log.Printf("FFmpeg not available for conversion of %s", anime.Name)
		anime.Status = "error"
		return
	}

	anime.Status = "converting"
	anime.Progress = 0
	
	inputFile := anime.HLSPath
	if inputFile == "" {
		log.Printf("No video file found for %s", anime.Name)
		anime.Status = "error"
		return
	}

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		log.Printf("Video file does not exist: %s", inputFile)
		anime.Status = "error"
		return
	}

	outputDir := filepath.Join("../storage/hls", anime.ID)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("Failed to create output directory: %v", err)
		anime.Status = "error"
		return
	}

	log.Printf("Starting HLS conversion for %s", anime.Name)
	log.Printf("Input file: %s", inputFile)
	log.Printf("Output directory: %s", outputDir)

	// Define all available quality settings
	allQualities := map[string]struct {
		resolution string
		bitrate    string
		crf        string
		preset     string
	}{
		"720p": {"1280x720", "2500k", "23", "medium"},
		"480p": {"854x480", "1200k", "26", "medium"},
		"360p": {"640x360", "600k", "28", "fast"},
	}

	// Filter to only convert selected qualities
	var selectedQualityConfigs []struct {
		name       string
		resolution string
		bitrate    string
		crf        string
		preset     string
	}

	for _, qualityName := range anime.Qualities {
		if config, exists := allQualities[qualityName]; exists {
			selectedQualityConfigs = append(selectedQualityConfigs, struct {
				name       string
				resolution string
				bitrate    string
				crf        string
				preset     string
			}{
				name:       qualityName,
				resolution: config.resolution,
				bitrate:    config.bitrate,
				crf:        config.crf,
				preset:     config.preset,
			})
		}
	}

	qualities := selectedQualityConfigs

	var qualityNames []string
	totalQualities := len(qualities)
	
	for i, quality := range qualities {
		// Update progress based on current quality being processed
		progressPercent := (i * 100) / totalQualities
		anime.Progress = progressPercent
		
		log.Printf("Converting %s to %s (%d/%d - %d%%)", anime.Name, quality.name, i+1, totalQualities, progressPercent)
		
		outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.m3u8", quality.name))
		segmentPattern := filepath.Join(outputDir, fmt.Sprintf("%s%%03d.ts", quality.name))
		
		cmd := fmt.Sprintf(
			"ffmpeg -i \"%s\" -c:v libx264 -preset %s -crf %s -maxrate %s -bufsize %s -c:a aac -b:a 128k -vf scale=%s -hls_time 10 -hls_list_size 0 -hls_segment_filename \"%s\" -f hls \"%s\"",
			inputFile, quality.preset, quality.crf, quality.bitrate, quality.bitrate, quality.resolution, segmentPattern, outputPath,
		)
		
		log.Printf("FFmpeg command: %s", cmd)
		
		// Execute FFmpeg command
		if err := executeCommand(cmd); err != nil {
			log.Printf("Error converting %s: %v", quality.name, err)
			continue
		}
		
		// Verify the output file was created
		if _, err := os.Stat(outputPath); err == nil {
			qualityNames = append(qualityNames, quality.name)
			log.Printf("Successfully created %s playlist", quality.name)
		} else {
			log.Printf("Failed to create %s playlist: %v", quality.name, err)
		}
	}
	
	// Set final progress to 100% when conversion is complete
	anime.Progress = 100

	// Create master playlist
	createMasterPlaylist(outputDir, qualityNames)
	
	// Save anime metadata for future restoration
	saveAnimeMetadata(outputDir, anime)
	
	// Create HLS URLs map for each quality
	hlsUrls := make(map[string]string)
	for _, quality := range qualityNames {
		hlsUrls[quality] = fmt.Sprintf("/hls/%s/%s.m3u8", anime.ID, quality)
	}
	
	anime.Status = "ready"
	anime.Qualities = qualityNames
	anime.HLSPath = fmt.Sprintf("/hls/%s/master.m3u8", anime.ID)  // Master playlist
	anime.HLSUrls = hlsUrls  // Individual quality URLs
	
	log.Printf("HLS conversion completed for: %s", anime.Name)
	log.Printf("Available qualities: %v", qualityNames)
	log.Printf("HLS URLs: %v", hlsUrls)
}

func createMasterPlaylist(outputDir string, qualities []string) {
	masterPath := filepath.Join(outputDir, "master.m3u8")
	file, err := os.Create(masterPath)
	if err != nil {
		log.Printf("Error creating master playlist: %v", err)
		return
	}
	defer file.Close()

	file.WriteString("#EXTM3U\n")
	file.WriteString("#EXT-X-VERSION:3\n\n")

	bandwidths := map[string]string{
		"720p": "2500000",
		"480p": "1000000", 
		"360p": "500000",
	}

	resolutions := map[string]string{
		"720p": "1280x720",
		"480p": "854x480",
		"360p": "640x360",
	}

	for _, quality := range qualities {
		file.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%s,RESOLUTION=%s\n", 
			bandwidths[quality], resolutions[quality]))
		file.WriteString(fmt.Sprintf("%s.m3u8\n", quality))
	}
}

func cleanupTorrent(animeID string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	
	if torrent, exists := activeTorrents[animeID]; exists {
		log.Printf("Cleaning up torrent for anime ID: %s", animeID)
		
		// Drop the torrent from the client
		if sharedTorrentClient != nil {
			torrent.Drop()
		}
		
		// Remove from active torrents map
		delete(activeTorrents, animeID)
		
		log.Printf("Torrent cleanup completed for anime ID: %s", animeID)
	}
}

func saveAnimeMetadata(outputDir string, anime *Anime) {
	metadataFile := filepath.Join(outputDir, "metadata.txt")
	content := fmt.Sprintf("%s\n%s\n%s", anime.Name, anime.TorrentURL, anime.CreatedAt)
	
	if err := os.WriteFile(metadataFile, []byte(content), 0644); err != nil {
		log.Printf("Failed to save metadata for %s: %v", anime.Name, err)
	} else {
		log.Printf("Saved metadata for %s", anime.Name)
	}
}

func isFFmpegAvailable() bool {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	return err == nil
}

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