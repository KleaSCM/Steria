// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: Login.tsx
// Description: Login page for Steria frontend, ported from Go web LoginHandler. Uses Tailwind CSS for styling.
import React, { useState } from 'react';

const Login: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const res = await fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({ username, password }).toString(),
        credentials: 'include',
      });
      if (res.redirected) {
        window.location.href = res.url;
      } else if (!res.ok) {
        setError('Invalid username or password.');
      }
    } catch {
      setError('Login failed.');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-pink-50">
      <form onSubmit={handleSubmit} className="bg-white p-8 rounded-lg shadow-md w-full max-w-sm">
        <h2 className="text-2xl font-bold mb-6 text-center text-pink-700">Login to Steria</h2>
        <input
          className="w-full mb-4 p-2 border border-pink-200 rounded focus:outline-none focus:ring-2 focus:ring-pink-400"
          name="username"
          placeholder="Username"
          value={username}
          onChange={e => setUsername(e.target.value)}
          required
        />
        <input
          className="w-full mb-4 p-2 border border-pink-200 rounded focus:outline-none focus:ring-2 focus:ring-pink-400"
          name="password"
          type="password"
          placeholder="Password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          required
        />
        {error && <div className="text-red-600 mb-4 text-center">{error}</div>}
        <button
          type="submit"
          className="w-full bg-pink-600 text-white py-2 rounded hover:bg-pink-700 transition"
        >
          Login
        </button>
      </form>
    </div>
  );
};

export default Login; 