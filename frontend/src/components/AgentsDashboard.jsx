import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { 
  Users, 
  Plus, 
  Edit, 
  Trash2, 
  Download, 
  FileText,
  Search,
  Grid,
  List,
  Sun,
  Moon,
  Copy,
  Eye,
  Upload,
  CheckSquare
} from 'lucide-react';
import LanguageSelector from './LanguageSelector';
import ImportAgentsModal from './ImportAgentsModal';
import ExportAgentsModal from './ExportAgentsModal';
import ValidationAgentsModal from './ValidationAgentsModal';

const AgentsDashboard = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  
  // Estados
  const [agents, setAgents] = useState([]);
  const [darkMode, setDarkMode] = useState(true);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterCategory, setFilterCategory] = useState('all');
  const [viewMode, setViewMode] = useState('grid'); // 'grid' or 'list'
  const [showGenerateModal, setShowGenerateModal] = useState(false);
  const [showImportModal, setShowImportModal] = useState(false);
  const [showExportModal, setShowExportModal] = useState(false);
  const [showValidationModal, setShowValidationModal] = useState(false);
  const [requirements, setRequirements] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [generatedMarkdown, setGeneratedMarkdown] = useState('');
  
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

  // Carregar agents
  useEffect(() => {
    loadAgents();
  }, []);

  const loadAgents = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/agents');
      if (!response.ok) throw new Error('Failed to load agents');
      const data = await response.json();
      setAgents(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Filtrar agents
  const filteredAgents = agents.filter(agent => {
    const matchesSearch = agent.Title?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         agent.Role?.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesFilter = filterCategory === 'all' || 
                         agent.Skills?.some(skill => skill.toLowerCase().includes(filterCategory.toLowerCase()));
    return matchesSearch && matchesFilter;
  });

  // Deletar agent
  const deleteAgent = async (id) => {
    try {
      const response = await fetch(`/api/agents/${id}`, { method: 'DELETE' });
      if (!response.ok) throw new Error('Failed to delete agent');
      setAgents(agents.filter(a => a.ID !== id));
    } catch (err) {
      setError(err.message);
    }
  };

  // Gerar AGENTS.md via LLM
  const generateAgentsFromRequirements = async () => {
    try {
      setIsGenerating(true);
      const response = await fetch('/api/agents/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ requirements })
      });
      
      if (!response.ok) throw new Error('Failed to generate agents');
      const data = await response.json();
      setGeneratedMarkdown(data.markdown);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsGenerating(false);
    }
  };

  // Exportar AGENTS.md
  const exportAgentsMarkdown = async () => {
    try {
      const response = await fetch('/api/agents/markdown');
      if (!response.ok) throw new Error('Failed to export markdown');
      const markdown = await response.text();
      
      const blob = new Blob([markdown], { type: 'text/markdown' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'AGENTS.md';
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      setError(err.message);
    }
  };

  // Handlers para os novos modais
  const handleImportSuccess = (result) => {
    // Recarregar a lista de agents após importação bem-sucedida
    loadAgents();
  };

  const handleExportAgents = () => {
    setShowExportModal(true);
  };

  const handleImportAgents = () => {
    setShowImportModal(true);
  };

  const handleValidateAgents = () => {
    setShowValidationModal(true);
  };

  // Copiar markdown para clipboard
  const copyMarkdownToClipboard = async () => {
    try {
      const response = await fetch('/api/agents/markdown');
      if (!response.ok) throw new Error('Failed to get markdown');
      const markdown = await response.text();
      await navigator.clipboard.writeText(markdown);
    } catch (err) {
      setError(err.message);
    }
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

  return (
    <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text}`}>
      {/* Header */}
      <header className={`${currentTheme.cardBg} ${currentTheme.border} border-b sticky top-0 z-10`}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center gap-4">
              <Users className="h-8 w-8 text-blue-500" />
              <h1 className="text-2xl font-bold">{t('agents.title')}</h1>
            </div>
            
            <div className="flex items-center gap-4">
              <LanguageSelector currentTheme={currentTheme} />
              
              <button
                onClick={() => setDarkMode(!darkMode)}
                className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                {darkMode ? <Sun size={20} /> : <Moon size={20} />}
              </button>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Controles */}
        <div className={`${currentTheme.cardBg} ${currentTheme.border} border rounded-lg p-6 mb-8`}>
          <div className="flex flex-col lg:flex-row gap-4 items-start lg:items-center justify-between">
            {/* Busca e Filtros */}
            <div className="flex flex-col sm:flex-row gap-4 flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <input
                  type="text"
                  placeholder={t('agents.search')}
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className={`pl-10 pr-4 py-2 rounded-lg ${currentTheme.input} ${currentTheme.border} border focus:ring-2 focus:ring-blue-500 focus:border-transparent w-full sm:w-64`}
                />
              </div>
              
              <select
                value={filterCategory}
                onChange={(e) => setFilterCategory(e.target.value)}
                className={`px-4 py-2 rounded-lg ${currentTheme.input} ${currentTheme.border} border focus:ring-2 focus:ring-blue-500 focus:border-transparent`}
              >
                <option value="all">{t('agents.all')}</option>
                <option value="backend">Backend</option>
                <option value="frontend">Frontend</option>
                <option value="devops">DevOps</option>
                <option value="qa">QA</option>
              </select>
            </div>

            {/* Ações */}
            <div className="flex gap-2">
              <div className="flex gap-1 p-1 bg-gray-100 dark:bg-gray-700 rounded-lg">
                <button
                  onClick={() => setViewMode('grid')}
                  className={`p-2 rounded ${viewMode === 'grid' ? 'bg-blue-500 text-white' : currentTheme.buttonSecondary}`}
                >
                  <Grid size={16} />
                </button>
                <button
                  onClick={() => setViewMode('list')}
                  className={`p-2 rounded ${viewMode === 'list' ? 'bg-blue-500 text-white' : currentTheme.buttonSecondary}`}
                >
                  <List size={16} />
                </button>
              </div>

              <button
                onClick={() => setShowGenerateModal(true)}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg ${currentTheme.button} transition-colors`}
              >
                <FileText size={16} />
                {t('agents.generateSquad')}
              </button>

              <button
                onClick={() => navigate('/agents/new')}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg ${currentTheme.button} transition-colors`}
              >
                <Plus size={16} />
                {t('agents.new')}
              </button>

              <button
                onClick={handleImportAgents}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                <Upload size={16} />
                {t('agents.import.title')}
              </button>

              <button
                onClick={handleExportAgents}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                <Download size={16} />
                {t('agents.export.title')}
              </button>

              <button
                onClick={handleValidateAgents}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                <CheckSquare size={16} />
                {t('agents.validation.title')}
              </button>

              <button
                onClick={exportAgentsMarkdown}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                <Download size={16} />
                {t('agents.exportButton')}
              </button>

              <button
                onClick={copyMarkdownToClipboard}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                <Copy size={16} />
                {t('copy')}
              </button>
            </div>
          </div>
        </div>

        {/* Lista de Agents */}
        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        {filteredAgents.length === 0 ? (
          <div className={`${currentTheme.cardBg} ${currentTheme.border} border rounded-lg p-12 text-center`}>
            <Users className="h-24 w-24 text-gray-400 mx-auto mb-4" />
            <h3 className="text-xl font-semibold mb-2">{t('agents.empty.title')}</h3>
            <p className={`${currentTheme.textSecondary} mb-6`}>{t('agents.empty.description')}</p>
            <button
              onClick={() => navigate('/agents/new')}
              className={`flex items-center gap-2 px-6 py-3 rounded-lg ${currentTheme.button} transition-colors mx-auto`}
            >
              <Plus size={20} />
              {t('agents.createFirst')}
            </button>
          </div>
        ) : (
          <div className={viewMode === 'grid' ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6' : 'space-y-4'}>
            {filteredAgents.map((agent) => (
              <div
                key={agent.ID}
                className={`${currentTheme.cardBg} ${currentTheme.border} border rounded-lg p-6 transition-all hover:shadow-lg ${
                  viewMode === 'list' ? 'flex items-center justify-between' : ''
                }`}
              >
                <div className={viewMode === 'list' ? 'flex-1' : ''}>
                  <h3 className="text-xl font-semibold mb-2">{agent.Title}</h3>
                  <p className={`${currentTheme.textSecondary} mb-3`}>{agent.Role}</p>
                  
                  {agent.Skills && agent.Skills.length > 0 && (
                    <div className="flex flex-wrap gap-2 mb-4">
                      {agent.Skills.map((skill, index) => (
                        <span
                          key={index}
                          className="px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-xs rounded-full"
                        >
                          {skill}
                        </span>
                      ))}
                    </div>
                  )}
                </div>

                <div className={`flex gap-2 ${viewMode === 'list' ? 'ml-4' : 'mt-4'}`}>
                  <button
                    onClick={() => navigate(`/agents/${agent.ID}`)}
                    className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
                    title={t('view')}
                  >
                    <Eye size={16} />
                  </button>
                  <button
                    onClick={() => navigate(`/agents/${agent.ID}/edit`)}
                    className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
                    title={t('edit')}
                  >
                    <Edit size={16} />
                  </button>
                  <button
                    onClick={() => deleteAgent(agent.ID)}
                    className="p-2 rounded-lg bg-red-600 hover:bg-red-700 text-white transition-colors"
                    title={t('delete')}
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Modal de Geração */}
      {showGenerateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
          <div className={`${currentTheme.cardBg} rounded-lg p-6 w-full max-w-2xl max-h-[80vh] overflow-y-auto`}>
            <h2 className="text-2xl font-bold mb-4">{t('agents.generate.title')}</h2>
            
            <div className="mb-4">
              <label className="block text-sm font-medium mb-2">
                {t('agents.generate.requirements')}
              </label>
              <textarea
                value={requirements}
                onChange={(e) => setRequirements(e.target.value)}
                placeholder={t('agents.generate.placeholder')}
                className={`w-full h-32 px-4 py-2 rounded-lg ${currentTheme.input} ${currentTheme.border} border focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none`}
              />
            </div>

            {generatedMarkdown && (
              <div className="mb-4">
                <label className="block text-sm font-medium mb-2">
                  {t('agents.generate.result')}
                </label>
                <pre className={`w-full h-48 px-4 py-2 rounded-lg ${currentTheme.input} ${currentTheme.border} border overflow-auto text-sm whitespace-pre-wrap`}>
                  {generatedMarkdown}
                </pre>
              </div>
            )}

            <div className="flex gap-3 justify-end">
              <button
                onClick={() => {
                  setShowGenerateModal(false);
                  setRequirements('');
                  setGeneratedMarkdown('');
                }}
                className={`px-4 py-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                {t('cancel')}
              </button>
              <button
                onClick={generateAgentsFromRequirements}
                disabled={!requirements.trim() || isGenerating}
                className={`px-4 py-2 rounded-lg ${currentTheme.button} transition-colors disabled:opacity-50`}
              >
                {isGenerating ? t('generating') : t('generate')}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Modais */}
      <ImportAgentsModal 
        isOpen={showImportModal}
        onClose={() => setShowImportModal(false)}
        onImport={handleImportSuccess}
        darkMode={darkMode}
      />
      
      <ExportAgentsModal 
        isOpen={showExportModal}
        onClose={() => setShowExportModal(false)}
        agents={agents}
        darkMode={darkMode}
      />
      
      <ValidationAgentsModal 
        isOpen={showValidationModal}
        onClose={() => setShowValidationModal(false)}
        darkMode={darkMode}
      />
    </div>
  );
};

export default AgentsDashboard;
