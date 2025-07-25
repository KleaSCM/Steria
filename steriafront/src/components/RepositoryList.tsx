/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: RepositoryList.tsx
Description: List of repositories for SteriaFront. Each repo is a clickable card with project details. Modular SCSS.
*/

import React from 'react';
import styles from './RepositoryList.module.scss';

const repos = [
  {
    name: 'steriafront',
    tech: ['TypeScript', 'React', 'SCSS'],
    date: '2024-05-01',
    desc: 'A modern GitHub clone',
    url: '/projects/steriafront',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'zig-cpp-magic',
    tech: ['Zig', 'C++'],
    date: '2024-04-15',
    desc: 'C++ with Zig build system',
    url: '/projects/zig-cpp-magic',
  },
  {
    name: 'ai-notebooks',
    tech: ['Python', 'Jupyter'],
    date: '2024-03-20',
    desc: 'Jupyter AI notebooks',
    url: '/projects/ai-notebooks',
  },
  
];

const RepositoryList: React.FC = () => (
  <div className={styles.repoListWrapper}>
    {repos.map(repo => (
      <a key={repo.name} href={repo.url} className={styles.repoCard}>
        <div className={styles.repoName}>{repo.name}</div>
        <div className={styles.repoDesc}>{repo.desc}</div>
        <div className={styles.repoMeta}>
          <span className={styles.repoDate}>{repo.date}</span>
          <span className={styles.techList}>{repo.tech.join(', ')}</span>
        </div>
      </a>
    ))}
  </div>
);

export default RepositoryList; 