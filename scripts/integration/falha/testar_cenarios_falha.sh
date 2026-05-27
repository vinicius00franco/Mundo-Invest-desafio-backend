#!/bin/bash

# Script para testar cenários de falha - Criação de Cliente
# Baseado em: bdd/criacao-cliente/falha.md

# Configurações
BASE_URL="http://localhost:8080"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT_DIR="${SCRIPT_DIR}/../../exports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_FILE="${OUTPUT_DIR}/cenarios_falha_${TIMESTAMP}.txt"

# Criar diretório de exportação se não existir
mkdir -p "$OUTPUT_DIR"

# Inicializar arquivo de saída
echo "========================================" > "$OUTPUT_FILE"
echo "TESTES DE CENÁRIOS DE FALHA" >> "$OUTPUT_FILE"
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

# Cenário 1: Criar cliente sem campo cliente_nome
UNIQUE_EMAIL="joao.silha.${TIMESTAMP}@example.com"
PAYLOAD1="{
  \"email\": \"${UNIQUE_EMAIL}\",
  \"tipoSolicitacao\": \"Atualização cadastral\",
  \"valorPatrimonio\": 250000
}"
executar_teste "CENÁRIO 1: Criar cliente sem campo cliente_nome" \
  "Tentar criar cliente sem fornecer o campo cliente_nome" \
  "$PAYLOAD1" \
  "400"

# Cenário 2: Criar cliente sem campo tipoSolicitacao
UNIQUE_EMAIL2="joao.silva2.${TIMESTAMP}@example.com"
PAYLOAD2="{
  \"nome\": \"João Silva\",
  \"email\": \"${UNIQUE_EMAIL2}\",
  \"valorPatrimonio\": 250000
}"
executar_teste "CENÁRIO 2: Criar cliente sem campo tipoSolicitacao" \
  "Tentar criar cliente sem fornecer o campo tipoSolicitacao" \
  "$PAYLOAD2" \
  "400"

# Cenário 3: Criar cliente sem campo valorPatrimonio
UNIQUE_EMAIL3="joao.silva3.${TIMESTAMP}@example.com"
PAYLOAD3="{
  \"nome\": \"João Silva\",
  \"email\": \"${UNIQUE_EMAIL3}\",
  \"tipoSolicitacao\": \"Atualização cadastral\"
}"
executar_teste "CENÁRIO 3: Criar cliente sem campo valorPatrimonio" \
  "Tentar criar cliente sem fornecer o campo valorPatrimonio" \
  "$PAYLOAD3" \
  "400"

# Cenário 4: Criar cliente sem campo email
PAYLOAD4='{
  "nome": "João Silva",
  "tipoSolicitacao": "Atualização cadastral",
  "valorPatrimonio": 250000
}'
executar_teste "CENÁRIO 4: Criar cliente sem campo email" \
  "Tentar criar cliente sem fornecer o campo email" \
  "$PAYLOAD4" \
  "400"

# Cenário 5: Criar cliente com erro de banco de dados
# Nota: Este cenário requer que o banco de dados esteja indisponível
# Para testar, pare o container do PostgreSQL antes de executar este teste
UNIQUE_EMAIL5="joao.silva5.${TIMESTAMP}@example.com"
PAYLOAD5="{
  \"nome\": \"João Silva\",
  \"email\": \"${UNIQUE_EMAIL5}\",
  \"tipoSolicitacao\": \"Atualização cadastral\",
  \"valorPatrimonio\": 250000
}"
echo "=== CENÁRIO 5: Criar cliente com erro de banco de dados ===" >> "$OUTPUT_FILE"
echo "Descrição: Tentar criar cliente quando banco de dados está indisponível" >> "$OUTPUT_FILE"
echo "Nota: Este cenário requer que o banco de dados esteja parado" >> "$OUTPUT_FILE"
echo "Para testar: docker-compose stop postgres" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "Payload enviado:" >> "$OUTPUT_FILE"
if echo "$PAYLOAD5" | jq . > /dev/null 2>&1; then
  echo "$PAYLOAD5" | jq . >> "$OUTPUT_FILE"
else
  echo "$PAYLOAD5" >> "$OUTPUT_FILE"
fi
echo "" >> "$OUTPUT_FILE"

echo "Requisição: POST ${BASE_URL}/clientes" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "${BASE_URL}/clientes" \
  -H "Content-Type: application/json" \
  -d "$PAYLOAD5")

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

echo "Verificações:" >> "$OUTPUT_FILE"
if [ "$HTTP_CODE" = "500" ]; then
  echo "✓ HTTP 500 Internal Server Error - PASSOU" >> "$OUTPUT_FILE"
else
  echo "✗ HTTP $HTTP_CODE - FALHOU (esperado 500)" >> "$OUTPUT_FILE"
  echo "  Nota: Se o banco está rodando, este teste falhará corretamente" >> "$OUTPUT_FILE"
fi
echo "" >> "$OUTPUT_FILE"

echo "========================================" >> "$OUTPUT_FILE"
echo "FIM DOS TESTES DE FALHA" >> "$OUTPUT_FILE"
echo "Resultado salvo em: $OUTPUT_FILE" >> "$OUTPUT_FILE"
echo "========================================" >> "$OUTPUT_FILE"

# Exibir resultado no terminal
cat "$OUTPUT_FILE"
