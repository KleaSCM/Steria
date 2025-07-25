/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: ProjectList.tsx
Description: List of projects for SteriaFront. Each project is a clickable card with details. Modular SCSS.
*/

import React from 'react';
import styles from './ProjectList.module.scss';

const projects = [
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
    name: 'ai-notebooks',
    tech: ['Python', 'Jupyter'],
    date: '2024-03-20',
    desc: 'Jupyter AI notebooks',
    url: '/projects/ai-notebooks',
  },
];

const ProjectList: React.FC = () => (
  <div className={styles.projectListWrapper}>
    {projects.map(project => (
      <a key={project.name} href={project.url} className={styles.projectCard}>
        <div className={styles.projectName}>{project.name}</div>
        <div className={styles.projectDesc}>{project.desc}</div>
        <div className={styles.projectMeta}>
          <span className={styles.projectDate}>{project.date}</span>
          <span className={styles.techList}>{project.tech.join(', ')}</span>
        </div>
      </a>
    ))}
  </div>
);

export default ProjectList; 