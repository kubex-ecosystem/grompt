/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./{components,contexts,hooks,services,types}/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'monospace'],
      },
    },
  },
  plugins: [],
}