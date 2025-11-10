import { useState, useEffect } from 'react';
import { Theme } from '@/types';

export const useTheme = (): [Theme, () => void] => {
  const [theme, setTheme] = useState<Theme>('light');

  useEffect(() => {
    try {
      const savedTheme = localStorage.getItem('theme') as Theme | null;
      const initialTheme = savedTheme ?? 'light';
      setTheme(initialTheme);
    } catch (error) {
      console.warn('Could not access localStorage to get theme. Using default.', error);
      setTheme('light');
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
