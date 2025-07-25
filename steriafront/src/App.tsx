/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: App.tsx
Description: Main application entry for SteriaFront. Sets up routing for Login, Signup, and Profile pages, with MeteorShower background, Navbar, and Footer.
*/

import React from 'react';
import { BrowserRouter, Routes, Route, useNavigate, Navigate } from 'react-router-dom';
import MeteorShower from './components/MeteorShower';
import Navbar from './components/Navbar';
import Login from './components/Login';
import Signup from './components/Signup';
import Footer from './components/Footer';
import Profile from './pages/Profile';
import Repositories from './pages/Repositories';
import Projects from './pages/Projects';

// Wrapper to handle login navigation
const LoginWithRedirect: React.FC = () => {
  const navigate = useNavigate();
  return <Login onLogin={() => navigate('/profile')} />;
};

const SignupWithRedirect: React.FC = () => {
  const navigate = useNavigate();
  return <Signup onSignup={() => navigate('/login')} />;
};

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <MeteorShower />
      <Navbar />
      <Routes>
        <Route path="/login" element={<LoginWithRedirect />} />
        <Route path="/signup" element={<SignupWithRedirect />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/repositories" element={<Repositories />} />
        <Route path="/projects" element={<Projects />} />
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
      <Footer />
    </BrowserRouter>
  );
};

export default App;
