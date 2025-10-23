/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx}',
    './components/**/*.{js,ts,jsx,tsx}',
    './types/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        display: ['Orbitron', 'Inter', 'sans-serif'],
        mono: ['"IBM Plex Mono"', 'monospace'],
      },
      colors: {
        background: '#f9fafb',
        surface: '#ffffff',
        border: '#e2e8f0',
        text: {
          head: '#111827',
          body: '#334155',
        },
        primary: {
          DEFAULT: '#06b6d4',
          hover: '#0891b2',
          subtle: '#ecfeff',
          foreground: '#024251',
        },
        accent: {
          DEFAULT: '#a855f7',
          hover: '#9333ea',
          secondary: '#d946ef',
          subtle: '#f5f3ff',
          foreground: '#3b0764',
        },
        success: '#16a34a',
        warning: '#f59e0b',
        danger: '#ef4444',
      },
      boxShadow: {
        'soft-card': '0 18px 40px -28px rgba(15, 23, 42, 0.45)',
      },
    },
  },
  plugins: [],
};
