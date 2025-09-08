import { useState, useEffect } from 'react';
import { Theme } from '../types';

export const useTheme = (): [Theme, () => void] => {
  const [theme, setTheme] = useState<Theme>('dark');

  useEffect(() => {
    try {
      const savedTheme = localStorage.getItem('theme') as Theme | null;
      const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
      const initialTheme = savedTheme || (prefersDark ? 'dark' : 'light');
      setTheme(initialTheme);
    } catch (error) {
      console.warn('Could not access localStorage to get theme. Using default.', error);
      // Fallback to prefers-color-scheme if localStorage is not available
      const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
      setTheme(prefersDark ? 'dark' : 'light');
    }
  }, []);

  useEffect(() => {
    try {
      const root = window.document.documentElement;
      const body = window.document.body;
      root.classList.remove('light', 'dark');
      root.classList.add(theme);
      body.classList.remove('light-theme', 'dark-theme');
      body.classList.add(`${theme}-theme`);
      localStorage.setItem('theme', theme);
    } catch (error) {
      console.warn('Could not access localStorage to save theme.', error);
    }
  }, [theme]);

  const toggleTheme = () => {
    setTheme(prevTheme => (prevTheme === 'light' ? 'dark' : 'light'));
  };

  return [theme, toggleTheme];
};