#!/bin/bash

# Script para testar cenários de validação - Criação de Cliente
# Baseado em: bdd/criacao-cliente/validacao.md

# Configurações
BASE_URL="http://localhost:8080"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT_DIR="${SCRIPT_DIR}/../../exports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_FILE="${OUTPUT_DIR}/cenarios_validacao_${TIMESTAMP}.txt"

# Criar diretório de exportação se não existir
mkdir -p "$OUTPUT_DIR"

# Inicializar arquivo de saída
echo "========================================" > "$OUTPUT_FILE"
echo "TESTES DE CENÁRIOS DE VALIDAÇÃO" >> "$OUTPUT_FILE"
echo "Data: $(date)" >> "$OUTPUT_FILE"
echo "========================================" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Função para executar teste
executar_teste() {
  local cenario=$1
  local descricao=$2
  local payload=$3
  local expected_code=$4
  
  echo "=== $cenario ===" >> "$OUTPUT_FILE"
  echo "Descrição: $descricao" >> "$OUTPUT_FILE"
  echo "" >> "$OUTPUT_FILE"
  
  echo "Payload enviado:" >> "$OUTPUT_FILE"
  if echo "$payload" | jq . > /dev/null 2>&1; then
    echo "$payload" | jq . >> "$OUTPUT_FILE"
  else
    echo "$payload" >> "$OUTPUT_FILE"
  fi
  echo "" >> "$OUTPUT_FILE"
  
  echo "Requisição: POST ${BASE_URL}/clientes" >> "$OUTPUT_FILE"
  echo "" >> "$OUTPUT_FILE"
  
  # Executar requisição curl e capturar resposta
  RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "${BASE_URL}/clientes" \
    -H "Content-Type: application/json" \
    -d "$payload")
  
  # Extrair código HTTP e corpo da resposta
  HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
  BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE:/d')
  
  echo "Código HTTP: $HTTP_CODE" >> "$OUTPUT_FILE"
  echo "Resposta JSON:" >> "$OUTPUT_FILE"
  if echo "$BODY" | jq . > /dev/null 2>&1; then
    echo "$BODY" | jq . >> "$OUTPUT_FILE"
  else
    echo "$BODY" >> "$OUTPUT_FILE"
  fi
  echo "" >> "$OUTPUT_FILE"
  
  # Verificações
  echo "Verificações:" >> "$OUTPUT_FILE"
  if [ "$HTTP_CODE" = "$expected_code" ]; then
    echo "✓ HTTP $expected_code - PASSOU" >> "$OUTPUT_FILE"
  else
    echo "✗ HTTP $HTTP_CODE - FALHOU (esperado $expected_code)" >> "$OUTPUT_FILE"
  fi
  echo "" >> "$OUTPUT_FILE"
}

# Cenário 1: Criar cliente com e-mail inválido
PAYLOAD1='{
  "nome": "João Silva",
  "email": "email-invalido",
  "tipoSolicitacao": "Atualização cadastral",
  "valorPatrimonio": 250000
}'
executar_teste "CENÁRIO 1: Criar cliente com e-mail inválido" \
  "Tentar criar cliente com e-mail em formato inválido" \
  "$PAYLOAD1" \
  "400"

# Cenário 2: Criar cliente com patrimônio negativo
UNIQUE_EMAIL2="joao.silva2.${TIMESTAMP}@example.com"
PAYLOAD2="{
  \"nome\": \"João Silva\",
  \"email\": \"${UNIQUE_EMAIL2}\",
  \"tipoSolicitacao\": \"Atualização cadastral\",
  \"valorPatrimonio\": -1000
}"
executar_teste "CENÁRIO 2: Criar cliente com patrimônio negativo" \
  "Tentar criar cliente com valor de patrimônio negativo" \
  "$PAYLOAD2" \
  "400"

# Cenário 3: Criar cliente com patrimônio zero
UNIQUE_EMAIL3="joao.silva3.${TIMESTAMP}@example.com"
PAYLOAD3="{
  \"nome\": \"João Silva\",
  \"email\": \"${UNIQUE_EMAIL3}\",
  \"tipoSolicitacao\": \"Atualização cadastral\",
  \"valorPatrimonio\": 0
}"
executar_teste "CENÁRIO 3: Criar cliente com patrimônio zero" \
  "Tentar criar cliente com valor de patrimônio igual a zero" \
  "$PAYLOAD3" \
  "400"

# Cenário 4: Criar cliente com nome vazio
UNIQUE_EMAIL4="joao.silva4.${TIMESTAMP}@example.com"
PAYLOAD4="{
  \"nome\": \"\",
  \"email\": \"${UNIQUE_EMAIL4}\",
  \"tipoSolicitacao\": \"Atualização cadastral\",
  \"valorPatrimonio\": 250000
}"
executar_teste "CENÁRIO 4: Criar cliente com nome vazio" \
  "Tentar criar cliente com nome em branco" \
  "$PAYLOAD4" \
  "400"

# Cenário 5: Criar cliente com e-mail vazio
PAYLOAD5='{
  "nome": "João Silva",
  "email": "",
  "tipoSolicitacao": "Atualização cadastral",
  "valorPatrimonio": 250000
}'
executar_teste "CENÁRIO 5: Criar cliente com e-mail vazio" \
  "Tentar criar cliente com e-mail em branco" \
  "$PAYLOAD5" \
  "400"

# Cenário 6: Criar cliente com tipoSolicitacao vazio
UNIQUE_EMAIL6="joao.silva6.${TIMESTAMP}@example.com"
PAYLOAD6="{
  \"nome\": \"João Silva\",
  \"email\": \"${UNIQUE_EMAIL6}\",
  \"tipoSolicitacao\": \"\",
  \"valorPatrimonio\": 250000
}"
executar_teste "CENÁRIO 6: Criar cliente com tipoSolicitacao vazio" \
  "Tentar criar cliente com tipo de solicitação em branco" \
  "$PAYLOAD6" \
  "400"

# Cenário 7: Criar cliente com nome muito curto
UNIQUE_EMAIL7="joao.silva7.${TIMESTAMP}@example.com"
PAYLOAD7="{
  \"nome\": \"AB\",
  \"email\": \"${UNIQUE_EMAIL7}\",
  \"tipoSolicitacao\": \"Atualização cadastral\",
  \"valorPatrimonio\": 250000
}"
executar_teste "CENÁRIO 7: Criar cliente com nome muito curto" \
  "Tentar criar cliente com nome com menos de 3 caracteres" \
  "$PAYLOAD7" \
  "400"

echo "========================================" >> "$OUTPUT_FILE"
echo "FIM DOS TESTES DE VALIDAÇÃO" >> "$OUTPUT_FILE"
echo "Resultado salvo em: $OUTPUT_FILE" >> "$OUTPUT_FILE"
echo "========================================" >> "$OUTPUT_FILE"

# Exibir resultado no terminal
cat "$OUTPUT_FILE"
