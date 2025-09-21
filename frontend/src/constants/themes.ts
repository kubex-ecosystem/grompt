/**
 * Theme definitions for Grompt UI
 * TypeScript version with proper type definitions
 */

export interface Theme {
  bg: string;
  cardBg: string;
  text: string;
  textSecondary: string;
  border: string;
  input: string;
  button: string;
  buttonSecondary: string;
  accent: string;
}

export interface Themes {
  dark: Theme;
  light: Theme;
}

export const themes: Themes = {
  dark: {
    bg: 'bg-gray-900',
    cardBg: 'bg-gray-800',
    text: 'text-gray-100',
    textSecondary: 'text-gray-300',
    border: 'border-gray-700',
    input: 'bg-gray-700 border-gray-600 text-gray-100',
    button: 'bg-blue-600 hover:bg-blue-700 text-white',
    buttonSecondary: 'bg-gray-700 hover:bg-gray-600 text-gray-200',
    accent: 'text-blue-400'
  },
  light: {
    bg: 'bg-gray-50',
    cardBg: 'bg-white',
    text: 'text-gray-900',
    textSecondary: 'text-gray-600',
    border: 'border-gray-300',
    input: 'bg-white border-gray-300 text-gray-900',
    button: 'bg-blue-600 hover:bg-blue-700 text-white',
    buttonSecondary: 'bg-gray-200 hover:bg-gray-300 text-gray-700',
    accent: 'text-blue-600'
  }
};