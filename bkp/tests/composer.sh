#!/bin/bash

# Define o diretório base
BASE_DIR="lib"

# Lista de arquivos a serem criados
FILES=(
  "config.sh"
  "utils.sh"
  "platform.sh"
  "build.sh"
  "validate.sh"
  "install_funcs.sh"
  "info.sh"
)

# Cria o diretório base, se ainda não existir
mkdir -p "$BASE_DIR"

# Cria os arquivos dentro do diretório
for file in "${FILES[@]}"; do
  FILE_PATH="$BASE_DIR/$file"
  if [[ ! -f "$FILE_PATH" ]]; then
    touch "$FILE_PATH"
    printf '%s' "#!/bin/bash" | tee "$FILE_PATH" >/dev/null
    printf '%s' "# $file - script placeholder" | tee -a "$FILE_PATH" >/dev/null
    chmod +x "$FILE_PATH"
    echo "Criado: $FILE_PATH"
  else
    echo "Já existe: $FILE_PATH"
  fi
done

