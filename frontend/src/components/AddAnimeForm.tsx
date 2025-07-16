import React, { useState } from 'react';
import { AddAnimeRequest } from '../types';

interface Props {
  onAdd: (request: AddAnimeRequest) => Promise<void>;
}

export const AddAnimeForm: React.FC<Props> = ({ onAdd }) => {
  const [name, setName] = useState('');
  const [torrentUrl, setTorrentUrl] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !torrentUrl.trim()) return;

    try {
      setLoading(true);
      await onAdd({ name: name.trim(), torrent_url: torrentUrl.trim() });
      setName('');
      setTorrentUrl('');
    } catch (err) {
      console.error('Failed to add anime:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="add-anime-form">
      <h2>Add New Anime</h2>
      <div className="form-group">
        <label htmlFor="name">Anime Name:</label>
        <input
          id="name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter anime name"
          required
        />
      </div>
      <div className="form-group">
        <label htmlFor="torrent">Torrent URL/Magnet Link:</label>
        <input
          id="torrent"
          type="text"
          value={torrentUrl}
          onChange={(e) => setTorrentUrl(e.target.value)}
          placeholder="magnet:?xt=urn:btih:..."
          required
        />
      </div>
      <button type="submit" disabled={loading || !name.trim() || !torrentUrl.trim()}>
        {loading ? 'Adding...' : 'Add Anime'}
      </button>
    </form>
  );
};