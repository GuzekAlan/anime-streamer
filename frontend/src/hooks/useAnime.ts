import { useState, useEffect } from 'react';
import { Anime, AddAnimeRequest } from '../types';

const API_BASE = '/api';

export const useAnime = () => {
  const [animeList, setAnimeList] = useState<Anime[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAnimeList = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE}/anime`);
      const data = await response.json();
      setAnimeList(data.anime || []);
    } catch (err) {
      setError('Failed to fetch anime list');
    } finally {
      setLoading(false);
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
    fetchAnimeList();
  }, []);

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