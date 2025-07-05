import React, { useState, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { 
  X, 
  Upload, 
  FileText, 
  AlertCircle, 
  CheckCircle, 
  Loader,
  AlertTriangle
} from 'lucide-react';

const ImportAgentsModal = ({ isOpen, onClose, onImport, darkMode }) => {
  const { t } = useTranslation();
  const fileInputRef = useRef(null);
  
  const [content, setContent] = useState('');
  const [mode, setMode] = useState('upload'); // 'upload' or 'paste'
  const [options, setOptions] = useState({
    merge: false,
    validate: true
  });
  const [previewData, setPreviewData] = useState(null);
  const [errors, setErrors] = useState([]);
  const [loading, setLoading] = useState(false);
  const [step, setStep] = useState('input'); // 'input', 'preview', 'importing'
  
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

  const handleFileSelect = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        setContent(e.target.result);
      };
      reader.readAsText(file);
    }
  };

  const validateContent = async () => {
    if (!content.trim()) return;
    
    setLoading(true);
    setErrors([]);
    
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
        setPreviewData(data);
        setErrors(data.errors || []);
        setStep('preview');
      } else {
        setErrors([{ type: 'error', message: data.error || 'Validation failed' }]);
      }
    } catch (error) {
      setErrors([{ type: 'error', message: 'Failed to validate content' }]);
    } finally {
      setLoading(false);
    }
  };

  const handleImport = async () => {
    if (!content.trim()) return;
    
    // Confirm replacement if not merging
    if (!options.merge) {
      if (!window.confirm(t('agents.import.confirmReplace'))) {
        return;
      }
    }
    
    setLoading(true);
    setStep('importing');
    
    try {
      const response = await fetch('/api/agents/import', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          content: content.trim(),
          merge: options.merge,
          validate: options.validate
        })
      });
      
      const data = await response.json();
      
      if (response.ok && data.success) {
        // Success - close modal and refresh agents list
        onImport(data);
        handleClose();
      } else {
        // Show errors
        setErrors(data.errors || [{ type: 'error', message: data.message || 'Import failed' }]);
        setStep('preview');
      }
    } catch (error) {
      setErrors([{ type: 'error', message: 'Failed to import agents' }]);
      setStep('preview');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setContent('');
    setPreviewData(null);
    setErrors([]);
    setStep('input');
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
            {t('agents.import.title')}
          </h2>
          <button
            onClick={handleClose}
            className={`p-2 rounded-lg ${currentTheme.buttonSecondary}`}
          >
            <X size={20} />
          </button>
        </div>

        <div className="overflow-y-auto max-h-[calc(90vh-140px)]">
          {step === 'input' && (
            <div className="p-6 space-y-6">
              {/* Mode Selection */}
              <div className="flex space-x-4">
                <button
                  onClick={() => setMode('upload')}
                  className={`flex items-center space-x-2 px-4 py-2 rounded-lg ${
                    mode === 'upload' ? currentTheme.button : currentTheme.buttonSecondary
                  }`}
                >
                  <Upload size={20} />
                  <span>{t('agents.import.upload')}</span>
                </button>
                <button
                  onClick={() => setMode('paste')}
                  className={`flex items-center space-x-2 px-4 py-2 rounded-lg ${
                    mode === 'paste' ? currentTheme.button : currentTheme.buttonSecondary
                  }`}
                >
                  <FileText size={20} />
                  <span>{t('agents.import.paste')}</span>
                </button>
              </div>

              {/* File Upload */}
              {mode === 'upload' && (
                <div className="space-y-4">
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept=".md,.markdown,.txt"
                    onChange={handleFileSelect}
                    className="hidden"
                  />
                  <button
                    onClick={() => fileInputRef.current?.click()}
                    className={`w-full p-8 border-2 border-dashed ${currentTheme.border} rounded-lg hover:bg-opacity-50 transition-colors`}
                  >
                    <Upload size={48} className={`mx-auto mb-4 ${currentTheme.textSecondary}`} />
                    <p className={currentTheme.text}>
                      {t('agents.import.upload')}
                    </p>
                    <p className={`text-sm ${currentTheme.textSecondary} mt-2`}>
                      {t('agents.import.supportedFormats')}
                    </p>
                  </button>
                </div>
              )}

              {/* Paste Content */}
              {mode === 'paste' && (
                <div className="space-y-4">
                  <label className={`block text-sm font-medium ${currentTheme.text}`}>
                    {t('agents.import.content')}
                  </label>
                  <textarea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="# Agents\n\n## Agent 1: Title\n- Role: ...\n- Skills: ...\n"
                    className={`w-full h-64 px-3 py-2 rounded-lg border ${currentTheme.input} resize-none`}
                  />
                </div>
              )}

              {/* Content Preview */}
              {content && (
                <div className="space-y-4">
                  <label className={`block text-sm font-medium ${currentTheme.text}`}>
                    {t('agents.import.preview')}
                  </label>
                  <div className={`p-4 rounded-lg border ${currentTheme.border} ${currentTheme.input} max-h-32 overflow-y-auto`}>
                    <pre className={`text-sm ${currentTheme.textSecondary}`}>
                      {content.substring(0, 500)}
                      {content.length > 500 && '...'}
                    </pre>
                  </div>
                </div>
              )}

              {/* Options */}
              <div className="space-y-4">
                <h3 className={`text-lg font-semibold ${currentTheme.text}`}>
                  {t('agents.import.options')}
                </h3>
                <div className="space-y-3">
                  <label className="flex items-center space-x-3">
                    <input
                      type="radio"
                      name="importMode"
                      checked={options.merge}
                      onChange={() => setOptions(prev => ({ ...prev, merge: true }))}
                      className="w-4 h-4"
                    />
                    <span className={currentTheme.text}>{t('agents.import.merge')}</span>
                  </label>
                  <label className="flex items-center space-x-3">
                    <input
                      type="radio"
                      name="importMode"
                      checked={!options.merge}
                      onChange={() => setOptions(prev => ({ ...prev, merge: false }))}
                      className="w-4 h-4"
                    />
                    <span className={currentTheme.text}>{t('agents.import.replace')}</span>
                  </label>
                  <label className="flex items-center space-x-3">
                    <input
                      type="checkbox"
                      checked={options.validate}
                      onChange={(e) => setOptions(prev => ({ ...prev, validate: e.target.checked }))}
                      className="w-4 h-4"
                    />
                    <span className={currentTheme.text}>{t('agents.import.validate')}</span>
                  </label>
                </div>
              </div>
            </div>
          )}

          {step === 'preview' && previewData && (
            <div className="p-6 space-y-6">
              {/* Validation Results */}
              <div className="space-y-4">
                <div className={`flex items-center space-x-2 p-4 rounded-lg ${
                  previewData.valid ? currentTheme.success : currentTheme.error
                }`}>
                  {previewData.valid ? (
                    <CheckCircle size={20} />
                  ) : (
                    <AlertCircle size={20} />
                  )}
                  <span className="font-medium">
                    {previewData.valid ? t('agents.validation.valid') : t('agents.validation.invalid')}
                  </span>
                  <span>
                    {previewData.agents_found} {t('agents.import.agentsFound')}
                  </span>
                </div>

                {/* Errors */}
                {errors.length > 0 && (
                  <div className="space-y-2">
                    <h4 className={`font-semibold ${currentTheme.text}`}>
                      {t('agents.import.errors')} ({errors.length})
                    </h4>
                    <div className="space-y-2 max-h-40 overflow-y-auto">
                      {errors.map((error, index) => (
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

                {/* Agents Preview */}
                {previewData.agents && previewData.agents.length > 0 && (
                  <div className="space-y-2">
                    <h4 className={`font-semibold ${currentTheme.text}`}>
                      {t('agents.import.preview')} ({previewData.agents.length})
                    </h4>
                    <div className="space-y-2 max-h-60 overflow-y-auto">
                      {previewData.agents.map((agent, index) => (
                        <div key={index} className={`p-3 rounded border ${currentTheme.border} ${currentTheme.modal}`}>
                          <h5 className={`font-medium ${currentTheme.text}`}>{agent.Title}</h5>
                          {agent.Role && (
                            <p className={`text-sm ${currentTheme.textSecondary}`}>{agent.Role}</p>
                          )}
                          {agent.Skills && agent.Skills.length > 0 && (
                            <div className="flex flex-wrap gap-1 mt-2">
                              {agent.Skills.slice(0, 3).map((skill, skillIndex) => (
                                <span key={skillIndex} className={`px-2 py-1 text-xs rounded ${currentTheme.buttonSecondary}`}>
                                  {skill}
                                </span>
                              ))}
                              {agent.Skills.length > 3 && (
                                <span className={`px-2 py-1 text-xs rounded ${currentTheme.textSecondary}`}>
                                  +{agent.Skills.length - 3} more
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
            </div>
          )}

          {step === 'importing' && (
            <div className="p-12 text-center">
              <Loader size={48} className={`mx-auto mb-4 animate-spin ${currentTheme.text}`} />
              <p className={`text-lg ${currentTheme.text}`}>
                {t('agents.import.importing')}
              </p>
            </div>
          )}
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
          
          {step === 'input' && (
            <button
              onClick={validateContent}
              disabled={!content.trim() || loading}
              className={`flex items-center space-x-2 px-4 py-2 rounded-lg ${currentTheme.button} disabled:opacity-50`}
            >
              {loading && <Loader size={16} className="animate-spin" />}
              <span>{t('agents.import.validate')}</span>
            </button>
          )}
          
          {step === 'preview' && (
            <button
              onClick={handleImport}
              disabled={loading || (options.validate && errors.length > 0)}
              className={`flex items-center space-x-2 px-4 py-2 rounded-lg ${currentTheme.button} disabled:opacity-50`}
            >
              {loading && <Loader size={16} className="animate-spin" />}
              <span>{t('agents.import.title')}</span>
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export default ImportAgentsModal;
