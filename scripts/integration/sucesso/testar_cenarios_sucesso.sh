#!/bin/bash

# Script para testar cenários de sucesso - Criação de Cliente
# Baseado em: bdd/criacao-cliente/sucesso.md

# Configurações
BASE_URL="http://localhost:8080"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT_DIR="${SCRIPT_DIR}/../../exports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_FILE="${OUTPUT_DIR}/cenarios_sucesso_${TIMESTAMP}.txt"

# Criar diretório de exportação se não existir
mkdir -p "$OUTPUT_DIR"

# Inicializar arquivo de saída
echo "========================================" > "$OUTPUT_FILE"
echo "TESTES DE CENÁRIOS DE SUCESSO" >> "$OUTPUT_FILE"
echo "Data: $(date)" >> "$OUTPUT_FILE"
echo "========================================" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Cenário 1: Criar cliente com dados válidos
echo "=== CENÁRIO 1: Criar cliente com dados válidos ===" >> "$OUTPUT_FILE"
echo "Descrição: Criar um novo cliente com todos os campos obrigatórios preenchidos corretamente" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Gerar email único para evitar conflitos
UNIQUE_EMAIL="joao.silva.${TIMESTAMP}@example.com"
PAYLOAD="{
  \"cliente_nome\": \"João Silva\",
  \"cliente_email\": \"${UNIQUE_EMAIL}\",
  \"tipo_solicitacao\": \"Atualização cadastral\",
  \"valor_patrimonio\": 250000
}"

echo "Payload enviado:" >> "$OUTPUT_FILE"
echo "$PAYLOAD" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "Requisição: POST ${BASE_URL}/clientes" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Executar requisição curl e capturar resposta
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "${BASE_URL}/clientes" \
  -H "Content-Type: application/json" \
  -d "$PAYLOAD")

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
if [ "$HTTP_CODE" = "201" ]; then
  echo "✓ HTTP 201 Created - PASSOU" >> "$OUTPUT_FILE"
else
  echo "✗ HTTP $HTTP_CODE - FALHOU (esperado 201)" >> "$OUTPUT_FILE"
fi

# Verificar se contém campos esperados
if echo "$BODY" | jq -e '.identificador_interno' > /dev/null 2>&1; then
  echo "✓ Campo identificador_interno presente - PASSOU" >> "$OUTPUT_FILE"
else
  echo "✗ Campo identificador_interno ausente - FALHOU" >> "$OUTPUT_FILE"
fi

if echo "$BODY" | jq -e '.identificador_externo' > /dev/null 2>&1; then
  echo "✓ Campo identificador_externo (card_id) presente - PASSOU" >> "$OUTPUT_FILE"
else
  echo "✗ Campo identificador_externo ausente - FALHOU" >> "$OUTPUT_FILE"
fi

if echo "$BODY" | jq -e '.status' > /dev/null 2>&1; then
  STATUS=$(echo "$BODY" | jq -r '.status')
  echo "✓ Campo status presente: $STATUS" >> "$OUTPUT_FILE"
  if [ "$STATUS" = "Aguardando Análise" ]; then
    echo "✓ Status é 'Aguardando Análise' - PASSOU" >> "$OUTPUT_FILE"
  else
    echo "✗ Status não é 'Aguardando Análise' - FALHOU" >> "$OUTPUT_FILE"
  fi
else
  echo "✗ Campo status ausente - FALHOU" >> "$OUTPUT_FILE"
fi

echo "" >> "$OUTPUT_FILE"
echo "========================================" >> "$OUTPUT_FILE"
echo "FIM DOS TESTES DE SUCESSO" >> "$OUTPUT_FILE"
echo "Resultado salvo em: $OUTPUT_FILE" >> "$OUTPUT_FILE"
echo "========================================" >> "$OUTPUT_FILE"

# Exibir resultado no terminal
cat "$OUTPUT_FILE"
