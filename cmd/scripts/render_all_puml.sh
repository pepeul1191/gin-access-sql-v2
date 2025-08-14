#!/bin/bash

# Crear carpeta destino si no existe
mkdir -p docs/pics

# Recorrer todos los .puml en docs/
for file in docs/*.puml; do
  filename=$(basename "$file" .puml)
  plantuml -tpng "$file" -o ./pics
  echo "Generado: docs/pics/${filename}.png"
done