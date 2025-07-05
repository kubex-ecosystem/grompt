import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { 
  X, 
  Download, 
  Loader,
  CheckSquare,
  Square
} from 'lucide-react';

const ExportAgentsModal = ({ isOpen, onClose, agents, darkMode }) => {
  const { t } = useTranslation();
  
  const [selectedAgents, setSelectedAgents] = useState(new Set());
  const [format, setFormat] = useState('markdown');
  const [filename, setFilename] = useState('agents');
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
      success: 'bg-green-900 border-green-700 text-green-100'
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
      success: 'bg-green-50 border-green-300 text-green-800'
    }
  };

  const currentTheme = darkMode ? theme.dark : theme.light;

  useEffect(() => {
    if (isOpen) {
      // Select all agents by default
      setSelectedAgents(new Set(agents.map(a => a.ID)));
    }
  }, [isOpen, agents]);

  const handleSelectAll = () => {
    setSelectedAgents(new Set(agents.map(a => a.ID)));
  };

  const handleDeselectAll = () => {
    setSelectedAgents(new Set());
  };

  const handleToggleAgent = (agentId) => {
    const newSelected = new Set(selectedAgents);
    if (newSelected.has(agentId)) {
      newSelected.delete(agentId);
    } else {
      newSelected.add(agentId);
    }
    setSelectedAgents(newSelected);
  };

  const getFileExtension = () => {
    switch (format) {
      case 'json':
        return '.json';
      case 'yaml':
        return '.yaml';
      default:
        return '.md';
    }
  };

  const handleExport = async () => {
    if (selectedAgents.size === 0) return;
    
    setLoading(true);
    
    try {
      const response = await fetch('/api/agents/export-advanced', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          format: format,
          agent_ids: Array.from(selectedAgents),
          filename: filename
        })
      });
      
      if (response.ok) {
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.style.display = 'none';
        a.href = url;
        a.download = filename + getFileExtension();
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
        
        // Close modal after successful export
        handleClose();
      } else {
        throw new Error('Export failed');
      }
    } catch (error) {
      console.error('Export error:', error);
      // Could show error message to user
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setSelectedAgents(new Set());
    setFormat('markdown');
    setFilename('agents');
    setLoading(false);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className={`${currentTheme.modal} rounded-lg shadow-xl w-full max-w-3xl max-h-[90vh] overflow-hidden`}>
        {/* Header */}
        <div className={`flex items-center justify-between p-6 ${currentTheme.border} border-b`}>
          <h2 className={`text-xl font-bold ${currentTheme.text}`}>
            {t('agents.export.title')}
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
            {/* Agent Selection */}
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h3 className={`text-lg font-semibold ${currentTheme.text}`}>
                  {t('agents.export.selectAgents')}
                </h3>
                <div className="flex space-x-2">
                  <button
                    onClick={handleSelectAll}
                    className={`px-3 py-1 text-sm rounded ${currentTheme.buttonSecondary}`}
                  >
                    {t('agents.export.selectAll')}
                  </button>
                  <button
                    onClick={handleDeselectAll}
                    className={`px-3 py-1 text-sm rounded ${currentTheme.buttonSecondary}`}
                  >
                    {t('agents.export.deselectAll')}
                  </button>
                </div>
              </div>

              <div className="space-y-2 max-h-64 overflow-y-auto">
                {agents.map((agent) => (
                  <div
                    key={agent.ID}
                    className={`flex items-center space-x-3 p-3 rounded border ${currentTheme.border} ${currentTheme.modal} cursor-pointer hover:bg-opacity-50`}
                    onClick={() => handleToggleAgent(agent.ID)}
                  >
                    {selectedAgents.has(agent.ID) ? (
                      <CheckSquare size={20} className="text-blue-600" />
                    ) : (
                      <Square size={20} className={currentTheme.textSecondary} />
                    )}
                    <div className="flex-1">
                      <h4 className={`font-medium ${currentTheme.text}`}>{agent.Title}</h4>
                      {agent.Role && (
                        <p className={`text-sm ${currentTheme.textSecondary}`}>{agent.Role}</p>
                      )}
                      {agent.Skills && agent.Skills.length > 0 && (
                        <div className="flex flex-wrap gap-1 mt-1">
                          {agent.Skills.slice(0, 3).map((skill, index) => (
                            <span key={index} className={`px-2 py-0.5 text-xs rounded ${currentTheme.buttonSecondary}`}>
                              {skill}
                            </span>
                          ))}
                          {agent.Skills.length > 3 && (
                            <span className={`px-2 py-0.5 text-xs ${currentTheme.textSecondary}`}>
                              +{agent.Skills.length - 3}
                            </span>
                          )}
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>

              <div className={`text-sm ${currentTheme.textSecondary}`}>
                {selectedAgents.size} / {agents.length} {t('agents.export.selectedAgents')}
              </div>
            </div>

            {/* Export Options */}
            <div className="space-y-4">
              <h3 className={`text-lg font-semibold ${currentTheme.text}`}>
                {t('agents.export.format')}
              </h3>
              
              <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                {Object.entries({
                  markdown: t('agents.export.formats.markdown'),
                  json: t('agents.export.formats.json'),
                  yaml: t('agents.export.formats.yaml')
                }).map(([formatKey, label]) => (
                  <label key={formatKey} className="flex items-center space-x-3 cursor-pointer">
                    <input
                      type="radio"
                      name="format"
                      value={formatKey}
                      checked={format === formatKey}
                      onChange={(e) => setFormat(e.target.value)}
                      className="w-4 h-4"
                    />
                    <span className={currentTheme.text}>{label}</span>
                  </label>
                ))}
              </div>
            </div>

            {/* Filename */}
            <div className="space-y-2">
              <label className={`block text-sm font-medium ${currentTheme.text}`}>
                {t('agents.export.filename')}
              </label>
              <div className="flex items-center space-x-2">
                <input
                  type="text"
                  value={filename}
                  onChange={(e) => setFilename(e.target.value)}
                  className={`flex-1 px-3 py-2 rounded-lg border ${currentTheme.input}`}
                  placeholder="agents"
                />
                <span className={`px-3 py-2 rounded-lg ${currentTheme.buttonSecondary} text-sm`}>
                  {getFileExtension()}
                </span>
              </div>
            </div>
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
            onClick={handleExport}
            disabled={selectedAgents.size === 0 || loading || !filename.trim()}
            className={`flex items-center space-x-2 px-4 py-2 rounded-lg ${currentTheme.button} disabled:opacity-50`}
          >
            {loading ? (
              <Loader size={16} className="animate-spin" />
            ) : (
              <Download size={16} />
            )}
            <span>{t('agents.export.download')}</span>
          </button>
        </div>
      </div>
    </div>
  );
};

export default ExportAgentsModal;
