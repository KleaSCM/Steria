// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: main.tsx
// Description: Entry point for Steria frontend. Renders the App component.
import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
