import { AlertCircle, Check, Eye, EyeOff, Key } from 'lucide-react';
import React, { useEffect, useState } from 'react';
import { useTranslations } from '../../i18n/useTranslations';

interface ApiKeyInputProps {
  onApiKeyChange: (apiKey: string) => void;
  isVisible: boolean;
}

const ApiKeyInput: React.FC<ApiKeyInputProps> = ({ onApiKeyChange, isVisible }) => {
  const { t } = useTranslations();
  const [apiKey, setApiKey] = useState('');
  const [showKey, setShowKey] = useState(false);
  const [isValidating, setIsValidating] = useState(false);
  const [isValid, setIsValid] = useState<boolean | null>(null);

  // Load API key from localStorage on mount
  useEffect(() => {
    const savedKey = localStorage.getItem('userApiKey');
    if (savedKey) {
      setApiKey(savedKey);
      onApiKeyChange(savedKey);
      setIsValid(true);
    }
  }, [onApiKeyChange]);

  const handleApiKeyChange = (value: string) => {
    setApiKey(value);
    setIsValid(null);

    if (value.trim()) {
      // Save to localStorage
      localStorage.setItem('userApiKey', value.trim());
      onApiKeyChange(value.trim());

      // Simple validation - check if it looks like a Gemini API key
      const isValidFormat = value.startsWith('AIza') && value.length > 35;
      setIsValid(isValidFormat);
    } else {
      localStorage.removeItem('userApiKey');
      onApiKeyChange('');
      setIsValid(null);
    }
  };

  const clearApiKey = () => {
    setApiKey('');
    setIsValid(null);
    localStorage.removeItem('userApiKey');
    onApiKeyChange('');
  };

  if (!isVisible) return null;

  return (
    <div className="bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 rounded-lg p-4 mb-6">
      <div className="flex items-center gap-2 mb-3">
        <Key className="h-5 w-5 text-blue-500" />
        <h3 className="font-bold text-lg">Connect Your Gemini API Key</h3>
      </div>

      <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
        To get AI-powered prompts, enter your Gemini API key. Get one free at{' '}
        <a
          href="https://ai.google.dev/"
          target="_blank"
          rel="noopener noreferrer"
          className="text-blue-500 hover:underline"
        >
          ai.google.dev
        </a>
      </p>

      <div className="relative">
        <input
          type={showKey ? 'text' : 'password'}
          value={apiKey}
          onChange={(e) => handleApiKeyChange(e.target.value)}
          placeholder="AIza..."
          className="w-full px-4 py-2 pr-20 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />

        <div className="absolute right-2 top-1/2 transform -translate-y-1/2 flex items-center gap-1">
          {isValid === true && (
            <Check className="h-4 w-4 text-green-500" />
          )}
          {isValid === false && (
            <AlertCircle className="h-4 w-4 text-red-500" />
          )}

          <button
            type="button"
            onClick={() => setShowKey(!showKey)}
            className="p-1 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
          >
            {showKey ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
          </button>
        </div>
      </div>

      {isValid === false && (
        <p className="text-sm text-red-600 dark:text-red-400 mt-2">
          Invalid API key format. Gemini API keys start with "AIza" and are longer than 35 characters.
        </p>
      )}

      {apiKey && (
        <div className="flex justify-between items-center mt-3">
          <p className="text-sm text-green-600 dark:text-green-400">
            ✓ API key saved locally (never sent to our servers)
          </p>
          <button
            onClick={clearApiKey}
            className="text-sm text-red-500 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
          >
            Clear
          </button>
        </div>
      )}
    </div>
  );
};

export default ApiKeyInput;
