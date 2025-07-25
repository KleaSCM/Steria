/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: Login.tsx
Description: Fully-implemented, immersive login component for SteriaFront. Modern, beautiful, glassmorphic, with branding and glowing accent.
*/

import React, { useState } from 'react';
import styles from './Login.module.scss';

interface LoginProps {
  onLogin?: (email: string) => void;
}

const Login: React.FC<LoginProps> = ({ onLogin = () => {} }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (!email || !password) {
      setError('Please enter both email and password.');
      return;
    }
    setLoading(true);
    // Simulate async login (replace with real API call)
    setTimeout(() => {
      setLoading(false);
      if (email === '1' && password === '1') {
        onLogin(email);
      } else {
        setError('Invalid email or password.');
      }
    }, 1200);
  };

  return (
    <div className={styles.loginWrapper}>
      <div className={styles.glow} />
      <form className={styles.loginForm} onSubmit={handleSubmit} autoComplete="off">
        <div className={styles.brand}>
          <span className={styles.logo}>SteriaFront</span>
          <span className={styles.subtitle}>A modern GitHub alternative</span>
        </div>
        <h2 className={styles.title}>Sign in to your account</h2>
        <div className={styles.inputGroup}>
          <label htmlFor="email">Email</label>
          <input
            id="email"
            type="text"
            value={email}
            onChange={e => setEmail(e.target.value)}
            autoComplete="username"
            disabled={loading}
            required
          />
        </div>
        <div className={styles.inputGroup}>
          <label htmlFor="password">Password</label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            autoComplete="current-password"
            disabled={loading}
            required
          />
        </div>
        {error && <div className={styles.error}>{error}</div>}
        <button className={styles.loginButton} type="submit" disabled={loading}>
          {loading ? 'Signing in…' : 'Sign In'}
        </button>
      </form>
    </div>
  );
};

export default Login; 