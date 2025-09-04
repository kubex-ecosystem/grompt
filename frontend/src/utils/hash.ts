export function hashString(input: string): string {
  // Simple DJB2 hash, hex output
  let hash = 5381;
  for (let i = 0; i < input.length; i++) {
    hash = ((hash << 5) + hash) + input.charCodeAt(i);
    hash |= 0; // force 32-bit
  }
  // Convert to unsigned and hex
  const unsigned = hash >>> 0;
  return unsigned.toString(16);
}

