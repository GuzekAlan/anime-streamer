package main

// Anime represents an anime with download and streaming information
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

// animeList stores all anime in memory
var animeList = make(map[string]*Anime)
