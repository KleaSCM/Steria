/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: Repositories.tsx
Description: User repositories page for SteriaFront. Professional, modular, and ready for expansion.
*/

import React from 'react';
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import RepositoryList from '../components/RepositoryList';

const Repositories: React.FC = () => {
  return (
    <>
      <Navbar />
      <main style={{ minHeight: '100vh', paddingTop: 80, paddingBottom: 80 }}>
        <div style={{ maxWidth: 1100, margin: '0 auto', padding: '2rem 1.5rem', textAlign: 'center' }}>
          <h1 style={{ fontSize: '2.2rem', fontWeight: 800, color: '#a21caf', marginBottom: '2rem' }}>
            Repositories
          </h1>
          <RepositoryList />
        </div>
      </main>
      <Footer />
    </>
  );
};

export default Repositories; 