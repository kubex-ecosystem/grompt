#!/usr/bin/env bash
# shellcheck disable=SC2154

# Enhanced TypeScript i18n key extractor
# This script parses nested TypeScript translation objects

echo "📁 Processing TypeScript locales with nested object support..."

# Clear previous results
rm -f i18n_avail_en.txt i18n_avail_ptBR.txt

# Function to extract nested keys from TypeScript file
extract_ts_keys() {
  local file="$1"
  local namespace="$2"
  local lang_file="$3"

  # Use a Python script to parse the TypeScript object structure
  python3 -c "
import re
import sys

def extract_nested_keys(content, prefix=''):
    lines = content.split('\n')
    keys = []
    current_path = []

    for line in lines:
        # Skip comments and imports
        if '//' in line or 'import' in line or 'export' in line:
            continue

        # Find object property definitions
        match = re.search(r'^\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*:', line)
        if match:
            key = match.group(1)
            # Check if it's a string value (leaf node)
            if ':' in line and ('\"' in line or \"'\" in line):
                full_key = f'${prefix}.{key}' if prefix else key
                keys.append(full_key)

    return keys

try:
    with open('$file', 'r') as f:
        content = f.read()

    # Extract all property keys
    keys = []
    lines = content.split('\n')
    current_section = ''

    for line in lines:
        line = line.strip()
        if not line or line.startswith('//') or 'import' in line or 'export' in line:
            continue

        # Look for section headers (objects)
        section_match = re.search(r'^([a-zA-Z_][a-zA-Z0-9_]*)\s*:\s*{', line)
        if section_match:
            current_section = section_match.group(1)
            continue

        # Look for leaf properties
        prop_match = re.search(r'^([a-zA-Z_][a-zA-Z0-9_]*)\s*:', line)
        if prop_match and ('\"' in line or \"'\" in line):
            prop = prop_match.group(1)
            if current_section:
                key = f'${current_section}.{prop}'
            else:
                key = prop
            print(f'$namespace.{key}')

except Exception as e:
    pass
" >> "$lang_file"
}


_process_files() {
  local _what_part_is_scanning="${1:-FRONTEND}"
  local _lang="${2:-en-US}"
  # Process English files
  echo "🇺🇸 Processing English (en-US) files..."
  for f in "${_ROOT_DIR}/${_what_part_is_scanning,,}/locales/$_lang"/*.ts; do
    if [[ -f "$f" && $(basename "$f") != "index.ts" ]]; then
      ns=$(basename "$f" .ts)
      echo "  📄 Processing $ns..."
      extract_ts_keys "$f" "$ns" "REPORT_${_what_part_is_scanning}-i18n_avail_${_lang}.txt"
    fi
  done
}

# Process Portuguese files
echo "🇧🇷 Processing Portuguese (pt-BR) files..."
for f in frontend/locales/pt-BR/*.ts; do
  if [[ -f "$f" && $(basename "$f") != "index.ts" ]]; then
    ns=$(basename "$f" .ts)
    echo "  📄 Processing $ns..."
    extract_ts_keys "$f" "$ns" "i18n_avail_ptBR.txt"
  fi
done

# Sort and deduplicate
if [[ -f i18n_avail_en.txt ]]; then
  sort -u -o i18n_avail_en.txt i18n_avail_en.txt
  echo "✅ EN keys found: $(wc -l < i18n_avail_en.txt)"
else
  touch i18n_avail_en.txt
  echo "⚠️  No EN keys found"
fi

if [[ -f i18n_avail_ptBR.txt ]]; then
  sort -u -o i18n_avail_ptBR.txt i18n_avail_ptBR.txt
  echo "✅ PT-BR keys found: $(wc -l < i18n_avail_ptBR.txt)"
else
  touch i18n_avail_ptBR.txt
  echo "⚠️  No PT-BR keys found"
fi

echo "🎯 Ready for debt analysis!"
