/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: ProfileLinks.tsx
Description: Profile links for SteriaFront. Renders buttons for Kaggle, OSF, LinkedIn, etc. Uses modular SCSS.
*/

import React from 'react';
import styles from './ProfileLinks.module.scss';

const links = [
  { label: 'Kaggle', url: 'https://kaggle.com/', icon: '📊' },
  { label: 'OSF', url: 'https://osf.io/', icon: '🧬' },
  { label: 'LinkedIn', url: 'https://linkedin.com/', icon: '💼' },
  { label: 'GitHub', url: 'https://github.com/', icon: '🐙' },
];

const ProfileLinks: React.FC = () => (
  <div className={styles.linksWrapper}>
    {links.map(link => (
      <a
        key={link.label}
        href={link.url}
        className={styles.linkBtn}
        target="_blank"
        rel="noopener noreferrer"
      >
        <span className={styles.icon}>{link.icon}</span>
        {link.label}
      </a>
    ))}
  </div>
);

export default ProfileLinks;
