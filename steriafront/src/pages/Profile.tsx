/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: Profile.tsx
Description: Main user profile page for SteriaFront. Responsive, centered, and modern. Imports Navbar, Footer, and all modular profile components.
*/

import React from 'react';
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import ProfileHeader from '../components/ProfileHeader';
import StatusInput from '../components/StatusInput';
import ProfileLinks from '../components/ProfileLinks';
import Badges from '../components/Badges';
import PinnedRepos from '../components/PinnedRepos';
import ProfileBio from '../components/ProfileBio';

const Profile: React.FC = () => {
  return (
    <>
      <Navbar />
      <main style={{ minHeight: '100vh', paddingTop: 80, paddingBottom: 80, background: 'transparent' }}>
        <div
          style={{
            maxWidth: 1100,
            width: 'min(90vw, 1100px)',
            margin: '0 auto',
            padding: '2rem 1.5rem',
            display: 'flex',
            flexDirection: 'column',
            gap: '2.5rem',
            textAlign: 'center',
          }}
        >
          <ProfileHeader />
          <StatusInput />
          <ProfileLinks />
          <Badges />
          <PinnedRepos />
          <ProfileBio />
        </div>
      </main>
      <Footer />
    </>
  );
};

export default Profile; 