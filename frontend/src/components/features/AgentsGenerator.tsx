import {
  BadgeCheck,
  ClipboardCheck,
  ClipboardCopy,
  Download,
  History,
  Loader2,
  Plus,
  RefreshCcw,
  Save,
  Settings,
  Sparkles,
  Trash2,
  X,
} from 'lucide-react';
import React, { useCallback, useContext, useEffect, useMemo, useState } from 'react';
import ReactMarkdown from 'react-markdown';
import type { AgentsGenerationResult, Idea, StoredAgent } from '../../../types';
import { LanguageContext } from '../../context/LanguageContext';
import { configService, type ProviderInfo } from '../../services/configService';
import { agentsService } from '../../services/agentsService';
import Card from '../ui/Card';

type AgentFramework = 'crewai' | 'autogen' | 'langchain' | 'semantic-kernel' | 'custom';
type PurposeKey = 'automation' | 'analysis' | 'support' | 'research' | 'delivery' | 'other';

type BlueprintEntry = {
  id: string;
  createdAt: number;
  requirements: string;
  result: AgentsGenerationResult;
};

const frameworks: { value: AgentFramework; label: string; description: string }[] = [
  { value: 'crewai', label: 'CrewAI', description: 'Orquestra squads multi-agentes focados em produtividade.' },
  { value: 'autogen', label: 'AutoGen', description: 'Fluxos cooperativos entre agentes conversacionais.' },
  { value: 'langchain', label: 'LangChain', description: 'Tooling maduro com LangGraph e integrações robustas.' },
  { value: 'semantic-kernel', label: 'Semantic Kernel', description: 'Pipeline .NET com planejamento dinâmico de tarefas.' },
  { value: 'custom', label: 'Custom', description: 'Adapter próprio seguindo padrões Kubex.' },
];

const toolCatalog: string[] = [
  'web_search',
  'file_handler',
  'calculator',
  'email_sender',
  'database',
  'api_caller',
  'code_executor',
  'image_generator',
  'git_ops',
  'docker_manager',
];

const mcpCatalog: { name: string; desc: string }[] = [
  { name: 'filesystem', desc: '📁 Filesystem' },
  { name: 'database', desc: '🗄️ Database' },
  { name: 'web-scraper', desc: '🕷️ Web Scraper' },
  { name: 'git', desc: '🔄 Git' },
  { name: 'docker', desc: '🐳 Docker' },
  { name: 'kubernetes', desc: '☸️ Kubernetes' },
  { name: 'slack', desc: '💬 Slack' },
  { name: 'github', desc: '🐙 GitHub' },
  { name: 'notion', desc: '📝 Notion' },
  { name: 'calendar', desc: '📅 Calendar' },
];

const copyTimerMs = 2200;

const i18n: Record<string, Record<string, string>> = {
  en: {
    overviewTitle: 'Agents Command Center',
    overviewDescription: 'Configure squads, generate AGENTS.md and govern your AI estate.',
    frameworksLabel: 'Framework',
    providersLabel: 'Provider',
    toolsLabel: 'Tools',
    mcpLabel: 'MCP Servers',
    workspaceTitle: 'Agents Workspace',
    workspaceDescription: 'Capture requirements and configure frameworks before generating squads.',
    contextCardTitle: 'Context & Ideas',
    ideaInputPlaceholder: 'Document a requirement, hypothesis or constraint...',
    addIdea: 'Add idea',
    ideasEmpty: 'Start with raw requirements – they will be translated into agents.',
    configurationTitle: 'Agent Configuration',
    roleLabel: 'Agent role / specialization',
    rolePlaceholder: 'Example: Senior Data Strategist for LATAM Retail',
    purposeLabel: 'Primary purpose',
    purpose_automation: 'Automation',
    purpose_analysis: 'Analysis',
    purpose_support: 'Customer Support',
    purpose_research: 'Research & Discovery',
    purpose_delivery: 'Software Delivery',
    purpose_other: 'Other',
    customPurposePlaceholder: 'Describe a custom purpose...',
    toolsSectionTitle: 'Operational toolkit',
    mcpSectionTitle: 'Model Context Protocol',
    mcpCustomPlaceholder: 'custom-server',
    generationTitle: 'Squad Blueprint',
    generateButton: 'Generate Agents',
    generatingButton: 'Generating agents...',
    openBlueprint: 'Open blueprint',
    requirementsPreview: 'Requirements preview',
    generatedAgents: 'Generated agents',
    noGeneratedAgents: 'Generate a squad to preview the Markdown and export AGENTS.md.',
    copyMarkdown: 'Copy Markdown',
    copyRequirements: 'Copy requirements',
    copyAgentsTable: 'Copy table',
    copied: 'Copied!',
    saveAgents: 'Persist agents',
    saving: 'Saving...',
    savedFeedback: 'Agents persisted to store.',
    saveFailed: 'Unable to persist agents. Review backend logs.',
    exportMarkdown: 'Export AGENTS.md',
    refreshAgents: 'Refresh list',
    storedAgentsTitle: 'Registered Agents',
    storedAgentsDescription: 'Agents persisted on the server (internal JSON store).',
    emptyStoredAgents: 'No agents registered yet.',
    idColumn: 'ID',
    titleColumn: 'Title',
    roleColumn: 'Role',
    skillsColumn: 'Skills',
    actionsColumn: 'Actions',
    deleteAction: 'Delete agent',
    providerUnavailable: 'Unavailable',
    providerNeedsKey: 'Needs API key',
    providerReady: 'Ready',
    generationFailed: 'Generation failed. Please confirm the unified backend is reachable.',
    requirementsMissing: 'Add at least one idea and define provider + role before generating.',
    markdownSection: 'Markdown Preview',
    blueprintHistoryTitle: 'Blueprint history',
    blueprintHistoryEmpty: 'No blueprints generated yet.',
    viewBlueprint: 'View blueprint',
    generatedAt: 'Generated at',
    close: 'Close',
  },
  pt: {
    overviewTitle: 'Central de Agents',
    overviewDescription: 'Configure squads, gere o AGENTS.md e governe seu ecossistema de IA.',
    frameworksLabel: 'Framework',
    providersLabel: 'Provider',
    toolsLabel: 'Ferramentas',
    mcpLabel: 'Servidores MCP',
    workspaceTitle: 'Espaço de Trabalho de Agents',
    workspaceDescription: 'Capture requisitos e ajuste frameworks antes de gerar squads.',
    contextCardTitle: 'Contexto & Ideias',
    ideaInputPlaceholder: 'Anote um requisito, hipótese ou constraint...',
    addIdea: 'Adicionar ideia',
    ideasEmpty: 'Comece registrando requisitos brutos — eles serão convertidos em agents.',
    configurationTitle: 'Configuração do Agent',
    roleLabel: 'Papel / especialização do agent',
    rolePlaceholder: 'Ex: Senior Data Strategist para varejo LATAM',
    purposeLabel: 'Propósito principal',
    purpose_automation: 'Automação',
    purpose_analysis: 'Análise',
    purpose_support: 'Suporte',
    purpose_research: 'Pesquisa',
    purpose_delivery: 'Entrega de Software',
    purpose_other: 'Outro',
    customPurposePlaceholder: 'Descreva um propósito customizado...',
    toolsSectionTitle: 'Ferramentas operacionais',
    mcpSectionTitle: 'Model Context Protocol',
    mcpCustomPlaceholder: 'servidor-personalizado',
    generationTitle: 'Blueprint do Squad',
    generateButton: 'Gerar Agents',
    generatingButton: 'Gerando agents...',
    openBlueprint: 'Abrir blueprint',
    requirementsPreview: 'Prévia das instruções',
    generatedAgents: 'Agents gerados',
    noGeneratedAgents: 'Gere um squad para visualizar o Markdown e exportar o AGENTS.md.',
    copyMarkdown: 'Copiar Markdown',
    copyRequirements: 'Copiar requisitos',
    copyAgentsTable: 'Copiar tabela',
    copied: 'Copiado!',
    saveAgents: 'Persistir agents',
    saving: 'Salvando...',
    savedFeedback: 'Agents persistidos com sucesso.',
    saveFailed: 'Não foi possível salvar os agents. Verifique o backend.',
    exportMarkdown: 'Exportar AGENTS.md',
    refreshAgents: 'Atualizar lista',
    storedAgentsTitle: 'Agents Registrados',
    storedAgentsDescription: 'Agents persistidos no servidor (store interno em JSON).',
    emptyStoredAgents: 'Nenhum agent registrado ainda.',
    idColumn: 'ID',
    titleColumn: 'Título',
    roleColumn: 'Papel',
    skillsColumn: 'Skills',
    actionsColumn: 'Ações',
    deleteAction: 'Remover agent',
    providerUnavailable: 'Indisponível',
    providerNeedsKey: 'Configurar API key',
    providerReady: 'Pronto',
    generationFailed: 'Falha na geração. Confirme se o backend unificado está acessível.',
    requirementsMissing: 'Adicione ao menos uma ideia e defina provider + papel antes de gerar.',
    markdownSection: 'Prévia do Markdown',
    blueprintHistoryTitle: 'Histórico de blueprints',
    blueprintHistoryEmpty: 'Nenhum blueprint gerado ainda.',
    viewBlueprint: 'Ver blueprint',
    generatedAt: 'Gerado em',
    close: 'Fechar',
  },
  es: {
    overviewTitle: 'Centro de Agents',
    overviewDescription: 'Configura squads, genera AGENTS.md y gobierna tu portafolio de IA.',
    frameworksLabel: 'Framework',
    providersLabel: 'Proveedor',
    toolsLabel: 'Herramientas',
    mcpLabel: 'Servidores MCP',
    workspaceTitle: 'Espacio de Trabajo de Agents',
    workspaceDescription: 'Captura requisitos y ajusta frameworks antes de generar squads.',
    contextCardTitle: 'Contexto e Ideas',
    ideaInputPlaceholder: 'Registra un requerimiento, hipótesis o restricción...',
    addIdea: 'Agregar idea',
    ideasEmpty: 'Comienza con requisitos en bruto — serán traducidos a agents.',
    configurationTitle: 'Configuración del Agent',
    roleLabel: 'Rol / especialización del agent',
    rolePlaceholder: 'Ej: Estratega de Datos Senior para retail LATAM',
    purposeLabel: 'Propósito principal',
    purpose_automation: 'Automatización',
    purpose_analysis: 'Análisis',
    purpose_support: 'Soporte al cliente',
    purpose_research: 'Investigación',
    purpose_delivery: 'Entrega de software',
    purpose_other: 'Otro',
    customPurposePlaceholder: 'Describe un propósito personalizado...',
    toolsSectionTitle: 'Kit operacional',
    mcpSectionTitle: 'Model Context Protocol',
    mcpCustomPlaceholder: 'servidor-personalizado',
    generationTitle: 'Blueprint del Squad',
    generateButton: 'Generar Agents',
    generatingButton: 'Generando agents...',
    openBlueprint: 'Abrir blueprint',
    requirementsPreview: 'Previsualización de instrucciones',
    generatedAgents: 'Agents generados',
    noGeneratedAgents: 'Genera un squad para ver el Markdown y exportar AGENTS.md.',
    copyMarkdown: 'Copiar Markdown',
    copyRequirements: 'Copiar requisitos',
    copyAgentsTable: 'Copiar tabla',
    copied: '¡Copiado!',
    saveAgents: 'Persistir agents',
    saving: 'Guardando...',
    savedFeedback: 'Agents guardados correctamente.',
    saveFailed: 'No fue posible guardar los agents. Revisa el backend.',
    exportMarkdown: 'Exportar AGENTS.md',
    refreshAgents: 'Actualizar lista',
    storedAgentsTitle: 'Agents Registrados',
    storedAgentsDescription: 'Agents persistidos en el servidor (store interno JSON).',
    emptyStoredAgents: 'Aún no hay agents registrados.',
    idColumn: 'ID',
    titleColumn: 'Título',
    roleColumn: 'Rol',
    skillsColumn: 'Skills',
    actionsColumn: 'Acciones',
    deleteAction: 'Eliminar agent',
    providerUnavailable: 'No disponible',
    providerNeedsKey: 'Configurar API key',
    providerReady: 'Listo',
    generationFailed: 'Falló la generación. Confirma que el backend unificado está disponible.',
    requirementsMissing: 'Agrega al menos una idea y define provider + rol antes de generar.',
    markdownSection: 'Vista previa Markdown',
    blueprintHistoryTitle: 'Historial de blueprints',
    blueprintHistoryEmpty: 'Aún no se generaron blueprints.',
    viewBlueprint: 'Ver blueprint',
    generatedAt: 'Generado el',
    close: 'Cerrar',
  },
  zh: {
    overviewTitle: 'Agents 控制中心',
    overviewDescription: '配置多代理小队，生成 AGENTS.md 并治理你的 AI 资产。',
    frameworksLabel: '框架',
    providersLabel: '提供方',
    toolsLabel: '工具',
    mcpLabel: 'MCP 服务器',
    workspaceTitle: 'Agents 工作区',
    workspaceDescription: '在生成小队前，先整理需求并配置框架。',
    contextCardTitle: '上下文与想法',
    ideaInputPlaceholder: '记录需求、假设或限制条件...',
    addIdea: '添加想法',
    ideasEmpty: '先记录原始需求——系统会将其转换为 agents。',
    configurationTitle: 'Agent 配置',
    roleLabel: 'Agent 角色 / 专长',
    rolePlaceholder: '示例：拉美零售行业资深数据策略师',
    purposeLabel: '主要目的',
    purpose_automation: '自动化',
    purpose_analysis: '分析',
    purpose_support: '客户支持',
    purpose_research: '研究探索',
    purpose_delivery: '软件交付',
    purpose_other: '其他',
    customPurposePlaceholder: '描述自定义目的...',
    toolsSectionTitle: '操作工具包',
    mcpSectionTitle: 'Model Context Protocol',
    mcpCustomPlaceholder: '自定义服务器',
    generationTitle: '小队蓝图',
    generateButton: '生成 Agents',
    generatingButton: '正在生成 agents...',
    openBlueprint: '打开 blueprint',
    requirementsPreview: '指令预览',
    generatedAgents: '已生成的 agents',
    noGeneratedAgents: '生成一个小队以查看 Markdown 并导出 AGENTS.md。',
    copyMarkdown: '复制 Markdown',
    copyRequirements: '复制需求',
    copyAgentsTable: '复制表格',
    copied: '已复制！',
    saveAgents: '持久化 agents',
    saving: '保存中...',
    savedFeedback: 'Agents 已保存。',
    saveFailed: '无法保存 agents。请检查后端日志。',
    exportMarkdown: '导出 AGENTS.md',
    refreshAgents: '刷新列表',
    storedAgentsTitle: '已注册 Agents',
    storedAgentsDescription: '存储在服务器上的 agents（内部 JSON）。',
    emptyStoredAgents: '暂无注册的 agents。',
    idColumn: 'ID',
    titleColumn: '标题',
    roleColumn: '角色',
    skillsColumn: '技能',
    actionsColumn: '操作',
    deleteAction: '删除 agent',
    providerUnavailable: '不可用',
    providerNeedsKey: '需要 API key',
    providerReady: '已就绪',
    generationFailed: '生成失败。请确认统一后端可访问。',
    requirementsMissing: '在生成前请至少添加一个想法并设置 provider + 角色。',
    markdownSection: 'Markdown 预览',
    blueprintHistoryTitle: 'Blueprint 历史',
    blueprintHistoryEmpty: '尚未生成 blueprint。',
    viewBlueprint: '查看 blueprint',
    generatedAt: '生成时间',
    close: '关闭',
  },
};

const purposes: PurposeKey[] = ['automation', 'analysis', 'support', 'research', 'delivery', 'other'];

const AgentsGenerator: React.FC = () => {
  const { language } = useContext(LanguageContext);
  const [ideas, setIdeas] = useState<Idea[]>([]);
  const [currentIdea, setCurrentIdea] = useState('');
  const [agentFramework, setAgentFramework] = useState<AgentFramework>('crewai');
  const [agentProvider, setAgentProvider] = useState<string>('');
  const [providers, setProviders] = useState<ProviderInfo[]>([]);
  const [providersLoading, setProvidersLoading] = useState<boolean>(false);
  const [providersError, setProvidersError] = useState<string | null>(null);
  const [agentRole, setAgentRole] = useState('');
  const [purpose, setPurpose] = useState<PurposeKey>('automation');
  const [customPurpose, setCustomPurpose] = useState('');
  const [agentTools, setAgentTools] = useState<string[]>([]);
  const [mcpServers, setMcpServers] = useState<string[]>([]);
  const [customMcp, setCustomMcp] = useState('');
  const [storedAgents, setStoredAgents] = useState<StoredAgent[]>([]);
  const [agentsLoading, setAgentsLoading] = useState<boolean>(true);
  const [agentsError, setAgentsError] = useState<string | null>(null);
  const [history, setHistory] = useState<BlueprintEntry[]>([]);
  const [currentBlueprint, setCurrentBlueprint] = useState<BlueprintEntry | null>(null);
  const [isBlueprintModalOpen, setBlueprintModalOpen] = useState(false);
  const [isGenerating, setIsGenerating] = useState(false);
  const [generationError, setGenerationError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [copiedRequirements, setCopiedRequirements] = useState(false);
  const [copiedAgentsTable, setCopiedAgentsTable] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [saveStatus, setSaveStatus] = useState<{ variant: 'success' | 'error'; message: string } | null>(null);
  const t = useCallback(
    (key: string) => i18n[language]?.[key] ?? i18n.en[key] ?? key,
    [language],
  );

  useEffect(() => {
    if (copied || copiedAgentsTable || copiedRequirements) {
      const timer = window.setTimeout(() => {
        setCopied(false);
        setCopiedAgentsTable(false);
        setCopiedRequirements(false);
      }, copyTimerMs);
      return () => clearTimeout(timer);
    }
    return undefined;
  }, [copied, copiedAgentsTable, copiedRequirements]);

  const loadProviders = useCallback(async () => {
    setProvidersLoading(true);
    setProvidersError(null);
    try {
      const config = await configService.getConfig(false);
      const available = config.available_providers
        .map((name) => config.providers[name])
        .filter((provider): provider is ProviderInfo => Boolean(provider));
      setProviders(available);
      if (!agentProvider && available.length > 0) {
        setAgentProvider(available[0].name);
      } else if (config.default_provider && !agentProvider) {
        setAgentProvider(config.default_provider);
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);
      setProvidersError(message);
    } finally {
      setProvidersLoading(false);
    }
  }, [agentProvider]);

  const fetchAgents = useCallback(async () => {
    setAgentsLoading(true);
    setAgentsError(null);
    try {
      const data = await agentsService.list();
      setStoredAgents(data);
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);
      setAgentsError(message);
    } finally {
      setAgentsLoading(false);
    }
  }, []);

  useEffect(() => {
    loadProviders();
    fetchAgents();
  }, [fetchAgents, loadProviders]);

  const toggleTool = (tool: string) => {
    setAgentTools((prev) =>
      prev.includes(tool) ? prev.filter((item) => item !== tool) : [...prev, tool],
    );
  };

  const toggleMcp = (server: string) => {
    setMcpServers((prev) =>
      prev.includes(server) ? prev.filter((item) => item !== server) : [...prev, server],
    );
  };

  const removeIdea = (id: string) => {
    setIdeas((prev) => prev.filter((idea) => idea.id !== id));
  };

  const addIdea = () => {
    if (!currentIdea.trim()) return;
    setIdeas((prev) => [
      ...prev,
      {
        id: Date.now().toString(),
        text: currentIdea.trim(),
      },
    ]);
    setCurrentIdea('');
  };

  const composedPurpose = useMemo(() => {
    if (purpose === 'other') {
      return customPurpose.trim() || 'Custom mission';
    }
    return t(`purpose_${purpose}`);
  }, [customPurpose, purpose, t]);

  const requirements = useMemo(() => {
    const lines: string[] = [];
    lines.push(`# Squad Mission: ${composedPurpose}`);
    if (agentRole) {
      lines.push(`Primary Role: ${agentRole}`);
    }
    lines.push(`Framework: ${agentFramework}`);
    if (agentProvider) {
      lines.push(`LLM Provider: ${agentProvider}`);
    }
    lines.push(`Tools: ${agentTools.length > 0 ? agentTools.join(', ') : 'none specified'}`);
    lines.push(`MCP Servers: ${mcpServers.length > 0 ? mcpServers.join(', ') : 'none selected'}`);
    if (ideas.length > 0) {
      lines.push('## Core Requirements');
      ideas.forEach((idea) => {
        lines.push(`- ${idea.text}`);
      });
    }
    return lines.join('\n');
  }, [agentFramework, agentProvider, agentRole, agentTools, composedPurpose, ideas, mcpServers]);

  const handleGenerate = async () => {
    if (ideas.length === 0 || !agentProvider || !agentRole.trim()) {
      setGenerationError(t('requirementsMissing'));
      return;
    }
    setGenerationError(null);
    setIsGenerating(true);
    try {
      const result = await agentsService.generate(requirements);
      const entry: BlueprintEntry = {
        id: Date.now().toString(),
        createdAt: Date.now(),
        requirements,
        result,
      };
      setHistory((prev) => [entry, ...prev].slice(0, 10));
      setCurrentBlueprint(entry);
      setBlueprintModalOpen(true);
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);
      setGenerationError(message || t('generationFailed'));
    } finally {
      setIsGenerating(false);
    }
  };

  const handleOpenBlueprint = (entry: BlueprintEntry) => {
    setCurrentBlueprint(entry);
    setSaveStatus(null);
    setBlueprintModalOpen(true);
  };

  const handleCloseBlueprint = () => {
    setBlueprintModalOpen(false);
  };

  const handleCopyMarkdown = async () => {
    if (!currentBlueprint?.result.markdown) return;
    await navigator.clipboard.writeText(currentBlueprint.result.markdown);
    setCopied(true);
  };

  const handleCopyRequirements = async () => {
    if (!currentBlueprint?.requirements) return;
    await navigator.clipboard.writeText(currentBlueprint.requirements);
    setCopiedRequirements(true);
  };

  const handlePersistBlueprint = async (entry: BlueprintEntry | null = currentBlueprint) => {
    if (!entry || !entry.result.agents?.length) return;
    setIsSaving(true);
    setSaveStatus(null);
    try {
      for (const agent of entry.result.agents) {
        await agentsService.create(agent);
      }
      await fetchAgents();
      setSaveStatus({ variant: 'success', message: t('savedFeedback') });
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);
      setSaveStatus({ variant: 'error', message: message || t('saveFailed') });
    } finally {
      setIsSaving(false);
    }
  };

  const handleCopyAgentsTable = async () => {
    if (storedAgents.length === 0) return;
    const header = ['ID', 'Title', 'Role', 'Skills'];
    const rows = storedAgents.map((agent) => [
      agent.ID,
      agent.Title,
      agent.Role,
      agent.Skills.join(', '),
    ]);
    const csv = [header.join('\t'), ...rows.map((row) => row.join('\t'))].join('\n');
    await navigator.clipboard.writeText(csv);
    setCopiedAgentsTable(true);
  };

  const handleExportMarkdown = async () => {
    try {
      const markdown = await agentsService.exportMarkdown();
      const blob = new Blob([markdown], { type: 'text/markdown' });
      const url = URL.createObjectURL(blob);
      const anchor = document.createElement('a');
      anchor.href = url;
      anchor.download = `AGENTS-${Date.now()}.md`;
      document.body.appendChild(anchor);
      anchor.click();
      document.body.removeChild(anchor);
      URL.revokeObjectURL(url);
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);
      setAgentsError(message);
    }
  };

  const handleDeleteAgent = async (id: number) => {
    try {
      await agentsService.remove(id);
      setStoredAgents((prev) => prev.filter((agent) => agent.ID !== id));
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);
      setAgentsError(message);
    }
  };

  const selectedProviderInfo = useMemo(
    () => providers.find((provider) => provider.name === agentProvider),
    [agentProvider, providers],
  );

  const providerStatusLabel = useMemo(() => {
    if (!selectedProviderInfo) return '';
    switch (selectedProviderInfo.status) {
      case 'ready':
        return t('providerReady');
      case 'needs_api_key':
        return t('providerNeedsKey');
      default:
        return t('providerUnavailable');
    }
  }, [selectedProviderInfo, t]);

  const providerStatusTone = selectedProviderInfo?.status === 'ready'
    ? 'text-emerald-600 dark:text-emerald-300 bg-emerald-50 dark:bg-emerald-500/10'
    : selectedProviderInfo?.status === 'needs_api_key'
      ? 'text-amber-600 dark:text-amber-300 bg-amber-50 dark:bg-amber-500/10'
      : 'text-rose-600 dark:text-rose-300 bg-rose-50 dark:bg-rose-500/10';

  const formatTimestamp = useCallback((timestamp: number) => {
    try {
      return new Intl.DateTimeFormat(language, {
        dateStyle: 'short',
        timeStyle: 'short',
      }).format(new Date(timestamp));
    } catch {
      return new Date(timestamp).toLocaleString();
    }
  }, [language]);

  return (
    <div className="space-y-6">
      <Card
        title={t('overviewTitle')}
        description={t('overviewDescription')}
        action={
          <button
            type="button"
            onClick={handleExportMarkdown}
            className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white px-4 py-2 text-sm font-semibold text-slate-700 transition hover:border-[#38cde4] hover:text-[#0f172a] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#cbd5f5] dark:hover:border-[#38cde4]"
          >
            <Download className="h-4 w-4" />
            {t('exportMarkdown')}
          </button>
        }
      >
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div className="rounded-xl border border-slate-200 bg-slate-50/80 p-4 text-sm text-slate-600 dark:border-[#13263a] dark:bg-[#0f172a]/60 dark:text-[#94a3b8]">
            <p className="text-xs uppercase tracking-[0.3em] text-[#475569]/70 dark:text-[#64748b]">
              {t('frameworksLabel')}
            </p>
            <p className="mt-2 text-lg font-semibold text-[#0f172a] dark:text-[#e5f2f2]">{frameworks.find((fw) => fw.value === agentFramework)?.label}</p>
          </div>
          <div className="rounded-xl border border-slate-200 bg-slate-50/80 p-4 text-sm text-slate-600 dark:border-[#13263a] dark:bg-[#0f172a]/60 dark:text-[#94a3b8]">
            <p className="text-xs uppercase tracking-[0.3em] text-[#475569]/70 dark:text-[#64748b]">
              {t('providersLabel')}
            </p>
            <div className="mt-2 flex items-center gap-2">
              <span className="text-lg font-semibold text-[#0f172a] dark:text-[#e5f2f2]">{(selectedProviderInfo?.display_name ?? agentProvider) || '--'}</span>
              {selectedProviderInfo && (
                <span className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[11px] font-semibold ${providerStatusTone}`}>
                  <BadgeCheck className="h-3 w-3" />
                  {providerStatusLabel}
                </span>
              )}
            </div>
          </div>
          <div className="rounded-xl border border-slate-200 bg-slate-50/80 p-4 text-sm text-slate-600 dark:border-[#13263a] dark:bg-[#0f172a]/60 dark:text-[#94a3b8]">
            <p className="text-xs uppercase tracking-[0.3em] text-[#475569]/70 dark:text-[#64748b]">
              {t('toolsLabel')}
            </p>
            <p className="mt-2 text-lg font-semibold text-[#0f172a] dark:text-[#e5f2f2]">{agentTools.length}</p>
          </div>
          <div className="rounded-xl border border-slate-200 bg-slate-50/80 p-4 text-sm text-slate-600 dark:border-[#13263a] dark:bg-[#0f172a]/60 dark:text-[#94a3b8]">
            <p className="text-xs uppercase tracking-[0.3em] text-[#475569]/70 dark:text-[#64748b]">
              {t('mcpLabel')}
            </p>
            <p className="mt-2 text-lg font-semibold text-[#0f172a] dark:text-[#e5f2f2]">{mcpServers.length}</p>
          </div>
        </div>
      </Card>

      <Card
        title={t('workspaceTitle')}
        description={t('workspaceDescription')}
        action={
          <div className="flex flex-wrap items-center gap-2">
            <button
              type="button"
              onClick={handleGenerate}
              disabled={isGenerating || ideas.length === 0}
              className="inline-flex items-center gap-2 rounded-full border border-[#06b6d4] bg-[#06b6d4] px-4 py-2 text-sm font-semibold text-white shadow-soft-card transition hover:bg-[#0891b2] disabled:cursor-not-allowed disabled:opacity-60 focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/30 dark:border-[#06b6d4] dark:bg-[#06b6d4] dark:text-[#0a1523]"
            >
              {isGenerating ? <Loader2 className="h-4 w-4 animate-spin" /> : <Sparkles className="h-4 w-4" />}
              {isGenerating ? t('generatingButton') : t('generateButton')}
            </button>
            <button
              type="button"
              onClick={() => currentBlueprint ? setBlueprintModalOpen(true) : history[0] ? handleOpenBlueprint(history[0]) : null}
              disabled={!currentBlueprint && history.length === 0}
              className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white px-3 py-1.5 text-xs font-semibold text-slate-600 transition hover:border-[#38cde4] hover:text-[#0f172a] disabled:cursor-not-allowed disabled:opacity-50 dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8] dark:hover:border-[#38cde4]"
            >
              <History className="h-3.5 w-3.5" />
              {t('openBlueprint')}
            </button>
          </div>
        }
      >
        <div className="grid gap-6 lg:grid-cols-2">
          <section>
            <h3 className="text-base font-semibold text-[#0f172a] dark:text-[#e5f2f2]">{t('contextCardTitle')}</h3>
            <p className="mt-1 text-sm text-[#64748b] dark:text-[#94a3b8]">
              {t('ideasEmpty')}
            </p>
            <div className="mt-4 space-y-4">
              <textarea
                value={currentIdea}
                onChange={(event) => setCurrentIdea(event.target.value)}
                placeholder={t('ideaInputPlaceholder')}
                rows={6}
                className="w-full resize-none rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-[#475569] shadow-inner transition focus:border-[#06b6d4] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/20 dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#e5f2f2]"
              />
              <button
                type="button"
                onClick={addIdea}
                disabled={!currentIdea.trim()}
                className="inline-flex items-center gap-2 rounded-full border border-[#06b6d4] bg-[#06b6d4] px-4 py-2 text-sm font-semibold text-white shadow-soft-card transition hover:bg-[#0891b2] disabled:cursor-not-allowed disabled:opacity-60 focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/30 dark:border-[#06b6d4] dark:bg-[#06b6d4] dark:text-[#0a1523]"
              >
                <Plus className="h-4 w-4" />
                {t('addIdea')}
              </button>
              {ideas.length > 0 && (
                <ul className="space-y-3">
                  {ideas.map((idea) => (
                    <li
                      key={idea.id}
                      className="flex items-start justify-between rounded-2xl border border-slate-200/90 bg-white/85 px-4 py-3 text-sm text-[#475569] shadow-soft-card dark:border-[#13263a]/70 dark:bg-[#0a1523]/70 dark:text-[#e5f2f2]"
                    >
                      <span className="pr-3">{idea.text}</span>
                      <button
                        type="button"
                        onClick={() => removeIdea(idea.id)}
                        className="ml-3 inline-flex h-7 w-7 items-center justify-center rounded-full border border-slate-200 text-slate-500 transition hover:border-rose-400 hover:text-rose-500 dark:border-[#13263a] dark:text-[#64748b] dark:hover:border-rose-500 dark:hover:text-rose-300"
                        aria-label="Remove idea"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </section>

          <section className="space-y-5">
            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
                {t('roleLabel')}
              </label>
              <input
                value={agentRole}
                onChange={(event) => setAgentRole(event.target.value)}
                placeholder={t('rolePlaceholder')}
                className="w-full rounded-2xl border border-slate-200 bg-white px-4 py-2 text-sm text-[#475569] shadow-inner transition focus:border-[#06b6d4] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/20 dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#e5f2f2]"
              />
            </div>

            <div>
              <p className="mb-2 text-xs font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
                {t('purposeLabel')}
              </p>
              <div className="flex flex-wrap gap-2">
                {purposes.map((item) => (
                  <button
                    type="button"
                    key={item}
                    onClick={() => setPurpose(item)}
                    className={`rounded-full border px-3 py-1.5 text-xs font-semibold transition ${
                      purpose === item
                        ? 'border-[#06b6d4] bg-[#06b6d4] text-white dark:border-[#38cde4] dark:bg-[#38cde4] dark:text-[#0a1523]'
                        : 'border-slate-200 bg-white text-[#475569] hover:border-[#bae6fd] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8]'
                    }`}
                  >
                    {t(`purpose_${item}`)}
                  </button>
                ))}
              </div>
              {purpose === 'other' && (
                <input
                  value={customPurpose}
                  onChange={(event) => setCustomPurpose(event.target.value)}
                  placeholder={t('customPurposePlaceholder')}
                  className="mt-3 w-full rounded-2xl border border-slate-200 bg-white px-4 py-2 text-sm text-[#475569] shadow-inner transition focus:border-[#06b6d4] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/20 dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#e5f2f2]"
                />
              )}
            </div>

            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
                {t('frameworksLabel')}
              </label>
              <div className="space-y-2">
                {frameworks.map((framework) => (
                  <button
                    type="button"
                    key={framework.value}
                    onClick={() => setAgentFramework(framework.value)}
                    className={`w-full rounded-2xl border px-4 py-3 text-left text-sm transition ${
                      agentFramework === framework.value
                        ? 'border-[#06b6d4] bg-[#06b6d4]/10 text-[#0f172a] shadow-soft-card dark:border-[#38cde4] dark:bg-[#38cde4]/20 dark:text-[#e5f2f2]'
                        : 'border-slate-200 bg-white text-[#475569] hover:border-[#bae6fd] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8]'
                    }`}
                  >
                    <span className="block font-semibold">{framework.label}</span>
                    <span className="text-xs text-[#64748b] dark:text-[#94a3b8]">{framework.description}</span>
                  </button>
                ))}
              </div>
            </div>

            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
                {t('providersLabel')}
              </label>
              <select
                value={agentProvider}
                onChange={(event) => setAgentProvider(event.target.value)}
                className="w-full rounded-2xl border border-slate-200 bg-white px-4 py-2 text-sm text-[#475569] shadow-inner transition focus:border-[#06b6d4] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/20 dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#e5f2f2]"
              >
                <option value="" disabled>
                  {providersLoading ? 'Loading...' : providersError ? providersError : '--'}
                </option>
                {providers.map((provider) => (
                  <option key={provider.name} value={provider.name}>
                    {provider.display_name || provider.name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <p className="mb-2 text-xs font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
                {t('toolsSectionTitle')}
              </p>
              <div className="grid grid-cols-2 gap-2">
                {toolCatalog.map((tool) => (
                  <button
                    type="button"
                    key={tool}
                    onClick={() => toggleTool(tool)}
                    className={`rounded-xl border px-3 py-2 text-xs font-semibold transition ${
                      agentTools.includes(tool)
                        ? 'border-emerald-400 bg-emerald-50 text-emerald-700 dark:border-emerald-400/70 dark:bg-emerald-400/10 dark:text-emerald-200'
                        : 'border-slate-200 bg-white text-[#475569] hover:border-[#bae6fd] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8]'
                    }`}
                  >
                    {tool}
                  </button>
                ))}
              </div>
            </div>

            <div>
              <p className="mb-2 text-xs font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
                {t('mcpSectionTitle')}
              </p>
              <div className="grid grid-cols-2 gap-2">
                {mcpCatalog.map((server) => (
                  <button
                    type="button"
                    key={server.name}
                    onClick={() => toggleMcp(server.name)}
                    className={`rounded-xl border px-3 py-2 text-xs font-semibold transition ${
                      mcpServers.includes(server.name)
                        ? 'border-[#7c4dff] bg-[#7c4dff]/10 text-[#4c1d95] dark:border-[#7c4dff]/70 dark:bg-[#7c4dff]/20 dark:text-[#e0d7ff]'
                        : 'border-slate-200 bg-white text-[#475569] hover:border-[#bae6fd] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8]'
                    }`}
                    title={server.desc}
                  >
                    {server.desc}
                  </button>
                ))}
              </div>
              <div className="mt-3 flex items-center gap-2">
                <input
                  value={customMcp}
                  onChange={(event) => setCustomMcp(event.target.value)}
                  placeholder={t('mcpCustomPlaceholder')}
                  className="flex-1 rounded-2xl border border-slate-200 bg-white px-4 py-2 text-sm text-[#475569] shadow-inner transition focus:border-[#06b6d4] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/20 dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#e5f2f2]"
                />
                <button
                  type="button"
                  onClick={() => {
                    if (!customMcp.trim()) return;
                    toggleMcp(customMcp.trim());
                    setCustomMcp('');
                  }}
                  className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white px-3 py-1.5 text-xs font-semibold text-slate-600 transition hover:border-[#38cde4] hover:text-[#0f172a] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8] dark:hover:border-[#38cde4]"
                >
                  <Plus className="h-3.5 w-3.5" />
                  MCP
                </button>
              </div>
            </div>
          </section>
        </div>

        {generationError && (
          <p className="mt-6 rounded-xl border border-rose-400/60 bg-rose-50/80 p-3 text-sm text-rose-600 dark:border-rose-500/50 dark:bg-rose-500/10 dark:text-rose-200">
            {generationError}
          </p>
        )}

        <details className="mt-6 overflow-hidden rounded-2xl border border-dashed border-slate-200/80 bg-white/90 text-sm text-[#475569] shadow-inner dark:border-[#13263a]/80 dark:bg-[#0a1523]/70 dark:text-[#cbd5f5]">
          <summary className="cursor-pointer px-4 py-3 text-xs font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
            {t('requirementsPreview')}
          </summary>
          <pre className="max-h-48 overflow-y-auto px-4 pb-4 pt-2 text-xs leading-relaxed">{requirements}</pre>
        </details>
      </Card>

      <Card
        title={t('blueprintHistoryTitle')}
        description={t('workspaceDescription')}
      >
        {history.length === 0 ? (
          <p className="rounded-xl border border-dashed border-slate-200 bg-slate-50/80 p-4 text-sm text-slate-500 dark:border-[#13263a] dark:bg-[#0f172a]/60 dark:text-[#94a3b8]">
            {t('blueprintHistoryEmpty')}
          </p>
        ) : (
          <div className="space-y-3">
            {history.map((entry) => (
              <div
                key={entry.id}
                className="flex flex-col gap-3 rounded-2xl border border-slate-200/80 bg-white/85 px-4 py-3 text-sm text-[#475569] shadow-soft-card dark:border-[#13263a]/70 dark:bg-[#0a1523]/70 dark:text-[#e5f2f2] md:flex-row md:items-center md:justify-between"
              >
                <div>
                  <p className="font-semibold text-[#0f172a] dark:text-[#e5f2f2]">
                    {t('generatedAt')}: {formatTimestamp(entry.createdAt)}
                  </p>
                  <p className="text-xs text-[#64748b] dark:text-[#94a3b8]">
                    {entry.result.agents.length} {t('generatedAgents').toLowerCase()}
                  </p>
                </div>
                <div className="flex flex-wrap items-center gap-2">
                  <button
                    type="button"
                    onClick={() => handleOpenBlueprint(entry)}
                    className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white px-3 py-1.5 text-xs font-semibold text-slate-600 transition hover:border-[#38cde4] hover:text-[#0f172a] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8] dark:hover:border-[#38cde4]"
                  >
                    <History className="h-3.5 w-3.5" />
                    {t('viewBlueprint')}
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>

      <Card
        title={t('storedAgentsTitle')}
        description={t('storedAgentsDescription')}
        action={
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={fetchAgents}
              className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white px-3 py-1.5 text-xs font-semibold text-slate-600 transition hover:border-[#38cde4] hover:text-[#0f172a] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8] dark:hover:border-[#38cde4]"
            >
              <RefreshCcw className="h-3.5 w-3.5" />
              {t('refreshAgents')}
            </button>
            <button
              type="button"
              onClick={handleCopyAgentsTable}
              disabled={storedAgents.length === 0}
              className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white px-3 py-1.5 text-xs font-semibold text-slate-600 transition hover:border-[#38cde4] hover:text-[#0f172a] disabled:cursor-not-allowed disabled:opacity-50 dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8] dark:hover:border-[#38cde4]"
            >
              {copiedAgentsTable ? <ClipboardCheck className="h-3.5 w-3.5" /> : <ClipboardCopy className="h-3.5 w-3.5" />}
              {copiedAgentsTable ? t('copied') : t('copyAgentsTable')}
            </button>
          </div>
        }
      >
        {agentsError && (
          <p className="mb-3 rounded-xl border border-rose-400/60 bg-rose-50/80 p-3 text-sm text-rose-600 dark:border-rose-500/50 dark:bg-rose-500/10 dark:text-rose-200">
            {agentsError}
          </p>
        )}
        {agentsLoading ? (
          <div className="flex items-center gap-2 text-sm text-[#475569] dark:text-[#94a3b8]">
            <Loader2 className="h-4 w-4 animate-spin" />
            {t('generatingButton')}
          </div>
        ) : storedAgents.length === 0 ? (
          <p className="rounded-xl border border-dashed border-slate-200 bg-slate-50/80 p-4 text-sm text-slate-500 dark:border-[#13263a] dark:bg-[#0f172a]/60 dark:text-[#94a3b8]">
            {t('emptyStoredAgents')}
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-slate-200 text-sm dark:divide-[#13263a]">
              <thead className="bg-slate-100/80 dark:bg-[#101926]">
                <tr>
                  <th className="px-3 py-2 text-left font-semibold text-[#475569] dark:text-[#cbd5f5]">{t('idColumn')}</th>
                  <th className="px-3 py-2 text-left font-semibold text-[#475569] dark:text-[#cbd5f5]">{t('titleColumn')}</th>
                  <th className="px-3 py-2 text-left font-semibold text-[#475569] dark:text-[#cbd5f5]">{t('roleColumn')}</th>
                  <th className="px-3 py-2 text-left font-semibold text-[#475569] dark:text-[#cbd5f5]">{t('skillsColumn')}</th>
                  <th className="px-3 py-2 text-left font-semibold text-[#475569] dark:text-[#cbd5f5]">{t('actionsColumn')}</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-200 dark:divide-[#13263a]">
                {storedAgents.map((agent) => (
                  <tr key={agent.ID} className="bg-white/90 dark:bg-[#0a1523]/70">
                    <td className="px-3 py-2 text-[#475569] dark:text-[#cbd5f5]">{agent.ID}</td>
                    <td className="px-3 py-2 text-[#0f172a] dark:text-[#e5f2f2]">{agent.Title}</td>
                    <td className="px-3 py-2 text-[#475569] dark:text-[#cbd5f5]">{agent.Role}</td>
                    <td className="px-3 py-2 text-[#475569] dark:text-[#cbd5f5]">{agent.Skills.join(', ')}</td>
                    <td className="px-3 py-2">
                      <button
                        type="button"
                        onClick={() => handleDeleteAgent(agent.ID)}
                        className="inline-flex items-center gap-1 rounded-full border border-rose-400 bg-rose-50 px-3 py-1 text-xs font-semibold text-rose-600 transition hover:bg-rose-100 dark:border-rose-400/70 dark:bg-rose-500/10 dark:text-rose-200"
                      >
                        <Trash2 className="h-3.5 w-3.5" />
                        {t('deleteAction')}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>

      {isBlueprintModalOpen && currentBlueprint && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-[#111827]/70 backdrop-blur-sm px-4 py-8">
          <div className="relative w-full max-w-4xl max-h-[90vh] overflow-y-auto rounded-3xl border border-slate-200 bg-white/95 shadow-2xl dark:border-[#13263a] dark:bg-[#0a1523]/95">
            <header className="sticky top-0 flex items-center justify-between border-b border-slate-200 bg-white/95 px-6 py-4 dark:border-[#13263a] dark:bg-[#0a1523]/95">
              <div>
                <h2 className="text-lg font-semibold text-[#0f172a] dark:text-[#e5f2f2]">{t('generationTitle')}</h2>
                <p className="text-xs text-[#64748b] dark:text-[#94a3b8]">
                  {t('generatedAt')}: {formatTimestamp(currentBlueprint.createdAt)}
                </p>
              </div>
              <button
                type="button"
                onClick={handleCloseBlueprint}
                className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-slate-200 text-slate-500 transition hover:border-[#bae6fd] hover:text-[#0f172a] dark:border-[#13263a] dark:text-[#94a3b8]"
                aria-label={t('close')}
              >
                <X className="h-4 w-4" />
              </button>
            </header>

            <div className="space-y-6 px-6 py-6">
              <section>
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
                    {t('requirementsPreview')}
                  </h3>
                  <button
                    type="button"
                    onClick={handleCopyRequirements}
                    className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white px-3 py-1.5 text-xs font-semibold text-slate-600 transition hover:border-[#38cde4] hover:text-[#0f172a] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8] dark:hover:border-[#38cde4]"
                  >
                    {copiedRequirements ? <ClipboardCheck className="h-3.5 w-3.5" /> : <ClipboardCopy className="h-3.5 w-3.5" />}
                    {copiedRequirements ? t('copied') : t('copyRequirements')}
                  </button>
                </div>
                <pre className="mt-3 max-h-64 overflow-y-auto rounded-2xl border border-dashed border-slate-200/80 bg-white/90 p-4 text-xs text-[#475569] shadow-inner dark:border-[#13263a]/80 dark:bg-[#0a1523]/70 dark:text-[#cbd5f5]">
                  {currentBlueprint.requirements}
                </pre>
              </section>

              <section>
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
                    {t('generatedAgents')}
                  </h3>
                  <div className="flex flex-wrap items-center gap-2">
                    <button
                      type="button"
                      onClick={handleCopyMarkdown}
                      disabled={!currentBlueprint.result.markdown}
                      className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white px-3 py-1.5 text-xs font-semibold text-slate-600 transition hover:border-[#38cde4] hover:text-[#0f172a] disabled:cursor-not-allowed disabled:opacity-50 dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8] dark:hover:border-[#38cde4]"
                    >
                      {copied ? <ClipboardCheck className="h-3.5 w-3.5" /> : <ClipboardCopy className="h-3.5 w-3.5" />}
                      {copied ? t('copied') : t('copyMarkdown')}
                    </button>
                    <button
                      type="button"
                      onClick={() => handlePersistBlueprint(currentBlueprint)}
                      disabled={isSaving || !currentBlueprint.result.agents.length}
                      className="inline-flex items-center gap-2 rounded-full border border-emerald-500 bg-emerald-500 px-3 py-1.5 text-xs font-semibold text-white transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-60 dark:border-emerald-400 dark:bg-emerald-400 dark:text-[#0a1523]"
                    >
                      {isSaving ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <Save className="h-3.5 w-3.5" />}
                      {isSaving ? t('saving') : t('saveAgents')}
                    </button>
                  </div>
                </div>

                {saveStatus && (
                  <p className={`mt-3 text-xs ${saveStatus.variant === 'success' ? 'text-emerald-600 dark:text-emerald-300' : 'text-rose-600 dark:text-rose-300'}`}>
                    {saveStatus.message}
                  </p>
                )}

                <div className="mt-3 space-y-3">
                  {currentBlueprint.result.agents.map((agent, index) => (
                    <div
                      key={`${agent.Title}-${index}`}
                      className="rounded-2xl border border-slate-200/80 bg-white/85 p-4 shadow-soft-card dark:border-[#13263a]/70 dark:bg-[#0a1523]/70"
                    >
                      <p className="text-sm font-semibold text-[#0f172a] dark:text-[#e5f2f2]">{agent.Title}</p>
                      <p className="mt-1 text-xs text-[#64748b] dark:text-[#94a3b8]">{agent.Role}</p>
                      {agent.Skills.length > 0 && (
                        <p className="mt-2 text-xs text-[#475569] dark:text-[#cbd5f5]">
                          <span className="font-semibold">Skills:</span> {agent.Skills.join(', ')}
                        </p>
                      )}
                      {agent.Restrictions.length > 0 && (
                        <p className="mt-1 text-xs text-[#475569] dark:text-[#cbd5f5]">
                          <span className="font-semibold">Restrictions:</span> {agent.Restrictions.join(', ')}
                        </p>
                      )}
                      {agent.PromptExample && (
                        <p className="mt-2 rounded-xl border border-slate-200 bg-slate-50/80 p-3 text-xs text-[#475569] dark:border-[#13263a] dark:bg-[#101926] dark:text-[#cbd5f5]">
                          {agent.PromptExample}
                        </p>
                      )}
                    </div>
                  ))}
                </div>
              </section>

              {currentBlueprint.result.markdown && currentBlueprint.result.agents.length > 0 && (
                <section>
                  <h3 className="text-sm font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
                    {t('markdownSection')}
                  </h3>
                  <div className="mt-3 max-h-72 overflow-y-auto rounded-2xl border border-slate-200/80 bg-white/90 p-4 text-xs text-[#475569] shadow-inner dark:border-[#13263a]/80 dark:bg-[#0a1523]/70 dark:text-[#cbd5f5]">
                    <ReactMarkdown>{currentBlueprint.result.markdown}</ReactMarkdown>
                  </div>
                </section>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AgentsGenerator;
