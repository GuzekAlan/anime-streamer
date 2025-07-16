import React, { useState, useEffect } from 'react';
import { Anime } from '../types';

interface Props {
  animeList: Anime[];
  onDelete: (id: string) => void;
  onPlay: (anime: Anime) => void;
  getProgress: (id: string) => Promise<any>;
}

export const AnimeList: React.FC<Props> = ({ animeList, onDelete, onPlay, getProgress }) => {
  const [progressData, setProgressData] = useState<Record<string, any>>({});

  useEffect(() => {
    const updateProgress = async () => {
      const newProgressData: Record<string, any> = {};
      
      for (const anime of animeList) {
        if (anime.status === 'downloading' || anime.status === 'converting') {
          const progress = await getProgress(anime.id);
          if (progress) {
            newProgressData[anime.id] = progress;
          }
        }
      }
      
      setProgressData(newProgressData);
    };

    updateProgress();
    const interval = setInterval(updateProgress, 2000);
    return () => clearInterval(interval);
  }, [animeList, getProgress]);

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'ready': return '#4CAF50';
      case 'downloading': return '#2196F3';
      case 'converting': return '#FF9800';
      case 'error': return '#F44336';
      default: return '#757575';
    }
  };

  return (
    <div className="anime-list">
      <h2>Anime Library</h2>
      {animeList.length === 0 ? (
        <p>No anime added yet. Add some using the form above!</p>
      ) : (
        <div className="anime-grid">
          {animeList.map((anime) => {
            const currentProgress = progressData[anime.id] || { progress: anime.progress };
            
            return (
              <div key={anime.id} className="anime-card">
                <div className="anime-header">
                  <h3>{anime.name}</h3>
                  <button 
                    className="delete-btn"
                    onClick={() => onDelete(anime.id)}
                    title="Delete anime"
                  >
                    ×
                  </button>
                </div>
                
                <div className="anime-status">
                  <span 
                    className="status-badge"
                    style={{ backgroundColor: getStatusColor(anime.status) }}
                  >
                    {anime.status}
                  </span>
                </div>

                {(anime.status === 'downloading' || anime.status === 'converting') && (
                  <div className="progress-bar">
                    <div 
                      className="progress-fill"
                      style={{ width: `${currentProgress.progress}%` }}
                    />
                    <span className="progress-text">{currentProgress.progress}%</span>
                  </div>
                )}

                {anime.status === 'ready' && anime.qualities && (
                  <div className="qualities">
                    <span>Available: {anime.qualities.join(', ')}</span>
                  </div>
                )}

                <div className="anime-actions">
                  {anime.status === 'ready' ? (
                    <button 
                      className="play-btn"
                      onClick={() => onPlay(anime)}
                    >
                      Play
                    </button>
                  ) : (
                    <button className="play-btn" disabled>
                      {anime.status === 'downloading' ? 'Downloading...' : 
                       anime.status === 'converting' ? 'Converting...' : 'Not Ready'}
                    </button>
                  )}
                </div>

                <div className="anime-meta">
                  <small>Added: {anime.created_at}</small>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};