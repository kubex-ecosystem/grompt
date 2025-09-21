'use client';

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
  if (isFolder) return <FolderIcon className="w-4 h-4 text-blue-500" />;
  if (['js', 'ts', 'jsx', 'tsx', 'go', 'rs', 'py', 'sh', 'json', 'css', 'html', 'md'].includes(ext || '')) {
    return <CodeBracketIcon className="w-4 h-4 text-yellow-500" />;
  }
  return <DocumentIcon className="w-4 h-4 text-gray-500" />;
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
  const [isLoading, setIsLoading] = useState(false);
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
    abortRef.current?.abort();
    const ctrl = new AbortController();
    abortRef.current = ctrl;
    t0.current = performance.now();

    try {
      const response = await fetch(`/api/extract-project?project=${encodeURIComponent(projectFile)}`, {
        signal: ctrl.signal,
        headers: { 'x-kubex-client': 'grompt/pe' }
      });
      const data = await response.json();

      if (data?.success) {
        setProjectData(data);
        setSelectedFile(data.files[0] || null);
      } else {
        setError(data?.error || 'Falha na extração.');
      }
    } catch (err: any) {
      if (err.name !== 'AbortError') setError('Erro ao extrair projeto.');
    } finally {
      setIsLoading(false);
      const dt = (performance.now() - t0.current).toFixed(0);
      // eslint-disable-next-line no-console
      console.log(`[PE] extract done in ${dt}ms — files=${projectData?.files?.length ?? 0}`);
    }
  }, [projectFile]);

  const downloadProject = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch('/api/extract-project', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'x-kubex-client': 'grompt/pe' },
        body: JSON.stringify({ projectFile, format: 'zip' })
      });

      if (response.ok) {
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `${projectName}.zip`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
      } else {
        setError(`Falha no download (${response.status})`);
      }
    } catch {
      setError('Erro ao baixar ZIP.');
    } finally {
      setIsLoading(false);
    }
  }, [projectFile, projectName]);

  const downloadSourceFile = useCallback(async () => {
    setError(null);
    try {
      const response = await fetch(`/projects/${encodeURIComponent(projectFile)}`, {
        headers: { 'x-kubex-client': 'grompt/pe' }
      });
      if (response.ok) {
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        const extSrc = projectFile.split('.').pop() || 'lkt';
        a.download = `${projectName}.${extSrc}`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
      } else {
        setError('Erro ao baixar arquivo de origem.');
      }
    } catch {
      setError('Erro ao baixar arquivo de origem.');
    }
  }, [projectFile, projectName]);

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
    <div className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden bg-white dark:bg-gray-900">
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-500 to-purple-600 text-white p-4 sticky top-0 z-10">
        <div className="flex items-center justify-between gap-3">
          <div className="min-w-0">
            <h3 className="text-lg font-bold truncate">{projectName}</h3>
            {description && <p className="text-sm opacity-90 truncate">{description}</p>}
          </div>
          <div className="flex gap-2">
            <button
              title="Extraction mode"
              onClick={() => setExtractionMode((m) => (m === 'preview' ? 'download' : 'preview'))}
              className="p-2 bg-white/20 rounded-lg hover:bg-white/30 transition-colors"
              aria-pressed={showStats}
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
              placeholder="Buscar (atalho: /)"
              className="w-full pl-8 pr-3 py-2 rounded-md text-sm text-gray-900"
            />
          </div>
          <div className="flex items-center gap-2">
            <FunnelIcon className="w-4 h-4 opacity-80" />
            <select
              value={extFilter}
              onChange={(e) => setExtFilter(e.target.value)}
              className="bg-white/90 text-gray-900 rounded-md px-2 py-1 text-sm"
              title="Filtrar por extensão"
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
                title='Ajustar tamanho da fonte'
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
            className="bg-gray-50 dark:bg-gray-800 p-4 border-b"
          >
            <div className="grid grid-cols-3 gap-4 text-center">
              <div>
                <div className="text-2xl font-bold text-blue-600">{projectData.stats.totalFiles}</div>
                <div className="text-sm text-gray-600 dark:text-gray-400">Arquivos</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-green-600">{formatFileSize(projectData.stats.totalBytes)}</div>
                <div className="text-sm text-gray-600 dark:text-gray-400">Tamanho</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-purple-600">{projectData.stats.totalMarkers}</div>
                <div className="text-sm text-gray-600 dark:text-gray-400">Marcadores</div>
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
      <div className="p-4 border-b bg-gray-50 dark:bg-gray-800">
        <div className="flex gap-2 flex-wrap">
          {!projectData ? (
            <>
              <button
                onClick={extractProject}
                disabled={isLoading}
                className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
              >
                {isLoading ? (
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                ) : (
                  <PlayIcon className="w-4 h-4" />
                )}
                {isLoading ? 'Extraindo...' : 'Extrair Projeto'}
              </button>
              <button
                onClick={downloadSourceFile}
                className="flex items-center gap-2 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
                title="Baixar arquivo original"
              >
                <DocumentIcon className="w-4 h-4" />
                Download origem
              </button>
            </>
          ) : (
            <>
              <button
                onClick={() => setExtractionMode('preview')}
                className={classNames(
                  'flex items-center gap-2 px-4 py-2 rounded-lg transition-colors',
                  extractionMode === 'preview'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-200 dark:bg-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-500'
                )}
              >
                <EyeIcon className="w-4 h-4" />
                Preview
              </button>
              <button
                onClick={downloadProject}
                disabled={isLoading}
                className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50 transition-colors"
                title="Ctrl/Cmd+S"
              >
                {isLoading ? (
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                ) : (
                  <ArrowDownTrayIcon className="w-4 h-4" />
                )}
                Download ZIP
              </button>
              <button
                onClick={downloadSourceFile}
                className="flex items-center gap-2 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
              >
                <DocumentIcon className="w-4 h-4" />
                Download origem
              </button>
              <button
                onClick={() => extractProject()}
                className="flex items-center gap-2 px-3 py-2 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
                title="Reprocessar"
              >
                <ArrowPathIcon className="w-4 h-4" />
                Reprocessar
              </button>
            </>
          )}
        </div>
        {error && (
          <div className="mt-2 text-sm text-red-600">
            {error}{' '}
            <button onClick={() => extractProject()} className="underline">
              Tentar novamente
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
            aria-label="Arquivos do projeto"
          >
            {/* Left: Tree */}
            <div className="col-span-4 border-r dark:border-gray-700 flex flex-col bg-gray-50 dark:bg-gray-800 overflow-hidden">
              <div className="p-3 border-b dark:border-gray-700 bg-gray-100 dark:bg-gray-700">
                <h4 className="font-semibold text-sm text-gray-900 dark:text-white flex items-center gap-2">
                  <FolderIcon className="w-4 h-4 text-blue-500" />
                  Arquivos ({filteredFiles.length})
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
            <div className="col-span-8 flex flex-col bg-white dark:bg-gray-900 overflow-hidden">
              {selectedFile ? (
                <>
                  {/* Header */}
                  <div className="p-3 border-b dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
                    <div className="flex items-center justify-between gap-2">
                      <div className="flex items-center gap-2 min-w-0 flex-1">
                        {getFileIcon(selectedFile.path)}
                        <span className="font-mono text-sm text-gray-900 dark:text-white truncate">
                          {selectedFile.path}
                        </span>
                      </div>
                      <div className="flex items-center gap-2 text-xs text-gray-600 dark:text-gray-400 ml-4">
                        <span className="bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 px-2 py-1 rounded">
                          {selectedFile.lines} linhas
                        </span>
                        <span className="bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 px-2 py-1 rounded">
                          {formatFileSize(selectedFile.size)}
                        </span>
                        <button
                          onClick={() => {
                            navigator.clipboard.writeText(selectedFile.content);
                          }}
                          className="ml-2 inline-flex items-center gap-1 px-2 py-1 border rounded hover:bg-gray-100 dark:hover:bg-gray-700"
                          title="Copiar conteúdo"
                        >
                          <ClipboardDocumentIcon className="w-4 h-4" /> Copiar
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
        <div className="p-8 text-center">
          <FolderIcon className="w-16 h-16 text-gray-400 mx-auto mb-4 opacity-50" />
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">Projeto ainda não extraído</h3>
          <p className="text-gray-600 dark:text-gray-400 mb-4">
            Clique em &quot;Extrair Projeto&quot; para visualizar os arquivos extraídos dos marcadores LookAtni
          </p>
        </div>
      )}
    </div>
  );
}

/* ========= Subcomponentes ========= */

function EmptyState() {
  return (
    <div className="text-center">
      <DocumentIcon className="w-20 h-20 mx-auto mb-4 opacity-30" />
      <h3 className="text-lg font-medium mb-2 text-gray-700 dark:text-gray-300">Nenhum arquivo selecionado</h3>
      <p className="text-sm max-w-xs mx-auto">
        Clique em um arquivo na lista à esquerda para visualizar seu conteúdo completo
      </p>
      <div className="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800 text-blue-700 dark:text-blue-300 text-xs">
        💡 <strong>Novo layout:</strong> Lista à esquerda, conteúdo detalhado à direita. Use <kbd>/</kbd> para buscar e <kbd>j/k</kbd> para navegar.
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
                  'w-full text-left px-3 py-2 border-b dark:border-gray-700 hover:bg-gray-100 dark:hover:bg-gray-700 transition-all',
                  sel && 'bg-blue-100 dark:bg-blue-900/30 border-r-2 border-r-blue-500 shadow-sm'
                )}
                aria-current={sel ? 'true' : undefined}
              >
                <div className="flex items-center gap-2 mb-1">
                  {getFileIcon(n.path)}
                  <span className={classNames('font-medium truncate', sel ? 'text-blue-700 dark:text-blue-300' : '')}>
                    {n.name}
                  </span>
                </div>
                <div className="text-xs text-gray-500 dark:text-gray-400 truncate">{n.path}</div>
              </button>
            </li>
          );
        }
        // folder
        return (
          <li key={n.path}>
            <button
              onClick={() => onToggle(n)}
              className="w-full text-left px-3 py-2 bg-gray-100/70 dark:bg-gray-700/30 hover:bg-gray-100 dark:hover:bg-gray-700 border-b dark:border-gray-700 flex items-center gap-2"
              aria-expanded={n.open ? 'true' : 'false'}
            >
              <FolderIcon className="w-4 h-4 text-blue-500" />
              <span className="font-semibold">{n.name}</span>
            </button>
            <AnimatePresence initial={false}>
              {n.open && (
                <motion.div initial={{ height: 0, opacity: 0 }} animate={{ height: 'auto', opacity: 1 }} exit={{ height: 0, opacity: 0 }}>
                  <div className="ml-4 border-l dark:border-gray-700">
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
  const html = useMemo(() => maybeHighlight(code, langHint), [code, langHint]);
  return (
    <pre

      title='Ajustar tamanho da fonte'
      className={classNames(
        'p-4 bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-200 overflow-auto leading-relaxed',
        wrap ? 'whitespace-pre-wrap break-words' : 'whitespace-pre'
      )}
      style={{ fontSize }}
    >
      {/* se highlight disponível, injeta html; senão, texto puro */}
      <code dangerouslySetInnerHTML={{ __html: html }} />
    </pre>
  );
}
