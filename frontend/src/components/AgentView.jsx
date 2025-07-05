import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { 
  ArrowLeft, 
  Edit, 
  Trash2, 
  Sun, 
  Moon, 
  User, 
  Code,
  Shield,
  FileText,
  Copy,
  Download
} from 'lucide-react';
import LanguageSelector from './LanguageSelector';

const AgentView = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { id } = useParams();

  // Estados
  const [darkMode, setDarkMode] = useState(true);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [agent, setAgent] = useState(null);
  const [copied, setCopied] = useState(false);

  // Temas
  const theme = {
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

  const currentTheme = darkMode ? theme.dark : theme.light;

  // Carregar agent
  useEffect(() => {
    const loadAgent = async () => {
      try {
        setLoading(true);
        const response = await fetch(`/api/agents/${id}`);
        if (!response.ok) throw new Error('Failed to load agent');
        const data = await response.json();
        setAgent(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    loadAgent();
  }, [id]);

  // Deletar agent
  const deleteAgent = async () => {
    if (!window.confirm(t('agents.confirmDelete'))) return;
    
    try {
      const response = await fetch(`/api/agents/${id}`, { method: 'DELETE' });
      if (!response.ok) throw new Error('Failed to delete agent');
      navigate('/agents');
    } catch (err) {
      setError(err.message);
    }
  };

  // Copiar prompt example
  const copyPromptExample = async () => {
    try {
      await navigator.clipboard.writeText(agent.PromptExample);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      setError('Failed to copy to clipboard');
    }
  };

  // Exportar agent como markdown
  const exportAsMarkdown = () => {
    const markdown = `# ${agent.Title}

## Role
${agent.Role || 'No role specified'}

## Skills
${agent.Skills?.map(skill => `- ${skill}`).join('\n') || 'No skills specified'}

## Restrictions
${agent.Restrictions?.map(restriction => `- ${restriction}`).join('\n') || 'No restrictions specified'}

## Prompt Example
\`\`\`
${agent.PromptExample || 'No prompt example provided'}
\`\`\`
`;
    
    const blob = new Blob([markdown], { type: 'text/markdown' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${agent.Title.replace(/[^a-zA-Z0-9]/g, '_')}.md`;
    a.click();
    URL.revokeObjectURL(url);
  };

  if (loading) {
    return (
      <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text} flex items-center justify-center`}>
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p>{t('loading')}...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text} flex items-center justify-center`}>
        <div className="text-center">
          <p className="text-red-500 mb-4">{error}</p>
          <button
            onClick={() => navigate('/agents')}
            className={`px-4 py-2 rounded-lg ${currentTheme.button} transition-colors`}
          >
            {t('agents.backToList')}
          </button>
        </div>
      </div>
    );
  }

  if (!agent) {
    return (
      <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text} flex items-center justify-center`}>
        <div className="text-center">
          <p className="mb-4">{t('agents.notFound')}</p>
          <button
            onClick={() => navigate('/agents')}
            className={`px-4 py-2 rounded-lg ${currentTheme.button} transition-colors`}
          >
            {t('agents.backToList')}
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text}`}>
      {/* Header */}
      <header className={`${currentTheme.cardBg} ${currentTheme.border} border-b sticky top-0 z-10`}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center gap-4">
              <button
                onClick={() => navigate('/agents')}
                className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                <ArrowLeft className="h-5 w-5" />
              </button>
              <User className="h-8 w-8 text-blue-500" />
              <h1 className="text-2xl font-bold">{agent.Title}</h1>
            </div>
            
            <div className="flex items-center gap-4">
              <LanguageSelector currentTheme={currentTheme} />
              
              <button
                onClick={() => setDarkMode(!darkMode)}
                className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                {darkMode ? <Sun size={20} /> : <Moon size={20} />}
              </button>
              
              <button
                onClick={() => navigate(`/agents/${id}/edit`)}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg ${currentTheme.button} transition-colors`}
              >
                <Edit size={16} />
                {t('edit')}
              </button>
              
              <button
                onClick={deleteAgent}
                className="flex items-center gap-2 px-4 py-2 rounded-lg bg-red-600 hover:bg-red-700 text-white transition-colors"
              >
                <Trash2 size={16} />
                {t('delete')}
              </button>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Informações Básicas */}
        <div className={`${currentTheme.cardBg} ${currentTheme.border} border rounded-lg p-6 mb-8`}>
          <h2 className="text-xl font-semibold mb-6 flex items-center gap-2">
            <User size={20} className="text-blue-500" />
            {t('agents.form.basic')}
          </h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium mb-2">
                {t('agents.form.title')}
              </label>
              <p className={`${currentTheme.textSecondary} text-lg`}>{agent.Title}</p>
            </div>
            
            {agent.Role && (
              <div>
                <label className="block text-sm font-medium mb-2">
                  {t('agents.form.role')}
                </label>
                <p className={`${currentTheme.textSecondary}`}>{agent.Role}</p>
              </div>
            )}
          </div>
        </div>

        {/* Skills */}
        {agent.Skills && agent.Skills.length > 0 && (
          <div className={`${currentTheme.cardBg} ${currentTheme.border} border rounded-lg p-6 mb-8`}>
            <h2 className="text-xl font-semibold mb-6 flex items-center gap-2">
              <Code size={20} className="text-green-500" />
              {t('agents.form.skills')}
            </h2>
            
            <div className="flex flex-wrap gap-2">
              {agent.Skills.map((skill, index) => (
                <span
                  key={index}
                  className="px-3 py-1 bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 rounded-full text-sm"
                >
                  {skill}
                </span>
              ))}
            </div>
          </div>
        )}

        {/* Restrictions */}
        {agent.Restrictions && agent.Restrictions.length > 0 && (
          <div className={`${currentTheme.cardBg} ${currentTheme.border} border rounded-lg p-6 mb-8`}>
            <h2 className="text-xl font-semibold mb-6 flex items-center gap-2">
              <Shield size={20} className="text-red-500" />
              {t('agents.form.restrictions')}
            </h2>
            
            <div className="flex flex-wrap gap-2">
              {agent.Restrictions.map((restriction, index) => (
                <span
                  key={index}
                  className="px-3 py-1 bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200 rounded-full text-sm"
                >
                  {restriction}
                </span>
              ))}
            </div>
          </div>
        )}

        {/* Prompt Example */}
        {agent.PromptExample && (
          <div className={`${currentTheme.cardBg} ${currentTheme.border} border rounded-lg p-6 mb-8`}>
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-xl font-semibold flex items-center gap-2">
                <FileText size={20} className="text-purple-500" />
                {t('agents.form.promptExample')}
              </h2>
              
              <div className="flex gap-2">
                <button
                  onClick={copyPromptExample}
                  className={`flex items-center gap-2 px-3 py-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
                >
                  <Copy size={16} />
                  {copied ? t('copied') : t('copy')}
                </button>
                
                <button
                  onClick={exportAsMarkdown}
                  className={`flex items-center gap-2 px-3 py-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
                >
                  <Download size={16} />
                  {t('export')}
                </button>
              </div>
            </div>
            
            <pre className={`${currentTheme.input} ${currentTheme.border} border rounded-lg p-4 whitespace-pre-wrap text-sm overflow-auto max-h-96`}>
              {agent.PromptExample}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
};

export default AgentView;
