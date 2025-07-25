// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: Home.tsx
// Description: Home/landing page for Steria frontend. Welcoming, with navigation links.
import React from 'react';
import { Link } from 'react-router-dom';

const Home: React.FC = () => (
  <div className="min-h-screen flex flex-col items-center justify-center bg-pink-50">
    <div className="bg-white p-10 rounded-lg shadow-md w-full max-w-lg text-center">
      <h1 className="text-4xl font-bold text-pink-700 mb-4">Welcome to Steria</h1>
      <p className="mb-8 text-pink-900">Your beautiful, modular, sapphic version control system!</p>
      <div className="flex flex-col gap-4">
        <Link to="/login" className="bg-pink-600 text-white py-2 rounded hover:bg-pink-700 transition">Login</Link>
        <Link to="/browser" className="bg-pink-100 text-pink-700 py-2 rounded hover:bg-pink-200 transition">Browse Files</Link>
      </div>
    </div>
  </div>
);

export default Home; 