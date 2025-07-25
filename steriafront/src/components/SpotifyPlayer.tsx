/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: SpotifyPlayer.tsx
Description: Modern Spotify player with dropdown song list for SteriaFront profile. Modular SCSS, ready for real integration.
*/

import React, { useState } from 'react';
import styles from './SpotifyPlayer.module.scss';

const songs = [
  { title: 'Levitating', artist: 'Dua Lipa', url: 'https://open.spotify.com/track/463CkQjx2Zk1yXoBuierM9' },
  { title: 'Blinding Lights', artist: 'The Weeknd', url: 'https://open.spotify.com/track/0VjIjW4GlUZAMYd2vXMi3b' },
  { title: 'good 4 u', artist: 'Olivia Rodrigo', url: 'https://open.spotify.com/track/6PERP62TejQjgHu81OHxgX' },
];

const SpotifyPlayer: React.FC = () => {
  const [selected, setSelected] = useState(songs[0]);
  const [open, setOpen] = useState(false);

  return (
    <div className={styles.playerWrapper}>
      <button className={styles.dropdownBtn} onClick={() => setOpen(o => !o)}>
        {selected.title} <span className={styles.artist}>by {selected.artist}</span> ▼
      </button>
      {open && (
        <div className={styles.dropdownList}>
          {songs.map(song => (
            <div
              key={song.title}
              className={styles.dropdownItem}
              onClick={() => { setSelected(song); setOpen(false); }}
            >
              {song.title} <span className={styles.artist}>by {song.artist}</span>
            </div>
          ))}
        </div>
      )}
      <div className={styles.iframeWrapper}>
        <iframe
          src={`https://open.spotify.com/embed/track/${selected.url.split('/').pop()}`}
          width="220"
          height="80"
          frameBorder="0"
          allow="encrypted-media"
          title={selected.title}
        />
      </div>
    </div>
  );
};

export default SpotifyPlayer; 