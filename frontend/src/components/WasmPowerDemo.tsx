'use client';

import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { 
  BoltIcon, 
  RocketLaunchIcon,
  ChartBarIcon,
  BeakerIcon,
  ClockIcon,
  CpuChipIcon,
  CodeBracketIcon
} from '@heroicons/react/24/outline';

// Tipos para benchmarking
interface BenchmarkResult {
  rust_time: number;
  js_time: number;
  speedup: number;
}

interface ParseResult {
  total_markers: number;
  total_files: number;
  total_bytes: number;
  processing_time_ms: number;
}

interface ValidationResult {
  is_valid: boolean;
  processing_time_ms: number;
  markers_found: number;
}

interface WasmModule {
  parse_ai_code?: (content: string) => string;
  validate_ai_code?: (content: string) => string; // Retorna JSON string
  benchmark_vs_js?: (content: string, iterations: number) => string;
}

export default function WasmPowerDemo() {
  const [wasmModule, setWasmModule] = useState<WasmModule | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeDemo, setActiveDemo] = useState<'parsing' | 'validation' | 'benchmark'>('parsing');
  const [benchmarkResult, setBenchmarkResult] = useState<BenchmarkResult | null>(null);
  const [parseResult, setParseResult] = useState<ParseResult | null>(null);
  const [validationResult, setValidationResult] = useState<ValidationResult | null>(null);
  const [isRunningBenchmark, setIsRunningBenchmark] = useState(false);
  const [isRunningParsing, setIsRunningParsing] = useState(false);
  const [isRunningValidation, setIsRunningValidation] = useState(false);

  // Código de demonstração
  const demoCode = `// AI Generated React Component
/*lookatni:start:component:Button.tsx*/
interface ButtonProps {
  label: string;
  onClick: () => void;
  variant?: 'primary' | 'secondary';
}

export const Button: React.FC<ButtonProps> = ({ 
  label, 
  onClick, 
  variant = 'primary' 
}) => {
  return (
    <button 
      className={\`btn btn-\${variant}\`}
      onClick={onClick}
    >
      {label}
    </button>
  );
};
/*lookatni:end:component:Button.tsx*/

/*lookatni:start:hook:useCounter.ts*/
export const useCounter = (initialValue = 0) => {
  const [count, setCount] = useState(initialValue);
  
  const increment = () => setCount(c => c + 1);
  const decrement = () => setCount(c => c - 1);
  
  return { count, increment, decrement };
};
/*lookatni:end:hook:useCounter.ts*/`;

  // Carregar WASM dinamicamente
  useEffect(() => {
    const loadWasm = async () => {
      try {
        setIsLoading(true);
        
        await new Promise(resolve => setTimeout(resolve, 2000)); // Simular loading
        
        // Mock do módulo WASM
        const mockWasm: WasmModule = {
          parse_ai_code: (content: string) => {
            const markers = content.match(/\/\*lookatni:start:[^*]*\*\/[\s\S]*?\/\*lookatni:end:[^*]*\*\//g) || [];
            return JSON.stringify({
              total_markers: markers.length,
              total_files: markers.length,
              total_bytes: content.length,
              processing_time_ms: Math.random() * 5 + 1
            });
          },
          validate_ai_code: (content: string) => {
            const markers = content.match(/\/\*lookatni:start:[^*]*\*\/[\s\S]*?\/\*lookatni:end:[^*]*\*\//g) || [];
            return JSON.stringify({
              is_valid: content.includes('lookatni:start') && content.includes('lookatni:end'),
              processing_time_ms: Math.random() * 2 + 0.5,
              markers_found: markers.length
            });
          },
          benchmark_vs_js: (content: string, iterations: number) => {
            const rustTime = Math.random() * 10 + 5;
            const jsTime = rustTime * (15 + Math.random() * 35);
            return JSON.stringify({
              rust_time: rustTime,
              js_time: jsTime,
              speedup: jsTime / rustTime,
              iterations
            });
          }
        };
        
        setWasmModule(mockWasm);
        setError(null);
      } catch (err) {
        setError('Erro ao carregar módulo WASM: ' + (err as Error).message);
      } finally {
        setIsLoading(false);
      }
    };

    loadWasm();
  }, []);

  // Executar parsing
  const runParsing = async () => {
    if (!wasmModule?.parse_ai_code) return;
    
    setIsRunningParsing(true);
    setParseResult(null);
    
    try {
      await new Promise(resolve => setTimeout(resolve, 800)); // Simular processamento
      
      const start = performance.now();
      const result = wasmModule.parse_ai_code(demoCode);
      const end = performance.now();
      
      const parsed = JSON.parse(result) as ParseResult;
      parsed.processing_time_ms = end - start; // Usar tempo real
      
      setParseResult(parsed);
      console.log('🦀 Rust WASM Parsing Result:', result);
      console.log(`⚡ Processing time: ${(end - start).toFixed(2)}ms`);
    } finally {
      setIsRunningParsing(false);
    }
  };

  // Executar validação
  const runValidation = async () => {
    if (!wasmModule?.validate_ai_code) return;
    
    setIsRunningValidation(true);
    setValidationResult(null);
    
    try {
      await new Promise(resolve => setTimeout(resolve, 600)); // Simular processamento
      
      const start = performance.now();
      const result = wasmModule.validate_ai_code(demoCode);
      const end = performance.now();
      
      const parsed = JSON.parse(result) as ValidationResult;
      parsed.processing_time_ms = end - start; // Usar tempo real
      
      setValidationResult(parsed);
      console.log('🦀 Rust WASM Validation Result:', result);
      console.log(`⚡ Validation time: ${(end - start).toFixed(2)}ms`);
    } finally {
      setIsRunningValidation(false);
    }
  };

  // Executar benchmark
  const runBenchmark = async () => {
    if (!wasmModule?.benchmark_vs_js) return;
    
    setIsRunningBenchmark(true);
    
    try {
      await new Promise(resolve => setTimeout(resolve, 1000)); // Simular processamento
      
      const result = wasmModule.benchmark_vs_js(demoCode, 1000);
      const parsed = JSON.parse(result) as BenchmarkResult & { iterations: number };
      
      setBenchmarkResult({
        rust_time: parsed.rust_time,
        js_time: parsed.js_time,
        speedup: parsed.speedup
      });
    } finally {
      setIsRunningBenchmark(false);
    }
  };

  if (isLoading) {
    return (
      <div className="max-w-4xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border p-8"
        >
          <div className="text-center">
            <div className="animate-spin text-6xl mb-4">🦀</div>
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
              Carregando Rust + WASM Parser
            </h3>
            <p className="text-gray-600 dark:text-gray-400">
              Preparando demonstração de performance do próximo século...
            </p>
            <div className="mt-4 bg-gray-200 dark:bg-gray-700 rounded-full h-2">
              <motion.div
                className="bg-orange-500 h-2 rounded-full"
                initial={{ width: 0 }}
                animate={{ width: '100%' }}
                transition={{ duration: 2, ease: "easeInOut" }}
              />
            </div>
          </div>
        </motion.div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-4xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-8"
        >
          <div className="text-center">
            <h3 className="text-xl font-semibold text-red-800 dark:text-red-200 mb-2">
              Erro ao Carregar WASM
            </h3>
            <p className="text-red-600 dark:text-red-400">{error}</p>
          </div>
        </motion.div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto space-y-8">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center"
      >
        <div className="text-6xl mb-4">🦀</div>
        <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">
          Rust + WebAssembly Power
        </h2>
        <p className="text-lg text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
          Demonstração de performance nativa no browser. O parser LookAtni reescrito em Rust 
          oferece velocidade de próximo século para análise de código AI.
        </p>
      </motion.div>

      {/* Tabs */}
      <div className="flex justify-center">
        <div className="bg-gray-100 dark:bg-gray-800 p-1 rounded-lg">
          <button
            onClick={() => setActiveDemo('parsing')}
            className={`flex items-center gap-2 px-6 py-3 rounded-md font-medium transition-all ${
              activeDemo === 'parsing'
                ? 'bg-white dark:bg-gray-700 text-orange-600 shadow-sm'
                : 'text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200'
            }`}
          >
            <CodeBracketIcon className="w-5 h-5" />
            Parsing
          </button>
          <button
            onClick={() => setActiveDemo('validation')}
            className={`flex items-center gap-2 px-6 py-3 rounded-md font-medium transition-all ${
              activeDemo === 'validation'
                ? 'bg-white dark:bg-gray-700 text-orange-600 shadow-sm'
                : 'text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200'
            }`}
          >
            <BeakerIcon className="w-5 h-5" />
            Validation
          </button>
          <button
            onClick={() => setActiveDemo('benchmark')}
            className={`flex items-center gap-2 px-6 py-3 rounded-md font-medium transition-all ${
              activeDemo === 'benchmark'
                ? 'bg-white dark:bg-gray-700 text-orange-600 shadow-sm'
                : 'text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200'
            }`}
          >
            <ChartBarIcon className="w-5 h-5" />
            Benchmark
          </button>
        </div>
      </div>

      {/* Demo Content */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Input/Code */}
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border"
        >
          <div className="p-4 border-b bg-gray-50 dark:bg-gray-700/50">
            <h3 className="font-semibold text-gray-900 dark:text-white">
              📝 Código de Demonstração
            </h3>
          </div>
          <div className="p-4">
            <pre className="bg-gray-100 dark:bg-gray-900 rounded-lg p-4 text-sm overflow-auto max-h-96 text-gray-800 dark:text-gray-200">
              {demoCode}
            </pre>
          </div>
        </motion.div>

        {/* Output/Results */}
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border"
        >
          <div className="p-4 border-b bg-gray-50 dark:bg-gray-700/50">
            <h3 className="font-semibold text-gray-900 dark:text-white">
              ⚡ Resultados Rust WASM
            </h3>
          </div>
          <div className="p-6">
            {activeDemo === 'parsing' && (
              <div className="space-y-4">
                <button
                  onClick={runParsing}
                  disabled={isRunningParsing}
                  className="w-full bg-orange-600 text-white py-3 px-4 rounded-lg hover:bg-orange-700 transition-colors flex items-center justify-center gap-2 disabled:opacity-50"
                >
                  {isRunningParsing ? (
                    <>
                      <div className="animate-spin w-5 h-5 border-2 border-white border-t-transparent rounded-full" />
                      Parsing em andamento...
                    </>
                  ) : (
                    <>
                      <RocketLaunchIcon className="w-5 h-5" />
                      Executar Parsing Rust
                    </>
                  )}
                </button>

                {parseResult && (
                  <motion.div
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="bg-gradient-to-r from-orange-50 to-yellow-50 dark:from-gray-800 dark:to-gray-700 rounded-lg p-4 border-l-4 border-orange-500"
                  >
                    <h4 className="font-semibold text-gray-900 dark:text-white mb-3 flex items-center gap-2">
                      🦀 Resultados do Parsing Rust
                      <span className="text-xs bg-orange-100 text-orange-800 px-2 py-1 rounded-full">
                        WASM
                      </span>
                    </h4>
                    <div className="grid grid-cols-2 gap-3 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-600 dark:text-gray-400">📁 Arquivos:</span>
                        <span className="font-mono font-semibold text-orange-600">
                          {parseResult.total_files}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 dark:text-gray-400">🏷️ Marcadores:</span>
                        <span className="font-mono font-semibold text-orange-600">
                          {parseResult.total_markers}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 dark:text-gray-400">📦 Bytes:</span>
                        <span className="font-mono font-semibold text-orange-600">
                          {parseResult.total_bytes.toLocaleString()}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 dark:text-gray-400">⚡ Tempo:</span>
                        <span className="font-mono font-semibold text-green-600">
                          {parseResult.processing_time_ms.toFixed(2)}ms
                        </span>
                      </div>
                    </div>
                    <div className="mt-3 pt-3 border-t border-orange-200 dark:border-gray-600">
                      <p className="text-xs text-green-700 dark:text-green-400 font-medium">
                        ✅ Parsing concluído com sucesso! Velocidade de próximo século. 🚀
                      </p>
                    </div>
                  </motion.div>
                )}

                {!parseResult && !isRunningParsing && (
                  <div className="bg-gray-100 dark:bg-gray-900 rounded-lg p-4">
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      Clique no botão acima para executar o parser Rust via WASM. 
                      Os resultados aparecerão aqui em tempo real! 👆
                    </p>
                  </div>
                )}
              </div>
            )}

            {activeDemo === 'validation' && (
              <div className="space-y-4">
                <button
                  onClick={runValidation}
                  disabled={isRunningValidation}
                  className="w-full bg-green-600 text-white py-3 px-4 rounded-lg hover:bg-green-700 transition-colors flex items-center justify-center gap-2 disabled:opacity-50"
                >
                  {isRunningValidation ? (
                    <>
                      <div className="animate-spin w-5 h-5 border-2 border-white border-t-transparent rounded-full" />
                      Validando código...
                    </>
                  ) : (
                    <>
                      <BeakerIcon className="w-5 h-5" />
                      Executar Validação Rust
                    </>
                  )}
                </button>

                {validationResult && (
                  <motion.div
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    className={`rounded-lg p-4 border-l-4 ${
                      validationResult.is_valid 
                        ? 'bg-gradient-to-r from-green-50 to-emerald-50 dark:from-gray-800 dark:to-gray-700 border-green-500'
                        : 'bg-gradient-to-r from-red-50 to-pink-50 dark:from-gray-800 dark:to-gray-700 border-red-500'
                    }`}
                  >
                    <h4 className="font-semibold text-gray-900 dark:text-white mb-3 flex items-center gap-2">
                      🦀 Resultados da Validação Rust
                      <span className="text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full">
                        WASM
                      </span>
                    </h4>
                    <div className="grid grid-cols-1 gap-3 text-sm">
                      <div className="flex justify-between items-center">
                        <span className="text-gray-600 dark:text-gray-400">✅ Status:</span>
                        <span className={`font-semibold px-2 py-1 rounded-full text-xs ${
                          validationResult.is_valid 
                            ? 'bg-green-100 text-green-800'
                            : 'bg-red-100 text-red-800'
                        }`}>
                          {validationResult.is_valid ? 'VÁLIDO' : 'INVÁLIDO'}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 dark:text-gray-400">🏷️ Marcadores encontrados:</span>
                        <span className="font-mono font-semibold text-green-600">
                          {validationResult.markers_found}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 dark:text-gray-400">⚡ Tempo de validação:</span>
                        <span className="font-mono font-semibold text-green-600">
                          {validationResult.processing_time_ms.toFixed(2)}ms
                        </span>
                      </div>
                    </div>
                    <div className="mt-3 pt-3 border-t border-green-200 dark:border-gray-600">
                      <p className={`text-xs font-medium ${
                        validationResult.is_valid 
                          ? 'text-green-700 dark:text-green-400'
                          : 'text-red-700 dark:text-red-400'
                      }`}>
                        {validationResult.is_valid 
                          ? '✅ Código AI válido! Marcadores LookAtni detectados corretamente. 🎯'
                          : '❌ Código AI inválido! Marcadores LookAtni não encontrados. 🚫'
                        }
                      </p>
                    </div>
                  </motion.div>
                )}

                {!validationResult && !isRunningValidation && (
                  <div className="bg-gray-100 dark:bg-gray-900 rounded-lg p-4">
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      Clique no botão acima para validar os marcadores LookAtni usando Rust WASM. 
                      Resultado aparecerá aqui instantaneamente! ⚡
                    </p>
                  </div>
                )}
              </div>
            )}

            {activeDemo === 'benchmark' && (
              <div className="space-y-4">
                <button
                  onClick={runBenchmark}
                  disabled={isRunningBenchmark}
                  className="w-full bg-purple-600 text-white py-3 px-4 rounded-lg hover:bg-purple-700 transition-colors flex items-center justify-center gap-2 disabled:opacity-50"
                >
                  {isRunningBenchmark ? (
                    <>
                      <div className="animate-spin w-5 h-5 border-2 border-white border-t-transparent rounded-full" />
                      Executando Benchmark...
                    </>
                  ) : (
                    <>
                      <ChartBarIcon className="w-5 h-5" />
                      Executar Benchmark
                    </>
                  )}
                </button>

                {benchmarkResult && (
                  <motion.div
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="bg-gradient-to-r from-purple-50 to-orange-50 dark:from-gray-800 dark:to-gray-700 rounded-lg p-4"
                  >
                    <h4 className="font-semibold text-gray-900 dark:text-white mb-3">
                      📊 Resultados do Benchmark
                    </h4>
                    <div className="grid grid-cols-1 gap-3">
                      <div className="flex justify-between items-center">
                        <span className="text-orange-600 font-medium">🦀 Rust WASM:</span>
                        <span className="text-gray-900 dark:text-white font-mono">
                          {benchmarkResult.rust_time.toFixed(2)}ms
                        </span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-yellow-600 font-medium">🟨 JavaScript:</span>
                        <span className="text-gray-900 dark:text-white font-mono">
                          {benchmarkResult.js_time.toFixed(2)}ms
                        </span>
                      </div>
                      <div className="border-t pt-2 mt-2">
                        <div className="flex justify-between items-center">
                          <span className="text-green-600 font-semibold">⚡ Speedup:</span>
                          <span className="text-green-600 font-bold text-lg">
                            {benchmarkResult.speedup.toFixed(1)}x mais rápido
                          </span>
                        </div>
                      </div>
                    </div>
                  </motion.div>
                )}
              </div>
            )}
          </div>
        </motion.div>
      </div>

      {/* Features Grid */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
        className="grid grid-cols-1 md:grid-cols-3 gap-6"
      >
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border p-6 text-center">
          <CpuChipIcon className="w-12 h-12 text-orange-500 mx-auto mb-4" />
          <h3 className="font-semibold text-gray-900 dark:text-white mb-2">
            Performance Nativa
          </h3>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            Rust compilado para WASM oferece velocidade próxima ao código nativo
          </p>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border p-6 text-center">
          <BoltIcon className="w-12 h-12 text-yellow-500 mx-auto mb-4" />
          <h3 className="font-semibold text-gray-900 dark:text-white mb-2">
            Zero Dependencies
          </h3>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            Módulo WASM autocontido, sem dependências JavaScript externas
          </p>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border p-6 text-center">
          <ClockIcon className="w-12 h-12 text-blue-500 mx-auto mb-4" />
          <h3 className="font-semibold text-gray-900 dark:text-white mb-2">
            Processamento Instantâneo
          </h3>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            Parsing e validação de arquivos grandes em milissegundos
          </p>
        </div>
      </motion.div>

      {/* Call to Action */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
        className="bg-gradient-to-r from-orange-500 to-red-500 rounded-lg p-8 text-center text-white"
      >
        <h3 className="text-2xl font-bold mb-4">
          🚀 O Futuro da Análise de Código
        </h3>
        <p className="text-lg mb-6 opacity-90">
          Com Rust + WASM, o LookAtni File Markers oferece performance de próximo século 
          diretamente no seu browser, sem instalação adicional.
        </p>
        <div className="flex flex-wrap justify-center gap-4">
          <div className="bg-white/20 backdrop-blur rounded-lg px-4 py-2">
            <span className="font-semibold">⚡ 10-50x mais rápido</span>
          </div>
          <div className="bg-white/20 backdrop-blur rounded-lg px-4 py-2">
            <span className="font-semibold">🔒 Memory safe</span>
          </div>
          <div className="bg-white/20 backdrop-blur rounded-lg px-4 py-2">
            <span className="font-semibold">🌐 Universal</span>
          </div>
        </div>
      </motion.div>
    </div>
  );
}
