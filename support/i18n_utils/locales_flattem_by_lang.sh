#!/usr/bin/env bash

# Clear previous results
rm -f i18n_avail_en.txt i18n_avail_ptBR.txt

# EN-US (TypeScript files, not JSON)
echo "📁 Processing English TypeScript locales..."
for f in frontend/locales/en-US/*.ts; do
  if [[ -f "$f" ]]; then
    ns=$(basename "$f" .ts)
    # Extract object keys from TypeScript export default
    grep -o "'[^']*':" "$f" | sed "s/'//g; s/://; s/^/$ns./" >> i18n_avail_en.txt
  fi
done
if [[ -f i18n_avail_en.txt ]]; then
  sort -u -o i18n_avail_en.txt i18n_avail_en.txt
  echo "✅ EN keys: $(wc -l < i18n_avail_en.txt)"
fi

# PT-BR (TypeScript files, not JSON)
echo "📁 Processing Portuguese TypeScript locales..."
for f in frontend/locales/pt-BR/*.ts; do
  if [[ -f "$f" ]]; then
    ns=$(basename "$f" .ts)
    # Extract object keys from TypeScript export default
    grep -o "'[^']*':" "$f" | sed "s/'//g; s/://; s/^/$ns./" >> i18n_avail_ptBR.txt
  fi
done
if [[ -f i18n_avail_ptBR.txt ]]; then
  sort -u -o i18n_avail_ptBR.txt i18n_avail_ptBR.txt
  echo "✅ PT-BR keys: $(wc -l < i18n_avail_ptBR.txt)"
fi
