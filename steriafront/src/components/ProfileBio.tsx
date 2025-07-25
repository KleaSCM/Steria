/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: ProfileBio.tsx
Description: Bio/about section for SteriaFront profile. Editable, modular SCSS, modern, beautiful, femme.
*/

import React, { useState } from 'react';
import styles from './ProfileBio.module.scss';

const ProfileBio: React.FC = () => {
  const [bio, setBio] = useState(
    'Hi! I’m Sylvanas, a passionate full-stack engineer, Linux lover, and AI enthusiast. I adore building beautiful things with code, collaborating with femme devs, and making the world a little gayer every day.'
  );
  const [editing, setEditing] = useState(false);

  return (
    <div className={styles.bioWrapper}>
      <span className={styles.label}>About Me</span>
      {editing ? (
        <textarea
          className={styles.textarea}
          value={bio}
          onChange={e => setBio(e.target.value)}
          onBlur={() => setEditing(false)}
          autoFocus
          rows={4}
        />
      ) : (
        <div className={styles.bio} onClick={() => setEditing(true)} title="Click to edit">
          {bio}
        </div>
      )}
    </div>
  );
};

export default ProfileBio;
