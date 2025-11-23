package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// convertVideoToHLS converts a video file to HLS format with multiple qualities
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
	anime.HLSPath = fmt.Sprintf("/hls/%s/master.m3u8", anime.ID) // Master playlist
	anime.HLSUrls = hlsUrls                                      // Individual quality URLs

	log.Printf("HLS conversion completed for: %s", anime.Name)
	log.Printf("Available qualities: %v", qualityNames)
	log.Printf("HLS URLs: %v", hlsUrls)
}

// createMasterPlaylist creates an HLS master playlist for adaptive streaming
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
