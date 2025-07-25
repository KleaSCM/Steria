// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: App.tsx
// Description: Main app shell for Steria frontend. Handles routing between Login and Browser pages.
import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Browser from './pages/Browser';
import Home from './pages/Home';

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/browser" element={<Browser />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
