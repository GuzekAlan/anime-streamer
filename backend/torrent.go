package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
)

var (
	sharedTorrentClient *torrent.Client
	clientMutex         sync.Mutex
	activeTorrents      = make(map[string]*torrent.Torrent)
)

// getSharedTorrentClient returns a shared torrent client instance
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

// downloadTorrent downloads a torrent and automatically starts HLS conversion
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

// cleanupTorrent removes a torrent from the active torrents list
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