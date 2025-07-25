/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: Badges.tsx
Description: Tech/language badges for SteriaFront profile. Modular SCSS, modern, beautiful, femme.
*/

import React from 'react';
import styles from './Badges.module.scss';

const badges = [
  { label: 'TypeScript', color: '#3178c6' },
  { label: 'Go', color: '#00ADD8' },
  { label: 'C++', color: '#00599C' },
  { label: 'React', color: '#61dafb' },
  { label: 'Zig', color: '#f7a41d' },
  { label: 'Python', color: '#3776ab' },
  { label: 'Linux', color: '#333' },
  { label: 'Docker', color: '#2496ed' },
];

const Badges: React.FC = () => (
  <div className={styles.badgesWrapper}>
    {badges.map(badge => (
      <span
        key={badge.label}
        className={styles.badge}
        style={{ background: badge.color }}
      >
        {badge.label}
      </span>
    ))}
  </div>
);

export default Badges;
