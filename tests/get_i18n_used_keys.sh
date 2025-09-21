#!/usr/bin/env bash

# Scan for i18n usage in frontend (no src folder)
rg -no --pcre2 "t\(\s*['\"\$(]([A-Za-z][\w-]+)\.([A-Za-z0-9_.-]+)['\")]\s*(?:,|\))" ./frontend \
| awk -F: '{print $3}' \
| sed -E "s/^t\(['\"\`]//; s/['\"\`].*$//" \
| sort -u > i18n_used_keys.txt
