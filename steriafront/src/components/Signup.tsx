/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: Signup.tsx
Description: Modern, professional signup component for SteriaFront. Handles user registration with validation and error display. Uses modular SCSS.
*/

import React, { useState } from 'react';
import styles from './Signup.module.scss';

interface SignupProps {
  onSignup?: (email: string) => void;
}

const Signup: React.FC<SignupProps> = ({ onSignup = () => {} }) => {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  // Email validation regex
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (!username || !email || !password || !confirm) {
      setError('All fields are required.');
      return;
    }
    if (!emailRegex.test(email)) {
      setError('Please enter a valid email address.');
      return;
    }
    if (password.length < 6) {
      setError('Password must be at least 6 characters.');
      return;
    }
    if (password !== confirm) {
      setError('Passwords do not match.');
      return;
    }
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
      setSuccess(true);
      onSignup(email);
    }, 1200);
  };

  return (
    <div className={styles.signupWrapper}>
      <form className={styles.signupForm} onSubmit={handleSubmit} autoComplete="off">
        <div className={styles.brand}>
          <span className={styles.logo}>SteriaFront</span>
          <span className={styles.subtitle}>Create your account</span>
        </div>
        <div className={styles.inputGroup}>
          <label htmlFor="username">Username</label>
          <input
            id="username"
            type="text"
            value={username}
            onChange={e => setUsername(e.target.value)}
            autoComplete="username"
            disabled={loading}
            required
          />
        </div>
        <div className={styles.inputGroup}>
          <label htmlFor="email">Email</label>
          <input
            id="email"
            type="email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            autoComplete="email"
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
            autoComplete="new-password"
            disabled={loading}
            required
          />
        </div>
        <div className={styles.inputGroup}>
          <label htmlFor="confirm">Confirm Password</label>
          <input
            id="confirm"
            type="password"
            value={confirm}
            onChange={e => setConfirm(e.target.value)}
            autoComplete="new-password"
            disabled={loading}
            required
          />
        </div>
        {error && <div className={styles.error}>{error}</div>}
        {success && <div className={styles.success}>Account created! You can now log in.</div>}
        <button className={styles.signupButton} type="submit" disabled={loading}>
          {loading ? 'Signing up…' : 'Sign Up'}
        </button>
      </form>
    </div>
  );
};

export default Signup; 