import React, { useState } from 'react';
import { AddAnimeRequest } from '../types';

interface Props {
  onAdd: (request: AddAnimeRequest) => Promise<void>;
}

export const AddAnimeForm: React.FC<Props> = ({ onAdd }) => {
  const [name, setName] = useState('');
  const [torrentUrl, setTorrentUrl] = useState('');
  const [selectedQualities, setSelectedQualities] = useState<string[]>(['720p', '480p', '360p']);
  const [loading, setLoading] = useState(false);

  const availableQualities = ['720p', '480p', '360p'];

  const handleQualityChange = (quality: string) => {
    setSelectedQualities(prev => 
      prev.includes(quality) 
        ? prev.filter(q => q !== quality)
        : [...prev, quality]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !torrentUrl.trim() || selectedQualities.length === 0) return;

    try {
      setLoading(true);
      await onAdd({ 
        name: name.trim(), 
        torrent_url: torrentUrl.trim(),
        qualities: selectedQualities
      });
      setName('');
      setTorrentUrl('');
      setSelectedQualities(['720p', '480p', '360p']); // Reset to default
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
      
      <div className="form-group">
        <label>Select Qualities to Convert:</label>
        <div className="quality-checkboxes">
          {availableQualities.map(quality => (
            <label key={quality} className="checkbox-label">
              <input
                type="checkbox"
                checked={selectedQualities.includes(quality)}
                onChange={() => handleQualityChange(quality)}
              />
              <span className="checkbox-text">{quality}</span>
            </label>
          ))}
        </div>
        {selectedQualities.length === 0 && (
          <p className="error-text">Please select at least one quality</p>
        )}
      </div>
      
      <button type="submit" disabled={loading || !name.trim() || !torrentUrl.trim() || selectedQualities.length === 0}>
        {loading ? 'Adding...' : 'Add Anime'}
      </button>
    </form>
  );
};