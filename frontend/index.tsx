import * as React from 'react';
import * as ReactDOM from 'react-dom/client';
// FIX: Corrected import path for App component
import App from './src/App';

const swExceptions = [
  '//ai.studio',
  'scf.usercontent.goog',
  'generativelanguage.googleapis.com'
]

// Register the service worker for PWA capabilities
if ('serviceWorker' in navigator && swExceptions.filter((v, i) => ((window.location || {}).origin || '').indexOf(v) < 0).length == 0) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.ts')
      .then(registration => {
        console.log('ServiceWorker registration successful with scope: ', registration.scope);
      })
      .catch(err => {
        console.log('ServiceWorker registration failed: ', err);
      });
  });
}

const rootElement = document.getElementById('root');
if (!rootElement) {
  throw new Error("Could not find root element to mount to");
}

const root = ReactDOM.createRoot(rootElement);
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
