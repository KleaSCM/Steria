/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: ProfileHeader.tsx
Description: Profile header for SteriaFront. Shows profile picture, user name, and editable title. Uses modular SCSS.
*/

import React, { useRef, useState } from 'react';
import styles from './ProfileHeader.module.scss';
import { Link } from 'react-router-dom';
import SpotifyPlayer from './SpotifyPlayer';

const DEFAULT_AVATAR = '/0c87901c947d7cd5a097d770eeefd1de.jpg';

const ProfileHeader: React.FC = () => {
  const [avatar, setAvatar] = useState(DEFAULT_AVATAR);
  const name = 'Sylvanas';
  const [title, setTitle] = useState('Full-Stack Software Engineer | Specializing in Typescript, Go, C++, and Cross-Platform Development');
  const [editingTitle, setEditingTitle] = useState(false);
  const fileInput = useRef<HTMLInputElement>(null);

  // Handle avatar upload
  const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const reader = new FileReader();
      reader.onload = ev => setAvatar(ev.target?.result as string);
      reader.readAsDataURL(e.target.files[0]);
    }
  };

  return (
    <div className={styles.header}>
      <div className={styles.topRow}>
        <div className={styles.avatarWrapper}>
          <img src={avatar} alt="Profile avatar" className={styles.avatar} />
          <button
            className={styles.avatarEdit}
            onClick={() => fileInput.current?.click()}
            aria-label="Change profile picture"
          >
            ✨
          </button>
          <input
            type="file"
            accept="image/*"
            ref={fileInput}
            style={{ display: 'none' }}
            onChange={handleAvatarChange}
          />
        </div>
        <div className={styles.spotifyPlayerWrapper}>
          <SpotifyPlayer />
        </div>
      </div>
      <div className={styles.info}>
        <h1 className={styles.name}>
          {name}
          <Link to="/repositories" className={styles.repoBadge} title="View repositories">
            Repositories
          </Link>
          <Link to="/projects" className={styles.repoBadge} title="View projects" style={{ marginLeft: '0.5rem' }}>
            Projects
          </Link>
        </h1>
        {editingTitle ? (
          <input
            className={styles.titleInput}
            value={title}
            onChange={e => setTitle(e.target.value)}
            onBlur={() => setEditingTitle(false)}
            autoFocus
          />
        ) : (
          <h2 className={styles.title} onClick={() => setEditingTitle(true)} title="Click to edit">
            {title}
          </h2>
        )}
      </div>
    </div>
  );
};

export default ProfileHeader;
