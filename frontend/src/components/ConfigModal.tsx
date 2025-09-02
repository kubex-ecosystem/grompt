'use client';

import { useEffect, useState } from 'react';
import { FaTimes } from 'react-icons/fa';

interface ConfigModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: () => void;
}

interface Config {
  gemini_api_key: string;
  openai_api_key: string;
  claude_api_key: string;
  deepseek_api_key: string;
  chatgpt_api_key: string;
  ollama_endpoint: string;
}

export default function ConfigModal({ isOpen, onClose, onSave }: ConfigModalProps) {
  const [config, setConfig] = useState<Config>({
    gemini_api_key: '',
    openai_api_key: '',
    claude_api_key: '',
    deepseek_api_key: '',
    chatgpt_api_key: '',
    ollama_endpoint: 'http://localhost:11434'
  });
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (isOpen) {
      loadConfig();
    }
  }, [isOpen]);

  const loadConfig = async () => {
    try {
      const response = await fetch('/config');
      if (response.ok) {
        const data = await response.json();
        setConfig(data);
      }
    } catch (error) {
      console.error('Erro ao carregar configuração:', error);
    }
  };

  const handleSave = async () => {
    setIsLoading(true);
    try {
      const response = await fetch('/config', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(config),
      });

      if (response.ok) {
        onSave();
        onClose();
      } else {
        alert('Erro ao salvar configuração');
      }
    } catch (error) {
      console.error('Erro ao salvar configuração:', error);
      alert('Erro ao salvar configuração');
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-gray-800 rounded-lg p-6 w-full max-w-2xl mx-4 max-h-[90vh] overflow-y-auto">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-xl font-bold text-white">⚙️ Configuração de APIs</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white"
            title="Fechar"
          >
            <FaTimes size={24} />
          </button>
        </div>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              🔶 Gemini API Key
            </label>
            <input
              type="password"
              value={config.gemini_api_key}
              onChange={(e) => setConfig({ ...config, gemini_api_key: e.target.value })}
              className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="AIzaSy..."
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              🤖 OpenAI API Key
            </label>
            <input
              type="password"
              value={config.openai_api_key}
              onChange={(e) => setConfig({ ...config, openai_api_key: e.target.value })}
              className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="sk-..."
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              🤖 ChatGPT API Key
            </label>
            <input
              type="password"
              value={config.chatgpt_api_key}
              onChange={(e) => setConfig({ ...config, chatgpt_api_key: e.target.value })}
              className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="sk-..."
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              🧠 Claude API Key
            </label>
            <input
              type="password"
              value={config.claude_api_key}
              onChange={(e) => setConfig({ ...config, claude_api_key: e.target.value })}
              className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="sk-ant-..."
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              🚀 DeepSeek API Key
            </label>
            <input
              type="password"
              value={config.deepseek_api_key}
              onChange={(e) => setConfig({ ...config, deepseek_api_key: e.target.value })}
              className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="sk-..."
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              🦙 Ollama Endpoint
            </label>
            <input
              type="text"
              value={config.ollama_endpoint}
              onChange={(e) => setConfig({ ...config, ollama_endpoint: e.target.value })}
              className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="http://localhost:11434"
            />
          </div>
        </div>

        <div className="flex justify-end space-x-3 mt-6">
          <button
            onClick={onClose}
            className="px-4 py-2 text-gray-400 hover:text-white transition-colors"
          >
            Cancelar
          </button>
          <button
            onClick={handleSave}
            disabled={isLoading}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 transition-colors"
          >
            {isLoading ? 'Salvando...' : 'Salvar'}
          </button>
        </div>
      </div>
    </div>
  );
}
