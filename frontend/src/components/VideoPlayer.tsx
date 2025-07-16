import React, { useRef, useEffect, useState } from 'react';
import Hls from 'hls.js';
import { Anime } from '../types';

interface Props {
  anime: Anime;
  onClose: () => void;
}

export const VideoPlayer: React.FC<Props> = ({ anime, onClose }) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const playerRef = useRef<HTMLDivElement>(null);
  const hlsRef = useRef<Hls | null>(null);
  const [currentQuality, setCurrentQuality] = useState<string>('auto');
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [isFullscreen, setIsFullscreen] = useState(false);

  useEffect(() => {
    const video = videoRef.current;
    if (!video || !anime.hls_path) return;

    if (Hls.isSupported()) {
      const hls = new Hls();
      hlsRef.current = hls;
      
      hls.loadSource(`http://localhost:8080${anime.hls_path}`);
      hls.attachMedia(video);
      
      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        console.log('HLS manifest loaded');
      });

      hls.on(Hls.Events.ERROR, (event, data) => {
        console.error('HLS error:', data);
      });

    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      // Safari native HLS support
      video.src = `http://localhost:8080${anime.hls_path}`;
    }

    return () => {
      if (hlsRef.current) {
        hlsRef.current.destroy();
      }
    };
  }, [anime.hls_path]);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    const handleTimeUpdate = () => setCurrentTime(video.currentTime);
    const handleDurationChange = () => setDuration(video.duration);
    const handlePlay = () => setIsPlaying(true);
    const handlePause = () => setIsPlaying(false);

    video.addEventListener('timeupdate', handleTimeUpdate);
    video.addEventListener('durationchange', handleDurationChange);
    video.addEventListener('play', handlePlay);
    video.addEventListener('pause', handlePause);

    // Fullscreen change listener
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };

    document.addEventListener('fullscreenchange', handleFullscreenChange);

    return () => {
      video.removeEventListener('timeupdate', handleTimeUpdate);
      video.removeEventListener('durationchange', handleDurationChange);
      video.removeEventListener('play', handlePlay);
      video.removeEventListener('pause', handlePause);
      document.removeEventListener('fullscreenchange', handleFullscreenChange);
    };
  }, []);

  const togglePlayPause = () => {
    const video = videoRef.current;
    if (!video) return;

    if (isPlaying) {
      video.pause();
    } else {
      video.play();
    }
  };

  const handleSeek = (e: React.ChangeEvent<HTMLInputElement>) => {
    const video = videoRef.current;
    if (!video) return;

    const newTime = (parseFloat(e.target.value) / 100) * duration;
    video.currentTime = newTime;
  };

  const changeQuality = (quality: string) => {
    const video = videoRef.current;
    if (!video) return;

    let sourceUrl: string;
    
    if (quality === 'auto') {
      // Use master playlist for auto quality
      sourceUrl = `http://localhost:8080${anime.hls_path}`;
    } else {
      // Use specific quality URL
      sourceUrl = anime.hls_urls?.[quality] 
        ? `http://localhost:8080${anime.hls_urls[quality]}`
        : `http://localhost:8080${anime.hls_path}`;
    }

    const currentTime = video.currentTime;
    const wasPlaying = !video.paused;

    if (hlsRef.current) {
      hlsRef.current.destroy();
    }

    if (Hls.isSupported()) {
      const hls = new Hls();
      hlsRef.current = hls;
      hls.loadSource(sourceUrl);
      hls.attachMedia(video);
      
      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        video.currentTime = currentTime;
        if (wasPlaying) {
          video.play();
        }
      });
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      video.src = sourceUrl;
      video.currentTime = currentTime;
      if (wasPlaying) {
        video.play();
      }
    }

    setCurrentQuality(quality);
  };

  const toggleFullscreen = async () => {
    const player = playerRef.current;
    if (!player) return;

    try {
      if (!document.fullscreenElement) {
        await player.requestFullscreen();
      } else {
        await document.exitFullscreen();
      }
    } catch (error) {
      console.error('Fullscreen error:', error);
    }
  };

  const handleDoubleClick = () => {
    toggleFullscreen();
  };

  const formatTime = (time: number) => {
    const hours = Math.floor(time / 3600);
    const minutes = Math.floor((time % 3600) / 60);
    const seconds = Math.floor(time % 60);
    
    if (hours > 0) {
      return `${hours}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    }
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  return (
    <div className="video-player-overlay">
      <div className={`video-player ${isFullscreen ? 'fullscreen' : ''}`} ref={playerRef}>
        <div className="video-header">
          <h3>{anime.name}</h3>
          <button className="close-btn" onClick={onClose}>×</button>
        </div>
        
        <div className="video-container">
          <video
            ref={videoRef}
            controls={false}
            className="video-element"
            onClick={togglePlayPause}
            onDoubleClick={handleDoubleClick}
          />
          
          <div className="video-controls">
            <button className="play-pause-btn" onClick={togglePlayPause}>
              {isPlaying ? '⏸️' : '▶️'}
            </button>
            
            <div className="time-display">
              {formatTime(currentTime)} / {formatTime(duration)}
            </div>
            
            <input
              type="range"
              className="seek-bar"
              min="0"
              max="100"
              value={duration ? (currentTime / duration) * 100 : 0}
              onChange={handleSeek}
            />
            
            <div className="quality-selector">
              <select 
                value={currentQuality} 
                onChange={(e) => changeQuality(e.target.value)}
              >
                <option value="auto">Auto</option>
                {anime.qualities?.map(quality => (
                  <option key={quality} value={quality}>{quality}</option>
                ))}
              </select>
            </div>
            
            <button className="fullscreen-btn" onClick={toggleFullscreen}>
              {isFullscreen ? '⛶' : '⛶'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};