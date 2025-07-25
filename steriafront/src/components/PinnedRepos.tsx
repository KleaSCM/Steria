/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: PinnedRepos.tsx
Description: Pinned repositories for SteriaFront profile. Shows 6 pinned repo cards. Modular SCSS, modern, beautiful, femme.
*/

import React from 'react';
import styles from './PinnedRepos.module.scss';

const repos = [
  { name: 'steriafront', desc: 'A modern GitHub clone', lang: 'TypeScript', color: '#3178c6' },
  { name: 'zig-cpp-magic', desc: 'C++ with Zig build system', lang: 'Zig', color: '#f7a41d' },
  { name: 'femmeshell', desc: 'A girly Linux shell', lang: 'Linux', color: '#333' },
  { name: 'ai-notebooks', desc: 'Jupyter AI notebooks', lang: 'Python', color: '#3776ab' },
  { name: 'go-micro', desc: 'Go microservices toolkit', lang: 'Go', color: '#00ADD8' },
  { name: 'react-ui-kit', desc: 'React UI kit for developers', lang: 'React', color: '#61dafb' },
];

const PinnedRepos: React.FC = () => (
  <div className={styles.pinnedWrapper}>
    {repos.map(repo => (
      <div key={repo.name} className={styles.repoCard}>
        <div className={styles.repoName}>{repo.name}</div>
        <div className={styles.repoDesc}>{repo.desc}</div>
        <div className={styles.repoMeta}>
          <span className={styles.langDot} style={{ background: repo.color }} />
          <span className={styles.lang}>{repo.lang}</span>
        </div>
      </div>
    ))}
  </div>
);

export default PinnedRepos;
