'use client';

// import { AlertCircle, CheckCircle, FileText, FolderOpen, Loader2 } from 'lucide-react';
// import { useRef, useState } from 'react';
// import { useTranslation } from 'react-i18next';

// // Importar estilos CSS module
// import styles from './LookatniGenerator.module.css';

// // Tipos locais
// interface LookatniGeneratorProps {
//   darkMode: boolean;
//   currentTheme: any;
// }

// interface ProcessingStats {
//   totalFiles: number;
//   processedFiles: number;
//   currentFile: string;
//   startTime: number;
// }

// // Componente auxiliar para barra de progresso
// const ProgressBar = ({ progress }: { progress: number }) => {
//   const safeProgress = Math.min(100, Math.max(0, progress));
//   return (
//     <div className={styles.progressBar}>
//       <div
//         className={styles.progressFill}
//         style={{ width: `${safeProgress}%` }}
//       />
//     </div>
//   );
// };

// export default function LookatniGenerator({ darkMode, currentTheme }: LookatniGeneratorProps) {
//   const { t } = useTranslation();
//   const [isProcessing, setIsProcessing] = useState(false);
//   const [processingStats, setProcessingStats] = useState<ProcessingStats | null>(null);
//   const [lastResult, setLastResult] = useState<string | null>(null);
//   const [error, setError] = useState<string | null>(null);
//   const fileInputRef = useRef<HTMLInputElement>(null);

//   // Função para verificar se o browser suporta File System Access API
//   const supportsFileSystemAccess = () => {
//     return 'showDirectoryPicker' in window;
//   };

//   // Função para ler arquivos usando File System Access API (Chrome/Edge moderno)
//   const readDirectoryModern = async (): Promise<File[]> => {
//     try {
//       // @ts-ignore - showDirectoryPicker ainda não tem tipos oficiais
//       const dirHandle = await window.showDirectoryPicker({
//         mode: 'read'
//       });

//       const files: File[] = [];

//       async function processDirectory(dirHandle: any, path = '') {
//         for await (const [name, handle] of dirHandle.entries()) {
//           const fullPath = path ? `${path}/${name}` : name;

//           if (handle.kind === 'file') {
//             try {
//               const file = await handle.getFile();
//               // Adicionar path info ao arquivo
//               Object.defineProperty(file, 'webkitRelativePath', {
//                 value: fullPath,
//                 writable: false
//               });
//               files.push(file);
//             } catch (error) {
//               console.warn(`Erro ao ler arquivo ${fullPath}:`, error);
//             }
//           } else if (handle.kind === 'directory') {
//             await processDirectory(handle, fullPath);
//           }
//         }
//       }

//       await processDirectory(dirHandle);
//       return files;
//     } catch (error) {
//       if (error.name === 'AbortError') {
//         throw new Error('Seleção de pasta cancelada pelo usuário');
//       }
//       throw error;
//     }
//   };

//   // Função para ler arquivos usando input file (fallback)
//   const readDirectoryFallback = (): Promise<File[]> => {
//     return new Promise((resolve, reject) => {
//       if (!fileInputRef.current) {
//         reject(new Error('Input de arquivo não encontrado'));
//         return;
//       }

//       const handleFileSelect = (event: Event) => {
//         const target = event.target as HTMLInputElement;
//         const files = Array.from(target.files || []);

//         // Limpar o event listener
//         fileInputRef.current?.removeEventListener('change', handleFileSelect);

//         if (files.length === 0) {
//           reject(new Error('Nenhum arquivo selecionado'));
//         } else {
//           resolve(files);
//         }
//       };

//       fileInputRef.current.addEventListener('change', handleFileSelect);
//       fileInputRef.current.click();
//     });
//   };

//   // Função para filtrar arquivos desnecessários
//   const shouldExcludeFile = (filePath: string): boolean => {
//     const excludePatterns = [
//       // Diretórios comuns a ignorar
//       /node_modules\//,
//       /\.git\//,
//       /\.next\//,
//       /dist\//,
//       /build\//,
//       /coverage\//,
//       /\.vscode\//,
//       /\.idea\//,
//       /tmp\//,
//       /temp\//,

//       // Tipos de arquivo a ignorar
//       /\.(log|tmp|cache|lock)$/,
//       /\.env(\.|$)/,
//       /package-lock\.json$/,
//       /yarn\.lock$/,
//       /pnpm-lock\.yaml$/,

//       // Arquivos binários comuns
//       /\.(png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$/i,
//       /\.(mp4|avi|mov|wmv|mp3|wav|pdf)$/i,
//       /\.(zip|rar|tar|gz|7z)$/i,

//       // Arquivos muito grandes (será verificado por tamanho também)
//       /\.min\.(js|css)$/
//     ];

//     return excludePatterns.some(pattern => pattern.test(filePath));
//   };

//   // Função para verificar se arquivo é muito grande
//   const isFileTooLarge = (file: File): boolean => {
//     const maxSize = 1024 * 1024; // 1MB
//     return file.size > maxSize;
//   };

//   // Função principal para processar diretório
//   const processDirectory = async () => {
//     setIsProcessing(true);
//     setError(null);
//     setLastResult(null);

//     const startTime = Date.now();

//     try {
//       console.log('🚀 Iniciando processamento do diretório...');

//       // Selecionar arquivos baseado no suporte do browser
//       let files: File[];

//       if (supportsFileSystemAccess()) {
//         console.log('📁 Usando File System Access API');
//         files = await readDirectoryModern();
//       } else {
//         console.log('📂 Usando fallback de input file');
//         files = await readDirectoryFallback();
//       }

//       console.log(`📊 ${files.length} arquivos encontrados`);

//       // Filtrar arquivos desnecessários
//       const filteredFiles = files.filter(file => {
//         const relativePath = file.webkitRelativePath || file.name;

//         if (shouldExcludeFile(relativePath)) {
//           console.log(`⏭️ Ignorando: ${relativePath} (filtro de padrão)`);
//           return false;
//         }

//         if (isFileTooLarge(file)) {
//           console.log(`⏭️ Ignorando: ${relativePath} (muito grande: ${(file.size / 1024).toFixed(1)}KB)`);
//           return false;
//         }

//         return true;
//       });

//       console.log(`✅ ${filteredFiles.length} arquivos após filtros (${files.length - filteredFiles.length} ignorados)`);

//       // Inicializar stats de processamento
//       setProcessingStats({
//         totalFiles: filteredFiles.length,
//         processedFiles: 0,
//         currentFile: '',
//         startTime
//       });

//       // Preparar estrutura de arquivos para o MarkerGenerator
//       const fileStructure: { [path: string]: string } = {};

//       // Processar arquivos um por um
//       for (let i = 0; i < filteredFiles.length; i++) {
//         const file = filteredFiles[i];
//         const relativePath = file.webkitRelativePath || file.name;

//         setProcessingStats(prev => ({
//           ...prev!,
//           processedFiles: i,
//           currentFile: relativePath
//         }));

//         try {
//           // Ler conteúdo do arquivo
//           const content = await readFileContent(file);
//           fileStructure[relativePath] = content;

//           console.log(`✅ Processado: ${relativePath} (${file.size} bytes)`);
//         } catch (error) {
//           console.warn(`⚠️ Erro ao ler ${relativePath}:`, error);
//           // Continuar mesmo com erro em arquivo individual
//         }

//         // Pequeno delay para não travar a UI
//         if (i % 10 === 0) {
//           await new Promise(resolve => setTimeout(resolve, 1));
//         }
//       }

//       // Gerar o arquivo .marked usando o MarkerGenerator
//       console.log('🔄 Gerando arquivo .marked...');

//       setProcessingStats(prev => ({
//         ...prev!,
//         currentFile: 'Gerando arquivo .marked...'
//       }));

//       // Criar um "projeto virtual" para o generator
//       const projectName = `projeto-${Date.now()}`;

//       // Criar estrutura temporária de arquivos para simular um diretório
//       const tempDir = `/tmp/${projectName}`;

//       // Por enquanto, vamos criar um .marked usando WASM (ultra-performance!)
//       const markedContent = await generateMarkedContentWasm(fileStructure, {
//         projectName,
//         totalFiles: filteredFiles.length,
//         originalTotal: files.length,
//         filteredOut: files.length - filteredFiles.length,
//         generatedBy: 'Grompt + LookAtni JavaScript Fallback',
//         generatedAt: new Date().toISOString()
//       });

//       // Criar e baixar o arquivo
//       const blob = new Blob([markedContent], { type: 'text/plain' });
//       const url = URL.createObjectURL(blob);

//       const a = document.createElement('a');
//       a.href = url;
//       a.download = `${projectName}.marked`;
//       document.body.appendChild(a);
//       a.click();
//       document.body.removeChild(a);

//       // Limpar URL
//       URL.revokeObjectURL(url);

//       const processingTime = Date.now() - startTime;
//       const resultMessage = `✅ Arquivo gerado com sucesso!
// 📦 ${filteredFiles.length} arquivos processados (${files.length - filteredFiles.length} filtrados)
// 📊 Total encontrado: ${files.length} arquivos
// ⏱️ ${Math.round(processingTime / 1000)}s de processamento
// 📥 Download iniciado: ${projectName}.marked`;

//       setLastResult(resultMessage);
//       console.log('🎉 Processamento concluído:', resultMessage);

//     } catch (error) {
//       const errorMessage = error instanceof Error ? error.message : 'Erro desconhecido';
//       setError(`❌ Erro durante o processamento: ${errorMessage}`);
//       console.error('💥 Erro no processamento:', error);
//     } finally {
//       setIsProcessing(false);
//       setProcessingStats(null);
//     }
//   };

//   // Função para gerar conteúdo .marked (usar JS fallback enquanto WASM não carrega)
//   const generateMarkedContentWasm = async (
//     fileStructure: { [path: string]: string },
//     metadata: any
//   ): Promise<string> => {
//     const wasmParser = await import('../wasm/lookatni_wasm_parser_bg.wasm');

//     if (!wasmParser) {
//       console.log('⚠️ WASM não carregado, usando fallback JavaScript');
//       return generateMarkedContentFallback(fileStructure, metadata);
//     }

//     try {
//       // Preparar dados para o WASM
//       const projectData = {
//         metadata,
//         files: fileStructure
//       };

//       // Usar o parser WASM para gerar o conteúdo .marked
//       const result = wasmParser.parse_ai_code(JSON.stringify(projectData));

//       return result;
//     } catch (error) {
//       console.error('❌ Erro no processamento WASM:', error);
//       // Fallback para JavaScript se WASM falhar
//       return generateMarkedContentFallback(fileStructure, metadata);
//     }
//   };

//   // Função fallback em JavaScript puro - FORMATO ORIGINAL DO LOOKATNI CLI
//   const generateMarkedContentFallback = async (
//     fileStructure: { [path: string]: string },
//     metadata: any
//   ): Promise<string> => {
//     console.log('🔄 Usando fallback JavaScript - FORMATO ORIGINAL do LookAtni CLI');

//     // Cabeçalho no formato EXATO do LookAtni CLI original
//     const fileEntries = Object.entries(fileStructure);
//     const header = `# LookAtni Code - Gerado automaticamente
// # Data: ${metadata.generatedAt}
// # Fonte: ./
// # Total de arquivos: ${fileEntries.length}

// `;

//     let content = header;

//     // ASCII 28 (File Separator) - ESSENCIAL para o parser do LookAtni
//     const FILE_SEPARATOR = String.fromCharCode(28); // \034 em octal

//     for (const [filePath, fileContent] of fileEntries) {
//       // Formato EXATO: //<ASCII28>/ FILENAME /<ASCII28>//
//       content += `//${FILE_SEPARATOR}/ ${filePath} /${FILE_SEPARATOR}//
// ${fileContent}

// `;
//     } return content;
//   };  // Função auxiliar para ler conteúdo de arquivo
//   const readFileContent = (file: File): Promise<string> => {
//     return new Promise((resolve, reject) => {
//       const reader = new FileReader();

//       reader.onload = (event) => {
//         resolve(event.target?.result as string);
//       };

//       reader.onerror = () => {
//         reject(new Error(`Erro ao ler arquivo: ${file.name}`));
//       };

//       // Tentar ler como texto UTF-8
//       reader.readAsText(file, 'utf-8');
//     });
//   };

//   // Função para limpar resultados
//   const clearResults = () => {
//     setLastResult(null);
//     setError(null);
//   };

//   return (
//     <div className={`${currentTheme.cardBg} rounded-xl border ${currentTheme.border} shadow-lg p-6`}>
//       {/* Input file oculto para fallback */}
//       <input
//         ref={fileInputRef}
//         type="file"
//         {...({ webkitdirectory: "", directory: "" } as any)}
//         multiple
//         className="hidden"
//         aria-label="Seletor de pasta para processamento"
//       />

//       {/* Header */}
//       <div className="flex items-center space-x-3 mb-6">
//         <div className="p-2 bg-gradient-to-r from-green-500 to-blue-600 rounded-lg">
//           <FileText className="h-6 w-6 text-white" />
//         </div>
//         <div>
//           <h3 className="text-lg font-semibold flex items-center space-x-2">
//             <span>LookAtni Generator</span>
//             <span className="px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-xs rounded-full flex items-center space-x-1">
//               <FileText className="h-3 w-3" />
//               <span>JS Ready</span>
//             </span>
//           </h3>
//           <p className="text-sm text-gray-500">
//             Transforme qualquer pasta em arquivo .marked - JavaScript Engine ativo!
//           </p>
//         </div>
//       </div>      {/* Botão principal */}
//       {!isProcessing && (
//         <div className="space-y-4">
//           <button
//             onClick={processDirectory}
//             disabled={isProcessing}
//             className={`w-full flex items-center justify-center space-x-3 px-6 py-4 ${currentTheme.button} rounded-lg font-medium transition-all hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed`}
//           >
//             <FolderOpen className="h-5 w-5" />
//             <span>Selecionar Pasta e Gerar .marked</span>
//             <FileText className="h-4 w-4" />
//           </button>

//           <div className={`text-xs ${currentTheme.textSecondary} bg-gray-100 dark:bg-gray-800 rounded-lg p-3`}>
//             <p className="font-medium mb-1">🚀 JavaScript Engine + 🔒 100% Local:</p>
//             <ul className="space-y-1 text-xs">
//               <li>• Seus arquivos não saem do seu computador</li>
//               <li>• Todo processamento acontece no seu navegador</li>
//               <li>• <strong>Engine JavaScript otimizada e funcional</strong></li>
//               <li>• O arquivo .marked é gerado e baixado diretamente</li>
//               <li>• Filtros inteligentes (exclui node_modules, .git, etc.)</li>
//             </ul>
//           </div>
//         </div>
//       )}

//       {/* Status do processamento */}
//       {isProcessing && processingStats && (
//         <div className="space-y-4">
//           <div className="flex items-center space-x-3">
//             <Loader2 className="h-5 w-5 animate-spin text-blue-500" />
//             <span className="font-medium">Processando arquivos...</span>
//           </div>

//           <div className="space-y-2">
//             <div className="flex justify-between text-sm">
//               <span>Progresso:</span>
//               <span>{processingStats.processedFiles} / {processingStats.totalFiles}</span>
//             </div>

//             <ProgressBar progress={(processingStats.processedFiles / processingStats.totalFiles) * 100} />            <div className="text-xs text-gray-500 truncate">
//               📄 {processingStats.currentFile}
//             </div>

//             <div className="text-xs text-gray-500">
//               ⏱️ {Math.round((Date.now() - processingStats.startTime) / 1000)}s decorridos
//             </div>
//           </div>
//         </div>
//       )}

//       {/* Resultado de sucesso */}
//       {lastResult && (
//         <div className="space-y-3">
//           <div className="flex items-start space-x-3 p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
//             <CheckCircle className="h-5 w-5 text-green-500 mt-0.5" />
//             <div className="flex-1">
//               <div className="font-medium text-green-800 dark:text-green-200 mb-2">
//                 Processamento Concluído!
//               </div>
//               <pre className="text-xs text-green-700 dark:text-green-300 whitespace-pre-wrap">
//                 {lastResult}
//               </pre>
//             </div>
//           </div>

//           <button
//             onClick={clearResults}
//             className={`text-sm ${currentTheme.buttonSecondary} px-4 py-2 rounded-lg`}
//           >
//             Processar Nova Pasta
//           </button>
//         </div>
//       )}

//       {/* Erro */}
//       {error && (
//         <div className="space-y-3">
//           <div className="flex items-start space-x-3 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
//             <AlertCircle className="h-5 w-5 text-red-500 mt-0.5" />
//             <div className="flex-1">
//               <div className="font-medium text-red-800 dark:text-red-200 mb-2">
//                 Erro no Processamento
//               </div>
//               <div className="text-sm text-red-700 dark:text-red-300">
//                 {error}
//               </div>
//             </div>
//           </div>

//           <button
//             onClick={clearResults}
//             className={`text-sm ${currentTheme.buttonSecondary} px-4 py-2 rounded-lg`}
//           >
//             Tentar Novamente
//           </button>
//         </div>
//       )}

//       {/* Info sobre compatibilidade */}
//       <div className="mt-4 text-xs text-gray-500 border-t pt-4">
//         <p className="mb-1">
//           <strong>Compatibilidade:</strong>
//         </p>
//         <ul className="space-y-1">
//           <li>
//             • <strong>Chrome/Edge:</strong> Seleção de pasta nativa
//           </li>
//           <li>
//             • <strong>Firefox/Safari:</strong> Seleção de arquivos (atributo webkitdirectory)
//           </li>
//           <li>
//             • <strong>Mobile:</strong> Limitado pelo suporte do navegador
//           </li>
//         </ul>
//       </div>
//     </div>
//   );
// }
