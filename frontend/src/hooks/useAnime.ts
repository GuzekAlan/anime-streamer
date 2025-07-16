import { useState, useEffect, useRef } from 'react';
import { Anime, AddAnimeRequest } from '../types';

const API_BASE = '/api';

export const useAnime = () => {
  const [animeList, setAnimeList] = useState<Anime[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAnimeList = async (showLoading = false) => {
    try {
      if (showLoading) {
        setLoading(true);
      }
      const response = await fetch(`${API_BASE}/anime`);
      const data = await response.json();
      setAnimeList(data.anime || []);
    } catch (err) {
      setError('Failed to fetch anime list');
    } finally {
      if (showLoading) {
        setLoading(false);
      }
    }
  };

  const addAnime = async (request: AddAnimeRequest) => {
    try {
      const response = await fetch(`${API_BASE}/anime`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });
      
      if (!response.ok) {
        throw new Error('Failed to add anime');
      }
      
      await fetchAnimeList();
    } catch (err) {
      setError('Failed to add anime');
      throw err;
    }
  };

  const deleteAnime = async (id: string) => {
    try {
      await fetch(`${API_BASE}/anime/${id}`, {
        method: 'DELETE',
      });
      await fetchAnimeList();
    } catch (err) {
      setError('Failed to delete anime');
    }
  };

  const getProgress = async (id: string) => {
    try {
      const response = await fetch(`${API_BASE}/anime/${id}/progress`);
      return await response.json();
    } catch (err) {
      console.error('Failed to get progress:', err);
      return null;
    }
  };

  useEffect(() => {
    fetchAnimeList(true);
  }, []);

  // Poll when there are downloading or converting anime, and refresh when they complete
  useEffect(() => {
    const hasActiveAnime = animeList.some(anime => 
      anime.status === 'downloading' || anime.status === 'converting'
    );
    
    if (!hasActiveAnime) {
      return; // No active anime, no need to poll
    }

    const checkCompletion = async () => {
      for (const anime of animeList) {
        if (anime.status === 'downloading' || anime.status === 'converting') {
          const progress = await getProgress(anime.id);
          
          // Only refresh if status changed to a different state
          if (progress && progress.status !== anime.status) {
            if (anime.status === 'downloading' && progress.status !== 'downloading') {
              console.log(`Download completed for ${anime.name}, refreshing list`);
              fetchAnimeList();
              return; // Exit after first completion found
            }
            if (anime.status === 'converting' && progress.status === 'ready') {
              console.log(`Conversion completed for ${anime.name}, refreshing list`);
              fetchAnimeList();
              return; // Exit after first completion found
            }
          }
        }
      }
    };

    const interval = setInterval(checkCompletion, 3000);
    return () => clearInterval(interval);
  }, [animeList.filter(anime => anime.status === 'downloading' || anime.status === 'converting').length]); // Only re-run when number of active anime changes

  return {
    animeList,
    loading,
    error,
    addAnime,
    deleteAnime,
    getProgress,
    refetch: fetchAnimeList,
  };
};