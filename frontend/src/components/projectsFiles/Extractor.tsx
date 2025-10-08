import { AnimatePresence, motion } from 'framer-motion';
import {
  Download as ArrowDownTrayIcon,
  RotateCw as ArrowPathIcon,
  BarChart2 as ChartBarIcon,
  Clipboard as ClipboardDocumentIcon,
  Code2 as CodeBracketIcon,
  FileText as DocumentIcon,
  Eye as EyeIcon,
  Folder as FolderIcon,
  Funnel as FunnelIcon,
  Search as MagnifyingGlassIcon,
  Play as PlayIcon
} from "lucide-react";
import * as React from 'react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslations } from '../../i18n/useTranslations';
import { pack } from '../../utils/lookatni';

/* ========= Types (mantidos / estendidos) ========= */

interface ProjectFile {
  path: string;
  content: string;
  size: number;
  lines: number;
}

interface ProjectStats {
  totalFiles: number;
  totalMarkers: number;
  totalBytes: number;
  errors: Array<{ line: number; message: string }>;
}

interface ProjectExtractorProps {
  projectFile: string;
  projectName: string;
  description?: string;
}

interface LookAtniExtractedFile {
  id: string;
  name: string;
  path: string;
  content: string;
  language?: string;
  size?: number;
  line_count?: number;
  fragments?: Array<Record<string, unknown>>;
  metadata?: Record<string, string>;
}

interface LookAtniMetadata {
  languages?: Record<string, number>;
  total_lines?: number;
  total_files?: number;
  total_fragments?: number;
  extraction_time?: number;
}

interface LookAtniExtractedProject {
  project_name: string;
  structure?: Record<string, unknown>;
  files: LookAtniExtractedFile[];
  fragments?: Array<Record<string, unknown>>;
  metadata?: LookAtniMetadata;
  download_url?: string;
  extracted_at?: string;
}

interface LookAtniArchiveResponse {
  success?: boolean;
  archive_path?: string;
  download_url?: string;
  file_name?: string;
}

/* ========= Utils locais (sem libs) ========= */

function classNames(...xs: Array<string | false | null | undefined>) {
  return xs.filter(Boolean).join(' ');
}

function formatFileSize(bytes: number) {
  if (bytes < 1024) return `${bytes}B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)}KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)}MB`;
}

function getFileIcon(filename: string) {
  const ext = filename.split('.').pop()?.toLowerCase();
  const isFolder = filename.includes('/') && !filename.split('/').pop()?.includes('.');
  if (isFolder) return <FolderIcon className="w-4 h-4 text-[#06b6d4]" />;
  if (['js', 'ts', 'jsx', 'tsx', 'go', 'rs', 'py', 'sh', 'json', 'css', 'html', 'md'].includes(ext || '')) {
    return <CodeBracketIcon className="w-4 h-4 text-[#f59e0b]" />;
  }
  return <DocumentIcon className="w-4 h-4 text-[#94a3b8]" />;
}

function extOf(path: string) {
  return path.includes('.') ? path.split('.').pop()!.toLowerCase() : '';
}

/** highlight opcional (se existir window.hljs) */
function maybeHighlight(code: string, langHint?: string) {
  const anyWin = window as any;
  const hljs = anyWin?.hljs;
  if (!hljs) return code;
  try {
    if (langHint && hljs.getLanguage(langHint)) {
      return hljs.highlight(code, { language: langHint }).value;
    }
    return hljs.highlightAuto(code).value;
  } catch {
    return code;
  }
}

function normalizeExtractedProject(project: LookAtniExtractedProject | null) {
  if (!project) {
    return null;
  }

  const files: ProjectFile[] = project.files?.map((file) => {
    const size = typeof file.size === 'number' ? file.size : file.content.length;
    const lineCount = typeof file.line_count === 'number'
      ? file.line_count
      : file.content.split(/\r?\n/).length;

    return {
      path: file.path || file.name,
      content: file.content,
      size,
      lines: lineCount,
    };
  }) ?? [];

  const totalBytes = files.reduce((acc, file) => acc + file.size, 0);
  const totalMarkers = project.metadata?.total_fragments
    ?? project.fragments?.length
    ?? 0;

  const stats: ProjectStats = {
    totalFiles: project.metadata?.total_files ?? files.length,
    totalMarkers,
    totalBytes,
    errors: [],
  };

  return { stats, files };
}

/* ========= File tree básico ========= */

type TreeNode = {
  name: string;
  path: string;
  children?: TreeNode[];
  file?: ProjectFile;
  open?: boolean;
};

function buildTree(files: ProjectFile[]): TreeNode {
  const root: TreeNode = { name: '', path: '', children: [], open: true };
  for (const f of files) {
    const parts = f.path.split('/');
    let cur = root;
    parts.forEach((part, idx) => {
      if (idx === parts.length - 1) {
        // file
        cur.children = cur.children || [];
        cur.children.push({ name: part, path: f.path, file: f });
      } else {
        const existing = cur.children?.find((n) => n.name === part && !n.file);
        if (existing) {
          cur = existing;
        } else {
          const n: TreeNode = { name: part, path: parts.slice(0, idx + 1).join('/'), children: [], open: idx < 2 };
          cur.children!.push(n);
          cur = n;
        }
      }
    });
  }
  // ordenar: pastas primeiro, depois arquivos; ambos por nome
  const sortRec = (node: TreeNode) => {
    if (!node.children) return;
    node.children.sort((a, b) => {
      const af = !!a.file;
      const bf = !!b.file;
      if (af !== bf) return af ? 1 : -1;
      return a.name.localeCompare(b.name);
    });
    node.children.forEach(sortRec);
  };
  sortRec(root);
  return root;
}

/* ========= Componente ========= */

export default function ProjectExtractor({ projectFile, projectName, description }: ProjectExtractorProps) {
  /* ========== Traduções ========= */
  const { t } = useTranslations();
  /* ========= Estados ========= */
  const [isLoading, setIsLoading] = useState(false);
  const [extractedProject, setExtractedProject] = useState<LookAtniExtractedProject | null>(null);
  const [projectData, setProjectData] = useState<{ stats: ProjectStats; files: ProjectFile[] } | null>(null);
  const [selectedFile, setSelectedFile] = useState<ProjectFile | null>(null);
  const [showStats, setShowStats] = useState(false);
  const [extractionMode, setExtractionMode] = useState<'preview' | 'download'>('preview');

  // novos estados
  const [error, setError] = useState<string | null>(null);
  const [query, setQuery] = useState('');
  const [extFilter, setExtFilter] = useState<string>('all');
  const [wrap, setWrap] = useState(true);
  const [fontSize, setFontSize] = useState<number>(13);

  const abortRef = useRef<AbortController | null>(null);
  const searchRef = useRef<HTMLInputElement | null>(null);

  // micro-telemetria
  const t0 = useRef<number>(0);

  const extractProject = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    setExtractedProject(null);
    setProjectData(null);
    setSelectedFile(null);
    abortRef.current?.abort();
    const ctrl = new AbortController();
    abortRef.current = ctrl;
    t0.current = performance.now();

    try {
      const localPath = projectFile.startsWith('.') || projectFile.startsWith('/')
        ? projectFile
        : `./${projectFile}`;

      const response = await fetch('/api/v1/lookatni/extract', {
        method: 'POST',
        signal: ctrl.signal,
        headers: {
          'Content-Type': 'application/json',
          'x-kubex-client': 'grompt/pe',
        },
        body: JSON.stringify({
          local_path: localPath,
          fragment_by: 'file',
          context_depth: 2,
          include_hidden: false,
        }),
      });

      if (!response.ok) {
        const message = await response.text();
        throw new Error(message || 'Falha na extração');
      }

      const raw: LookAtniExtractedProject = await response.json();
      setExtractedProject(raw);

      const normalized = normalizeExtractedProject(raw);
      if (!normalized || normalized.files.length === 0) {
        setProjectData(null);
        setSelectedFile(null);
        setError('Nenhum arquivo extraído do projeto.');
        return;
      }

      setProjectData(normalized);
      setSelectedFile(normalized.files[0]);
    } catch (err: any) {
      if (err.name !== 'AbortError') {
        setError(err?.message ?? 'Erro ao extrair projeto.');
      }
    } finally {
      setIsLoading(false);
      const dt = (performance.now() - t0.current).toFixed(0);
      // eslint-disable-next-line no-console
      console.log(`[PE] extract done in ${dt}ms`);
    }
  }, [projectFile]);

  const downloadProject = useCallback(async () => {
    if (!extractedProject) {
      setError('Extraia o projeto antes de gerar o pacote.');
      return;
    }

    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch('/api/v1/lookatni/archive', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'x-kubex-client': 'grompt/pe',
        },
        body: JSON.stringify(extractedProject),
      });

      if (!response.ok) {
        throw new Error(`Falha no download (${response.status})`);
      }

      const result: LookAtniArchiveResponse = await response.json();
      const downloadUrl = result.download_url || result.archive_path;

      if (!downloadUrl) {
        throw new Error('Download indisponível.');
      }

      const anchor = document.createElement('a');
      anchor.href = downloadUrl;
      anchor.download = result.file_name || `${projectName}.zip`;
      document.body.appendChild(anchor);
      anchor.click();
      document.body.removeChild(anchor);
    } catch (err: any) {
      setError(err?.message ?? 'Erro ao gerar pacote navegável.');
    } finally {
      setIsLoading(false);
    }
  }, [extractedProject, projectName]);

  const downloadSourceFile = useCallback(() => {
    if (!projectData?.files?.length) {
      setError('Extraia o projeto antes de exportar.');
      return;
    }

    setError(null);

    try {
      const lookatniPackage = pack(
        projectData.files.map((file) => ({ path: file.path, content: file.content }))
      );

      const blob = new Blob([lookatniPackage], { type: 'text/plain;charset=utf-8' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${projectName}.lkt.txt`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch {
      setError('Erro ao gerar pacote lookatni.');
    }
  }, [projectData, projectName]);

  // filtros/busca
  const filteredFiles = useMemo(() => {
    if (!projectData?.files) return [];
    const q = query.trim().toLowerCase();
    return projectData.files.filter((f) => {
      const okExt = extFilter === 'all' ? true : extOf(f.path) === extFilter;
      if (!okExt) return false;
      if (!q) return true;
      // fuzzy-lite: procura em path e primeiras linhas do conteúdo
      if (f.path.toLowerCase().includes(q)) return true;
      return f.content.slice(0, 1000).toLowerCase().includes(q);
    });
  }, [projectData, query, extFilter]);

  const fileTree = useMemo(() => buildTree(filteredFiles), [filteredFiles]);

  // atalhos de teclado
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === '/' && !e.metaKey && !e.ctrlKey) {
        e.preventDefault();
        searchRef.current?.focus();
      }
      if (!projectData?.files?.length) return;
      const idx = selectedFile ? filteredFiles.findIndex((x) => x.path === selectedFile.path) : -1;
      if ((e.key === 'j' || e.key === 'ArrowDown') && idx < filteredFiles.length - 1) {
        e.preventDefault();
        setSelectedFile(filteredFiles[idx + 1]);
      }
      if ((e.key === 'k' || e.key === 'ArrowUp') && idx > 0) {
        e.preventDefault();
        setSelectedFile(filteredFiles[idx - 1]);
      }
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
        e.preventDefault();
        downloadProject();
      }
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'p') {
        e.preventDefault();
        setExtractionMode((m) => (m === 'preview' ? 'download' : 'preview'));
      }
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [projectData, filteredFiles, selectedFile, downloadProject]);

  // tree toggle
  const toggleNode = (node: TreeNode) => {
    node.open = !node.open;
    // força re-render
    setShowStats((s) => s);
  };

  /* ======== Render ======== */

  return (
    <div className="rounded-2xl border border-[#e2e8f0] bg-white/95 shadow-soft-card overflow-hidden dark:border-[#13263a] dark:bg-[#0a1523]/90">
      {/* Header */}
      <div className="bg-gradient-to-r from-[#06b6d4] via-[#38bdf8] to-[#a855f7] p-4 text-white sticky top-0 z-10">
        <div className="flex items-center justify-between gap-3">
          <div className="min-w-0">
            <h3 className="text-lg font-bold truncate">{projectName}</h3>
            {description && <p className="text-sm opacity-90 truncate">{description}</p>}
          </div>
          <div className="flex gap-2">
            <button
              aria-pressed={showStats}
              onClick={() => setShowStats((s) => !s)}
              className="p-2 bg-white/20 rounded-lg hover:bg-white/30 transition-colors"
              title={t(showStats ? 'hideStats' : 'showStats')}
            >
              <ChartBarIcon className="w-5 h-5" />
            </button>
          </div>
        </div>
        {/* Search + Filters */}
        <div className="mt-3 flex items-center gap-2">
          <div className="relative flex-1">
            <MagnifyingGlassIcon className="w-4 h-4 absolute left-2 top-2.5 opacity-70" />
            <input
              ref={searchRef}
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder={t('searchFiles')}
              // placeholder="Buscar (atalho: /)"
              className="w-full pl-8 pr-3 py-2 rounded-md text-sm text-gray-900"
            />
          </div>
          <div className="flex items-center gap-2">
            <FunnelIcon className="w-4 h-4 opacity-80" />
            <select
              value={extFilter}
              onChange={(e) => setExtFilter(e.target.value)}
              className="bg-white/90 text-gray-900 rounded-md px-2 py-1 text-sm"
              title={t('filterByExtension')}
            // title="Filtrar por extensão"
            >
              <option value="all">Todas</option>
              <option value="ts">.ts</option>
              <option value="tsx">.tsx</option>
              <option value="js">.js</option>
              <option value="json">.json</option>
              <option value="md">.md</option>
              <option value="css">.css</option>
              <option value="go">.go</option>
            </select>
            <label className="ml-2 text-xs opacity-90">
              <input type="checkbox" checked={wrap} onChange={(e) => setWrap(e.target.checked)} className="mr-1" />
              wrap
            </label>
            <div className="flex items-center gap-1">
              <span className="text-xs">A</span>
              <input
                title={t('adjustFontSize')}
                type="range"
                min={11}
                max={18}
                value={fontSize}
                onChange={(e) => setFontSize(parseInt(e.target.value))}
              />
              <span className="text-sm">A</span>
            </div>
          </div>
        </div>
      </div>

      {/* Stats Panel */}
      <AnimatePresence>
        {showStats && projectData && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            className="bg-[#f9fafb] dark:bg-[#0a1e2b] p-4 border-b border-[#e2e8f0] dark:border-[#13263a]"
          >
            <div className="grid grid-cols-1 gap-4 text-center sm:grid-cols-3">
              <div className="rounded-xl border border-[#e2e8f0] bg-white/80 p-3 shadow-sm dark:border-[#13263a] dark:bg-[#0a1523]/80">
                <div className="text-2xl font-semibold text-[#06b6d4]">{projectData.stats.totalFiles}</div>
                <div className="text-sm text-[#64748b] dark:text-[#94a3b8]">{t('files')}</div>
              </div>
              <div className="rounded-xl border border-[#e2e8f0] bg-white/80 p-3 shadow-sm dark:border-[#13263a] dark:bg-[#0a1523]/80">
                <div className="text-2xl font-semibold text-[#0891b2]">{formatFileSize(projectData.stats.totalBytes)}</div>
                <div className="text-sm text-[#64748b] dark:text-[#94a3b8]">{t('totalSizeLabel')}</div>
              </div>
              <div className="rounded-xl border border-[#e2e8f0] bg-white/80 p-3 shadow-sm dark:border-[#13263a] dark:bg-[#0a1523]/80">
                <div className="text-2xl font-semibold text-[#a855f7]">{projectData.stats.totalMarkers}</div>
                <div className="text-sm text-[#64748b] dark:text-[#94a3b8]">{t('fragmentsLabel')}</div>
              </div>
            </div>
            {!!projectData.stats.errors?.length && (
              <div className="mt-3 text-left text-red-600 text-sm">
                {projectData.stats.errors.slice(0, 3).map((e, i) => (
                  <div key={i}>Linha {e.line}: {e.message}</div>
                ))}
              </div>
            )}
          </motion.div>
        )}
      </AnimatePresence>

      {/* Action Buttons */}
      <div className="p-4 border-b bg-[#f9fafb] dark:bg-[#0a1e2b] border-[#e2e8f0] dark:border-[#13263a]">
        <div className="flex gap-2 flex-wrap">
          {!projectData ? (
            <>
              <button
                onClick={extractProject}
                disabled={isLoading}
                className="flex items-center gap-2 rounded-full border border-transparent bg-[#06b6d4] px-5 py-2 text-sm font-semibold text-white shadow-soft-card transition hover:bg-[#0891b2] disabled:cursor-not-allowed disabled:opacity-60"
              >
                {isLoading ? (
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                ) : (
                  <PlayIcon className="w-4 h-4" />
                )}
                {isLoading ? <span>{t('extracting')}</span> : t('extractProject')}
              </button>
              <button
                onClick={downloadSourceFile}
                className="flex items-center gap-2 rounded-full border border-[#a855f7]/40 bg-white px-5 py-2 text-sm font-semibold text-[#6b21a8] transition hover:border-[#a855f7] hover:bg-[#f5f3ff] dark:border-[#5b21b6]/60 dark:bg-[#0a1523] dark:text-[#d8b4fe] dark:hover:bg-[#1b2534]"
                title={t('downloadOriginalFile')}
              >
                <DocumentIcon className="w-4 h-4" />
                {t('downloadOriginal')}
              </button>
            </>
          ) : (
            <>
              <button
                onClick={() => setExtractionMode('preview')}
                className={classNames(
                  'flex items-center gap-2 rounded-full px-5 py-2 text-sm font-semibold transition-colors',
                  extractionMode === 'preview'
                    ? 'bg-[#06b6d4] text-white shadow-soft-card'
                    : 'border border-[#e2e8f0] bg-white text-[#475569] hover:border-[#bae6fd] hover:bg-[#ecfeff] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8] dark:hover:bg-[#1b2534]'
                )}
              >
                <EyeIcon className="w-4 h-4" />
                {t('preview')}
              </button>
              <button
                onClick={downloadProject}
                disabled={isLoading}
                className="flex items-center gap-2 rounded-full border border-transparent bg-[#16a34a] px-5 py-2 text-sm font-semibold text-white shadow-soft-card transition hover:bg-[#15803d] disabled:cursor-not-allowed disabled:opacity-60"
                title="Ctrl/Cmd+S"
              >
                {isLoading ? (
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                ) : (
                  <ArrowDownTrayIcon className="w-4 h-4" />
                )}
                {t('downloadZip')}
              </button>
              <button
                onClick={downloadSourceFile}
                className="flex items-center gap-2 rounded-full border border-[#a855f7]/40 bg-white px-5 py-2 text-sm font-semibold text-[#6b21a8] transition hover:border-[#a855f7] hover:bg-[#f5f3ff] dark:border-[#5b21b6]/60 dark:bg-[#0a1523] dark:text-[#d8b4fe] dark:hover:bg-[#1b2534]"
              >
                <DocumentIcon className="w-4 h-4" />
                {t('downloadOriginal')}
              </button>
              <button
                onClick={() => extractProject()}
                className="flex items-center gap-2 rounded-full border border-[#e2e8f0] px-4 py-2 text-sm font-semibold text-[#475569] hover:border-[#06b6d4]/60 hover:text-[#0f172a] dark:border-[#13263a] dark:text-[#94a3b8] dark:hover:border-[#38cde4]"
                title={t('reprocess')}
              >
                <ArrowPathIcon className="w-4 h-4" />
                {t('reprocess')}
              </button>
            </>
          )}
        </div>
        {error && (
          <div className="mt-2 text-sm text-red-600">
            {error}{' '}
            <button onClick={() => extractProject()} className="underline">
              {t('tryAgain')}
            </button>
          </div>
        )}
      </div>

      {/* Master-Detail */}
      <AnimatePresence>
        {projectData && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: '500px' }}
            exit={{ opacity: 0, height: 0 }}
            className="grid grid-cols-12 h-[500px]"
            role="region"
            aria-label={t('fileExplorer')}
          >
            {/* Left: Tree */}
            <div className="col-span-4 border-r border-[#e2e8f0] dark:border-[#13263a] flex flex-col bg-[#f9fafb] dark:bg-[#0a1523] overflow-hidden">
              <div className="p-3 border-b border-[#e2e8f0] dark:border-[#13263a] bg-white/85 dark:bg-[#0a1e2b]">
                <h4 className="font-semibold text-sm text-[#111827] dark:text-[#e5f2f2] flex items-center gap-2">
                  <FolderIcon className="w-4 h-4 text-[#06b6d4]" />
                  {t('files')} ({filteredFiles.length})
                </h4>
              </div>
              <div className="flex-1 overflow-auto p-1">
                <TreeView
                  node={fileTree}
                  selectedPath={selectedFile?.path || ''}
                  onSelect={(f) => setSelectedFile(f)}
                  onToggle={(n) => toggleNode(n)}
                />
              </div>
            </div>

            {/* Right: Content */}
            <div className="col-span-8 flex flex-col bg-white/95 dark:bg-[#050a12] overflow-hidden">
              {selectedFile ? (
                <>
                  {/* Header */}
                  <div className="p-3 border-b border-[#e2e8f0] bg-[#f9fafb] dark:border-[#13263a] dark:bg-[#0a1e2b]">
                    <div className="flex items-center justify-between gap-2">
                      <div className="flex items-center gap-2 min-w-0 flex-1">
                        {getFileIcon(selectedFile.path)}
                        <span className="font-mono text-sm text-[#111827] dark:text-[#e5f2f2] truncate">
                          {selectedFile.path}
                        </span>
                      </div>
                      <div className="flex items-center gap-2 text-xs text-[#475569] dark:text-[#94a3b8] ml-4">
                        <span className="rounded-full bg-[#ecfeff] px-2 py-1 text-[#036672] dark:bg-[#13263a] dark:text-[#38cde4]">
                          {selectedFile.lines} {t('lines')}
                        </span>
                        <span className="rounded-full bg-emerald-100 px-2 py-1 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300">
                          {formatFileSize(selectedFile.size)}
                        </span>
                        <button
                          onClick={() => {
                            navigator.clipboard.writeText(selectedFile.content);
                          }}
                          className="ml-2 inline-flex items-center gap-1 rounded-full border border-[#e2e8f0] px-3 py-1 text-[#475569] transition hover:border-[#bae6fd] hover:bg-[#ecfeff] dark:border-[#13263a] dark:text-[#94a3b8] dark:hover:bg-[#1b2534]"
                          title={t('copyToClipboard')}
                        >
                          <ClipboardDocumentIcon className="w-4 h-4" /> {t('copy')}
                        </button>
                      </div>
                    </div>
                  </div>

                  {/* Content */}
                  <div className="flex-1 overflow-auto">
                    <AnimatePresence mode="wait">
                      <motion.div
                        key={selectedFile.path + String(wrap) + String(fontSize)}
                        initial={{ opacity: 0, x: 20 }}
                        animate={{ opacity: 1, x: 0 }}
                        exit={{ opacity: 0, x: -20 }}
                        transition={{ duration: 0.15 }}
                      >
                        <CodeViewer
                          code={selectedFile.content}
                          wrap={wrap}
                          fontSize={fontSize}
                          langHint={extOf(selectedFile.path)}
                        />
                      </motion.div>
                    </AnimatePresence>
                  </div>
                </>
              ) : (
                <motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  className="flex-1 flex items-center justify-center text-gray-500 dark:text-gray-400"
                >
                  <EmptyState />
                </motion.div>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Empty initial state */}
      {!projectData && !isLoading && (
        <div className="p-8 text-center bg-[#f9fafb] dark:bg-[#050a12]">
          <FolderIcon className="w-16 h-16 mx-auto mb-4 text-[#06b6d4]/40" />
          <h3 className="text-lg font-medium text-[#111827] dark:text-[#e5f2f2] mb-2">{t('projectNotExtracted')}</h3>
          <p className="text-[#475569] dark:text-[#94a3b8] mb-4">
            {t('extractProjectInstructions')}
          </p>
        </div>
      )}
    </div>
  );
}

/* ========= Subcomponentes ========= */

function EmptyState() {
  const { t } = useTranslations();
  return (
    <div className="text-center">
      <DocumentIcon className="w-20 h-20 mx-auto mb-4 text-[#06b6d4]/40" />
      <h3 className="text-lg font-medium mb-2 text-[#475569] dark:text-[#cbd5f5]">{t('noFileSelected')}</h3>
      <p className="text-sm max-w-xs mx-auto text-[#64748b] dark:text-[#94a3b8]">
        {t('clickFileToView')}
      </p>
      <div className="mt-4 p-3 rounded-lg border border-[#bae6fd] bg-[#ecfeff] text-xs text-[#0f172a] dark:border-[#13263a] dark:bg-[#0a1e2b] dark:text-[#38cde4]">
        💡 <strong>{t('newLayout')}:</strong> {t('newLayoutInstructions')}
      </div>
    </div>
  );
}

function TreeView({
  node,
  onToggle,
  onSelect,
  selectedPath
}: {
  node: TreeNode;
  onToggle: (n: TreeNode) => void;
  onSelect: (f: ProjectFile) => void;
  selectedPath: string;
}) {
  if (!node.children) return null;
  return (
    <ul className="text-sm">
      {node.children.map((n) => {
        const isFile = !!n.file;
        if (isFile) {
          const sel = selectedPath === n.path;
          return (
            <li key={n.path}>
              <button
                onClick={() => n.file && onSelect(n.file)}
                className={classNames(
                  'w-full text-left px-3 py-2 border-b border-transparent dark:border-[#13263a] hover:bg-[#ecfeff] dark:hover:bg-[#1b2534] transition-all',
                  sel && 'bg-[#ecfeff] border-r-2 border-r-[#06b6d4] shadow-[inset_2px_0_0_0_#06b6d4] dark:bg-[#0a1e2b] dark:border-r-[#38cde4]'
                )}
                aria-current={sel ? 'true' : undefined}
              >
                <div className="flex items-center gap-2 mb-1">
                  {getFileIcon(n.path)}
                  <span className={classNames('font-medium truncate', sel ? 'text-[#0f172a] dark:text-[#e5f2f2]' : 'text-[#475569] dark:text-[#94a3b8]')}>
                    {n.name}
                  </span>
                </div>
                <div className="text-xs text-[#94a3b8] dark:text-[#64748b] truncate">{n.path}</div>
              </button>
            </li>
          );
        }
        // folder
        return (
          <li key={n.path}>
            <button
              onClick={() => onToggle(n)}
              className="w-full text-left px-3 py-2 bg-white/70 dark:bg-[#0a1e2b] hover:bg-[#ecfeff] dark:hover:bg-[#1b2534] border-b border-[#e2e8f0] dark:border-[#13263a] flex items-center gap-2"
              aria-expanded={n.open ? 'true' : 'false'}
            >
              <FolderIcon className="w-4 h-4 text-[#06b6d4]" />
              <span className="font-semibold text-[#111827] dark:text-[#e5f2f2]">{n.name}</span>
            </button>
            <AnimatePresence initial={false}>
              {n.open && (
                <motion.div initial={{ height: 0, opacity: 0 }} animate={{ height: 'auto', opacity: 1 }} exit={{ height: 0, opacity: 0 }}>
                  <div className="ml-4 border-l border-[#e2e8f0] dark:border-[#13263a]">
                    <TreeView node={n} onToggle={onToggle} onSelect={onSelect} selectedPath={selectedPath} />
                  </div>
                </motion.div>
              )}
            </AnimatePresence>
          </li>
        );
      })}
    </ul>
  );
}

function CodeViewer({ code, wrap, fontSize, langHint }: { code: string; wrap: boolean; fontSize: number; langHint?: string }) {
  const { t } = useTranslations();
  const html = useMemo(() => maybeHighlight(code, langHint), [code, langHint]);
  return (
    <pre

      title={t('adjustFontSize')}
      className={classNames(
        'p-4 bg-[#f9fafb] dark:bg-[#050a12] text-[#1f2937] dark:text-[#e5f2f2] overflow-auto leading-relaxed',
        wrap ? 'whitespace-pre-wrap break-words' : 'whitespace-pre'
      )}
      style={{ fontSize }}
    >
      {/* se highlight disponível, injeta html; senão, texto puro */}
      <code dangerouslySetInnerHTML={{ __html: html }} />
    </pre>
  );
}
