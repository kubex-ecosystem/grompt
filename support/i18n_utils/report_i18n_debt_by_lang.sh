#!/usr/bin/env bash

echo "📊 Generating i18n technical debt report..."

# Ensure input files exist
if [[ ! -f i18n_used_keys.txt ]]; then
  echo "❌ i18n_used_keys.txt not found! Run get_i18n_used_keys.sh first"
  exit 1
fi

if [[ ! -f i18n_avail_en.txt ]]; then
  echo "❌ i18n_avail_en.txt not found! Run locales_flattem_by_lang.sh first"
  exit 1
fi

# English analysis
echo "🔍 Analyzing English translations..."
comm -23 i18n_used_keys.txt i18n_avail_en.txt > i18n_missing_en.txt
comm -13 i18n_used_keys.txt i18n_avail_en.txt > i18n_unused_en.txt

# Portuguese analysis (if available)
if [[ -f i18n_avail_ptBR.txt ]]; then
  echo "🔍 Analyzing Portuguese translations..."
  comm -23 i18n_used_keys.txt i18n_avail_ptBR.txt > i18n_missing_ptBR.txt
  comm -13 i18n_used_keys.txt i18n_avail_ptBR.txt > i18n_unused_ptBR.txt
fi

# Report results
echo ""
echo "📋 I18N DEBT REPORT"
echo "==================="
echo "🔑 Used keys: $(wc -l < i18n_used_keys.txt)"
echo ""
echo "🇺🇸 ENGLISH:"
echo "  ❌ Missing: $(wc -l < i18n_missing_en.txt)"
echo "  🗑️  Unused:  $(wc -l < i18n_unused_en.txt)"

if [[ -f i18n_missing_ptBR.txt ]]; then
  echo ""
  echo "🇧🇷 PORTUGUÊS:"
  echo "  ❌ Missing: $(wc -l < i18n_missing_ptBR.txt)"
  echo "  🗑️  Unused:  $(wc -l < i18n_unused_ptBR.txt)"
fi

echo ""
echo "📁 Generated files:"
echo "  - i18n_missing_en.txt"
echo "  - i18n_unused_en.txt"
if [[ -f i18n_missing_ptBR.txt ]]; then
  echo "  - i18n_missing_ptBR.txt"
  echo "  - i18n_unused_ptBR.txt"
fi

echo ""
echo "🎯 Next steps:"
echo "  1. Review missing keys and add translations"
echo "  2. Remove unused keys to reduce bundle size"env bash

# faltando no EN
comm -23 i18n_used_keys.txt i18n_avail_en.txt   > i18n_missing_en.txt
# chaves “sobrando” no EN (não usadas)
comm -13 i18n_used_keys.txt i18n_avail_en.txt   > i18n_unused_en.txt

# faltando no PT-BR
comm -23 i18n_used_keys.txt i18n_avail_ptBR.txt > i18n_missing_ptBR.txt
# sobras no PT-BR
comm -13 i18n_used_keys.txt i18n_avail_ptBR.txt > i18n_unused_ptBR.txt

