/* tslint:disable */
/* eslint-disable */
export function main(): void;
export class LookAtniWasmParser {
  free(): void;
  constructor();
  /**
   * Parse código AI com marcadores invisíveis - SUPER PERFORMANCE!
   */
  parse_ai_code(code: string): any;
  /**
   * Validar código AI - SUPER INTELIGENTE!
   */
  validate_ai_code(code: string): any;
  /**
   * Extrair estatísticas avançadas - ANALYTICS POWER!
   */
  get_advanced_stats(code: string): any;
  /**
   * Benchmark performance contra JavaScript
   */
  benchmark_vs_js(code: string, iterations: number): any;
}

export type InitInput = RequestInfo | URL | Response | BufferSource | WebAssembly.Module;

export interface InitOutput {
  readonly memory: WebAssembly.Memory;
  readonly __wbg_lookatniwasmparser_free: (a: number, b: number) => void;
  readonly lookatniwasmparser_new: () => number;
  readonly lookatniwasmparser_parse_ai_code: (a: number, b: number, c: number) => any;
  readonly lookatniwasmparser_validate_ai_code: (a: number, b: number, c: number) => any;
  readonly lookatniwasmparser_get_advanced_stats: (a: number, b: number, c: number) => any;
  readonly lookatniwasmparser_benchmark_vs_js: (a: number, b: number, c: number, d: number) => any;
  readonly main: () => void;
  readonly __wbindgen_exn_store: (a: number) => void;
  readonly __externref_table_alloc: () => number;
  readonly __wbindgen_export_2: WebAssembly.Table;
  readonly __wbindgen_malloc: (a: number, b: number) => number;
  readonly __wbindgen_realloc: (a: number, b: number, c: number, d: number) => number;
  readonly __wbindgen_start: () => void;
}

export type SyncInitInput = BufferSource | WebAssembly.Module;
/**
* Instantiates the given `module`, which can either be bytes or
* a precompiled `WebAssembly.Module`.
*
* @param {{ module: SyncInitInput }} module - Passing `SyncInitInput` directly is deprecated.
*
* @returns {InitOutput}
*/
export function initSync(module: { module: SyncInitInput } | SyncInitInput): InitOutput;

/**
* If `module_or_path` is {RequestInfo} or {URL}, makes a request and
* for everything else, calls `WebAssembly.instantiate` directly.
*
* @param {{ module_or_path: InitInput | Promise<InitInput> }} module_or_path - Passing `InitInput` directly is deprecated.
*
* @returns {Promise<InitOutput>}
*/
export default function __wbg_init (module_or_path?: { module_or_path: InitInput | Promise<InitInput> } | InitInput | Promise<InitInput>): Promise<InitOutput>;
