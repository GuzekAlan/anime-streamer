export interface Anime {
  id: string;
  name: string;
  torrent_url: string;
  status: 'downloading' | 'converting' | 'ready' | 'error';
  progress: number;
  hls_path?: string;
  hls_urls?: Record<string, string>;  // Quality -> URL mapping
  video_path?: string;
  qualities?: string[];
  created_at: string;
}

export interface AddAnimeRequest {
  name: string;
  torrent_url: string;
  qualities: string[];
}