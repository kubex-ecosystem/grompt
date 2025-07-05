import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { 
  Save, 
  ArrowLeft, 
  Plus, 
  X, 
  Sun, 
  Moon, 
  User, 
  Code,
  Shield,
  FileText
} from 'lucide-react';
import LanguageSelector from './LanguageSelector';

const AgentForm = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { id } = useParams();
  const isEdit = Boolean(id);

  // Estados
  const [darkMode, setDarkMode] = useState(true);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);
  
  // Form data
  const [formData, setFormData] = useState({
    Title: '',
    Role: '',
    Skills: [],
    Restrictions: [],
    PromptExample: ''
  });
  
  const [currentSkill, setCurrentSkill] = useState('');
  const [currentRestriction, setCurrentRestriction] = useState('');

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

  // Carregar agent existente se for edição
  useEffect(() => {
    const loadAgent = async () => {
      try {
        setLoading(true);
        const response = await fetch(`/api/agents/${id}`);
        if (!response.ok) throw new Error('Failed to load agent');
        const agent = await response.json();
        setFormData({
          Title: agent.Title || '',
          Role: agent.Role || '',
          Skills: agent.Skills || [],
          Restrictions: agent.Restrictions || [],
          PromptExample: agent.PromptExample || ''
        });
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    if (isEdit) {
      loadAgent();
    }
  }, [id, isEdit]);

  // Manipular campos do formulário
  const handleInputChange = (field, value) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }));
  };

  // Adicionar skill
  const addSkill = () => {
    if (currentSkill.trim() && !formData.Skills.includes(currentSkill.trim())) {
      setFormData(prev => ({
        ...prev,
        Skills: [...prev.Skills, currentSkill.trim()]
      }));
      setCurrentSkill('');
    }
  };

  // Remover skill
  const removeSkill = (skill) => {
    setFormData(prev => ({
      ...prev,
      Skills: prev.Skills.filter(s => s !== skill)
    }));
  };

  // Adicionar restriction
  const addRestriction = () => {
    if (currentRestriction.trim() && !formData.Restrictions.includes(currentRestriction.trim())) {
      setFormData(prev => ({
        ...prev,
        Restrictions: [...prev.Restrictions, currentRestriction.trim()]
      }));
      setCurrentRestriction('');
    }
  };

  // Remover restriction
  const removeRestriction = (restriction) => {
    setFormData(prev => ({
      ...prev,
      Restrictions: prev.Restrictions.filter(r => r !== restriction)
    }));
  };

  // Submeter formulário
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.Title.trim()) {
      setError('Title is required');
      return;
    }

    try {
      setSaving(true);
      setError(null);
      
      const url = isEdit ? `/api/agents/${id}` : '/api/agents';
      const method = isEdit ? 'PUT' : 'POST';
      
      const response = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData)
      });
      
      if (!response.ok) throw new Error('Failed to save agent');
      
      navigate('/agents');
    } catch (err) {
      setError(err.message);
    } finally {
      setSaving(false);
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
              <button
                onClick={() => navigate('/agents')}
                className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                <ArrowLeft className="h-5 w-5" />
              </button>
              <User className="h-8 w-8 text-blue-500" />
              <h1 className="text-2xl font-bold">
                {isEdit ? t('agents.edit') : t('agents.new')}
              </h1>
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

      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-8">
          {/* Informações Básicas */}
          <div className={`${currentTheme.cardBg} ${currentTheme.border} border rounded-lg p-6`}>
            <h2 className="text-xl font-semibold mb-6 flex items-center gap-2">
              <User size={20} className="text-blue-500" />
              {t('agents.form.basic')}
            </h2>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium mb-2">
                  {t('agents.form.title')} *
                </label>
                <input
                  type="text"
                  value={formData.Title}
                  onChange={(e) => handleInputChange('Title', e.target.value)}
                  placeholder={t('agents.form.titlePlaceholder')}
                  className={`w-full px-4 py-3 rounded-lg ${currentTheme.input} ${currentTheme.border} border focus:ring-2 focus:ring-blue-500 focus:border-transparent`}
                  required
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-2">
                  {t('agents.form.role')}
                </label>
                <input
                  type="text"
                  value={formData.Role}
                  onChange={(e) => handleInputChange('Role', e.target.value)}
                  placeholder={t('agents.form.rolePlaceholder')}
                  className={`w-full px-4 py-3 rounded-lg ${currentTheme.input} ${currentTheme.border} border focus:ring-2 focus:ring-blue-500 focus:border-transparent`}
                />
              </div>
            </div>
          </div>

          {/* Skills */}
          <div className={`${currentTheme.cardBg} ${currentTheme.border} border rounded-lg p-6`}>
            <h2 className="text-xl font-semibold mb-6 flex items-center gap-2">
              <Code size={20} className="text-green-500" />
              {t('agents.form.skills')}
            </h2>
            
            <div className="space-y-4">
              <div className="flex gap-2">
                <input
                  type="text"
                  value={currentSkill}
                  onChange={(e) => setCurrentSkill(e.target.value)}
                  placeholder={t('agents.form.skillPlaceholder')}
                  className={`flex-1 px-4 py-3 rounded-lg ${currentTheme.input} ${currentTheme.border} border focus:ring-2 focus:ring-blue-500 focus:border-transparent`}
                  onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addSkill())}
                />
                <button
                  type="button"
                  onClick={addSkill}
                  className={`px-4 py-3 rounded-lg ${currentTheme.button} transition-colors`}
                >
                  <Plus size={16} />
                </button>
              </div>
              
              {formData.Skills.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {formData.Skills.map((skill, index) => (
                    <span
                      key={index}
                      className="inline-flex items-center gap-2 px-3 py-1 bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 rounded-full text-sm"
                    >
                      {skill}
                      <button
                        type="button"
                        onClick={() => removeSkill(skill)}
                        className="hover:bg-green-200 dark:hover:bg-green-800 rounded-full p-1"
                      >
                        <X size={12} />
                      </button>
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Restrictions */}
          <div className={`${currentTheme.cardBg} ${currentTheme.border} border rounded-lg p-6`}>
            <h2 className="text-xl font-semibold mb-6 flex items-center gap-2">
              <Shield size={20} className="text-red-500" />
              {t('agents.form.restrictions')}
            </h2>
            
            <div className="space-y-4">
              <div className="flex gap-2">
                <input
                  type="text"
                  value={currentRestriction}
                  onChange={(e) => setCurrentRestriction(e.target.value)}
                  placeholder={t('agents.form.restrictionPlaceholder')}
                  className={`flex-1 px-4 py-3 rounded-lg ${currentTheme.input} ${currentTheme.border} border focus:ring-2 focus:ring-blue-500 focus:border-transparent`}
                  onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addRestriction())}
                />
                <button
                  type="button"
                  onClick={addRestriction}
                  className={`px-4 py-3 rounded-lg ${currentTheme.button} transition-colors`}
                >
                  <Plus size={16} />
                </button>
              </div>
              
              {formData.Restrictions.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {formData.Restrictions.map((restriction, index) => (
                    <span
                      key={index}
                      className="inline-flex items-center gap-2 px-3 py-1 bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200 rounded-full text-sm"
                    >
                      {restriction}
                      <button
                        type="button"
                        onClick={() => removeRestriction(restriction)}
                        className="hover:bg-red-200 dark:hover:bg-red-800 rounded-full p-1"
                      >
                        <X size={12} />
                      </button>
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Prompt Example */}
          <div className={`${currentTheme.cardBg} ${currentTheme.border} border rounded-lg p-6`}>
            <h2 className="text-xl font-semibold mb-6 flex items-center gap-2">
              <FileText size={20} className="text-purple-500" />
              {t('agents.form.promptExample')}
            </h2>
            
            <textarea
              value={formData.PromptExample}
              onChange={(e) => handleInputChange('PromptExample', e.target.value)}
              placeholder={t('agents.form.promptExamplePlaceholder')}
              rows={6}
              className={`w-full px-4 py-3 rounded-lg ${currentTheme.input} ${currentTheme.border} border focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none`}
            />
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-4">
            <button
              type="button"
              onClick={() => navigate('/agents')}
              className={`px-6 py-3 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
            >
              {t('cancel')}
            </button>
            <button
              type="submit"
              disabled={saving || !formData.Title.trim()}
              className={`flex items-center gap-2 px-6 py-3 rounded-lg ${currentTheme.button} transition-colors disabled:opacity-50`}
            >
              <Save size={16} />
              {saving ? t('saving') : (isEdit ? t('update') : t('create'))}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AgentForm;
