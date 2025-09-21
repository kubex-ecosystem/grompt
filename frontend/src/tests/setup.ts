/**
 * Test setup configuration for Grompt frontend tests
 */

import '@testing-library/jest-dom';

// Mock window objects that might not be available in test environment
Object.defineProperty(window, 'navigator', {
  value: {
    clipboard: {
      writeText: jest.fn(() => Promise.resolve())
    }
  },
  writable: true
});

// Mock console methods to reduce noise in tests
global.console = {
  ...console,
  // Uncomment to suppress console output in tests
  // log: jest.fn(),
  // warn: jest.fn(),
  // error: jest.fn(),
};

// Mock EventSource for streaming tests
class MockEventSource {
  url: string;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  readyState: number = 1; // OPEN

  constructor(url: string) {
    this.url = url;
  }

  close() {
    this.readyState = 2; // CLOSED
  }

  // Helper method for tests to trigger events
  triggerMessage(data: any) {
    if (this.onmessage) {
      this.onmessage(new MessageEvent('message', { data: JSON.stringify(data) }));
    }
  }

  triggerError() {
    if (this.onerror) {
      this.onerror(new Event('error'));
    }
  }
}

global.EventSource = MockEventSource as any;

// Mock fetch if not already mocked
if (!global.fetch) {
  global.fetch = jest.fn();
}

// Clean up after each test
afterEach(() => {
  jest.clearAllMocks();
});