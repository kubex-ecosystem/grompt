// Extensões de tipos para suportar propriedades específicas do browser

declare global {
  interface HTMLInputElement {
    webkitdirectory?: string;
    directory?: string;
  }

  interface Window {
    showDirectoryPicker?: (options?: {
      mode?: 'read' | 'readwrite';
      startIn?: 'desktop' | 'documents' | 'downloads' | 'music' | 'pictures' | 'videos';
    }) => Promise<FileSystemDirectoryHandle>;

    // WASM bindgen
    wasm_bindgen?: (wasmPath: string) => Promise<void>;
    LookAtniWasmParser?: new () => {
      free(): void;
      parse_ai_code(code: string): any;
      validate_ai_code(code: string): any;
    };
  }

  interface FileSystemDirectoryHandle {
    entries(): AsyncIterableIterator<[string, FileSystemHandle]>;
    getFileHandle(name: string): Promise<FileSystemFileHandle>;
    getDirectoryHandle(name: string): Promise<FileSystemDirectoryHandle>;
  }

  interface FileSystemFileHandle {
    kind: 'file';
    getFile(): Promise<File>;
  }

  interface FileSystemHandle {
    kind: 'file' | 'directory';
    name: string;
  }
}

export { };

