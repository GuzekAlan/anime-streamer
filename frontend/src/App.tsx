import React, { useState } from 'react';
import { AddAnimeForm } from './components/AddAnimeForm';
import { AnimeList } from './components/AnimeList';
import { VideoPlayer } from './components/VideoPlayer';
import { Auth } from './components/Auth';
import { useAnime } from './hooks/useAnime';
import { Anime } from './types';
import './App.css';

function App() {
  const { animeList, loading, error, addAnime, deleteAnime, getProgress } = useAnime();
  const [selectedAnime, setSelectedAnime] = useState<Anime | null>(null);

  const handlePlayAnime = (anime: Anime) => {
    setSelectedAnime(anime);
  };

  const handleClosePlayer = () => {
    setSelectedAnime(null);
  };

  return (
    <Auth>
      <div className="App">
        <header className="App-header">
          <h1>🎌 Anime Streaming Platform</h1>
          <p>Download torrents and stream anime with HLS</p>
        </header>

        <main className="App-main">
          {error && (
            <div className="error-message">
              {error}
            </div>
          )}

          <AddAnimeForm onAdd={addAnime} />

          {loading ? (
            <div className="loading">Loading anime list...</div>
          ) : (
            <AnimeList
              animeList={animeList}
              onDelete={deleteAnime}
              onPlay={handlePlayAnime}
              getProgress={getProgress}
            />
          )}
        </main>

        {selectedAnime && (
          <VideoPlayer
            anime={selectedAnime}
            onClose={handleClosePlayer}
          />
        )}
      </div>
    </Auth>
  );
}

export default App;