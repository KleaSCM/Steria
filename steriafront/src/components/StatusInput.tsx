/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: StatusInput.tsx
Description: Editable status input for SteriaFront profile. Lets user set 'What are you working on?'. Uses modular SCSS.
*/

import React, { useState } from 'react';
import styles from './StatusInput.module.scss';

const StatusInput: React.FC = () => {
  const [status, setStatus] = useState('Building something innovative and beautiful!');
  const [editing, setEditing] = useState(false);

  return (
    <div className={styles.statusWrapper}>
      <span className={styles.label}>What are you working on?</span>
      {editing ? (
        <input
          className={styles.input}
          value={status}
          onChange={e => setStatus(e.target.value)}
          onBlur={() => setEditing(false)}
          autoFocus
        />
      ) : (
        <span className={styles.status} onClick={() => setEditing(true)} title="Click to edit">
          {status}
        </span>
      )}
    </div>
  );
};

export default StatusInput;
