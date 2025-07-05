import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { 
  X, 
  FileText, 
  AlertCircle, 
  CheckCircle, 
  Loader,
  AlertTriangle
} from 'lucide-react';

const ValidationAgentsModal = ({ isOpen, onClose, darkMode }) => {
  const { t } = useTranslation();
  
  const [content, setContent] = useState('');
  const [validationResult, setValidationResult] = useState(null);
  const [loading, setLoading] = useState(false);
  
  const theme = {
    dark: {
      bg: 'bg-gray-900',
      modal: 'bg-gray-800',
      text: 'text-gray-100',
      textSecondary: 'text-gray-300',
      border: 'border-gray-700',
      input: 'bg-gray-700 border-gray-600 text-gray-100',
      button: 'bg-blue-600 hover:bg-blue-700 text-white',
      buttonSecondary: 'bg-gray-700 hover:bg-gray-600 text-gray-200',
      success: 'bg-green-900 border-green-700 text-green-100',
      error: 'bg-red-900 border-red-700 text-red-100',
      warning: 'bg-yellow-900 border-yellow-700 text-yellow-100'
    },
    light: {
      bg: 'bg-gray-50',
      modal: 'bg-white',
      text: 'text-gray-900',
      textSecondary: 'text-gray-600',
      border: 'border-gray-300',
      input: 'bg-white border-gray-300 text-gray-900',
      button: 'bg-blue-600 hover:bg-blue-700 text-white',
      buttonSecondary: 'bg-gray-200 hover:bg-gray-300 text-gray-700',
      success: 'bg-green-50 border-green-300 text-green-800',
      error: 'bg-red-50 border-red-300 text-red-800',
      warning: 'bg-yellow-50 border-yellow-300 text-yellow-800'
    }
  };

  const currentTheme = darkMode ? theme.dark : theme.light;

  const handleValidate = async () => {
    if (!content.trim()) return;
    
    setLoading(true);
    setValidationResult(null);
    
    try {
      const response = await fetch('/api/agents/validate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ content: content.trim() })
      });
      
      const data = await response.json();
      
      if (response.ok) {
        setValidationResult(data);
      } else {
        setValidationResult({
          valid: false,
          agents_found: 0,
          errors: [{ type: 'error', message: data.error || 'Validation failed' }],
          error_count: 1
        });
      }
    } catch (error) {
      setValidationResult({
        valid: false,
        agents_found: 0,
        errors: [{ type: 'error', message: 'Failed to validate content' }],
        error_count: 1
      });
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setContent('');
    setValidationResult(null);
    setLoading(false);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className={`${currentTheme.modal} rounded-lg shadow-xl w-full max-w-4xl max-h-[90vh] overflow-hidden`}>
        {/* Header */}
        <div className={`flex items-center justify-between p-6 ${currentTheme.border} border-b`}>
          <h2 className={`text-xl font-bold ${currentTheme.text}`}>
            {t('agents.validation.title')}
          </h2>
          <button
            onClick={handleClose}
            className={`p-2 rounded-lg ${currentTheme.buttonSecondary}`}
          >
            <X size={20} />
          </button>
        </div>

        <div className="overflow-y-auto max-h-[calc(90vh-140px)]">
          <div className="p-6 space-y-6">
            {/* Content Input */}
            <div className="space-y-4">
              <label className={`block text-sm font-medium ${currentTheme.text}`}>
                {t('agents.validation.content')}
              </label>
              <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="# Agents&#10;&#10;## Agent 1: Title&#10;- Role: ...&#10;- Skills: ...&#10;"
                className={`w-full h-64 px-3 py-2 rounded-lg border ${currentTheme.input} resize-none font-mono text-sm`}
              />
            </div>

            {/* Validation Results */}
            {loading && (
              <div className="text-center py-8">
                <Loader size={48} className={`mx-auto mb-4 animate-spin ${currentTheme.text}`} />
                <p className={`text-lg ${currentTheme.text}`}>
                  {t('agents.validation.validating')}
                </p>
              </div>
            )}

            {validationResult && !loading && (
              <div className="space-y-4">
                {/* Summary */}
                <div className={`flex items-center space-x-2 p-4 rounded-lg ${
                  validationResult.valid ? currentTheme.success : currentTheme.error
                }`}>
                  {validationResult.valid ? (
                    <CheckCircle size={20} />
                  ) : (
                    <AlertCircle size={20} />
                  )}
                  <div className="flex-1">
                    <span className="font-medium">
                      {validationResult.valid ? t('agents.validation.valid') : t('agents.validation.invalid')}
                    </span>
                    <div className="text-sm mt-1">
                      {t('agents.validation.summary')}: {validationResult.agents_found} {t('agents.import.agentsFound')}
                      {validationResult.error_count > 0 && (
                        <span>, {validationResult.error_count} {t('agents.validation.errors')}</span>
                      )}
                    </div>
                  </div>
                </div>

                {/* Errors */}
                {validationResult.errors && validationResult.errors.length > 0 && (
                  <div className="space-y-2">
                    <h4 className={`font-semibold ${currentTheme.text}`}>
                      {t('agents.validation.errors')} ({validationResult.errors.length})
                    </h4>
                    <div className="space-y-2 max-h-60 overflow-y-auto">
                      {validationResult.errors.map((error, index) => (
                        <div key={index} className={`flex items-start space-x-2 p-3 rounded ${currentTheme.error}`}>
                          <AlertTriangle size={16} className="mt-0.5 flex-shrink-0" />
                          <div>
                            <p className="font-medium">{error.section || 'General'}</p>
                            <p className="text-sm">{error.message}</p>
                            {error.line && (
                              <p className="text-xs opacity-75">Line {error.line}</p>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Valid Agents Preview */}
                {validationResult.agents && validationResult.agents.length > 0 && (
                  <div className="space-y-2">
                    <h4 className={`font-semibold ${currentTheme.text}`}>
                      Valid Agents ({validationResult.agents.length})
                    </h4>
                    <div className="space-y-2 max-h-60 overflow-y-auto">
                      {validationResult.agents.map((agent, index) => (
                        <div key={index} className={`p-3 rounded border ${currentTheme.border} ${currentTheme.modal}`}>
                          <h5 className={`font-medium ${currentTheme.text}`}>{agent.Title}</h5>
                          {agent.Role && (
                            <p className={`text-sm ${currentTheme.textSecondary}`}>{agent.Role}</p>
                          )}
                          {agent.Skills && agent.Skills.length > 0 && (
                            <div className="flex flex-wrap gap-1 mt-2">
                              {agent.Skills.slice(0, 5).map((skill, skillIndex) => (
                                <span key={skillIndex} className={`px-2 py-1 text-xs rounded ${currentTheme.buttonSecondary}`}>
                                  {skill}
                                </span>
                              ))}
                              {agent.Skills.length > 5 && (
                                <span className={`px-2 py-1 text-xs rounded ${currentTheme.textSecondary}`}>
                                  +{agent.Skills.length - 5} more
                                </span>
                              )}
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Footer */}
        <div className={`flex justify-end space-x-3 p-6 ${currentTheme.border} border-t`}>
          <button
            onClick={handleClose}
            className={`px-4 py-2 rounded-lg ${currentTheme.buttonSecondary}`}
            disabled={loading}
          >
            {t('ideas.cancel')}
          </button>
          
          <button
            onClick={handleValidate}
            disabled={!content.trim() || loading}
            className={`flex items-center space-x-2 px-4 py-2 rounded-lg ${currentTheme.button} disabled:opacity-50`}
          >
            {loading ? (
              <Loader size={16} className="animate-spin" />
            ) : (
              <FileText size={16} />
            )}
            <span>{t('agents.validation.validate')}</span>
          </button>
        </div>
      </div>
    </div>
  );
};

export default ValidationAgentsModal;
