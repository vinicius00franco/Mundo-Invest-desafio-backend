#!/bin/bash

# Script para testar cenários de webhook - Processamento de Eventos
# Baseado em: bdd/processamento-eventos/webhook.md

# Configurações
BASE_URL="http://localhost:8080"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT_DIR="${SCRIPT_DIR}/../../exports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_FILE="${OUTPUT_DIR}/cenarios_webhook_${TIMESTAMP}.txt"

# Criar diretório de exportação se não existir
mkdir -p "$OUTPUT_DIR"

# Inicializar arquivo de saída
echo "========================================" > "$OUTPUT_FILE"
echo "TESTES DE CENÁRIOS DE WEBHOOK" >> "$OUTPUT_FILE"
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
  
  echo "Requisição: POST ${BASE_URL}/webhooks/pipefy/card-updated" >> "$OUTPUT_FILE"
  echo "" >> "$OUTPUT_FILE"
  
  # Executar requisição curl e capturar resposta
  RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "${BASE_URL}/webhooks/pipefy/card-updated" \
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

# Cenário 1: Processar webhook com cliente de prioridade alta
PAYLOAD1='{
  "event_id": "evt_prioridade_alta_001",
  "card_id": "card_123",
  "cliente_email": "maria.santos@example.com",
  "timestamp": "2026-05-27T10:00:00Z"
}'
executar_teste "CENÁRIO 1: Processar webhook com cliente de prioridade alta" \
  "Processar webhook para cliente com patrimônio >= 200.000" \
  "$PAYLOAD1" \
  "200"

# Cenário 2: Processar webhook com cliente de prioridade normal
PAYLOAD2='{
  "event_id": "evt_prioridade_normal_001",
  "card_id": "card_456",
  "cliente_email": "joao.silva@example.com",
  "timestamp": "2026-05-27T10:00:00Z"
}'
executar_teste "CENÁRIO 2: Processar webhook com cliente de prioridade normal" \
  "Processar webhook para cliente com patrimônio < 200.000" \
  "$PAYLOAD2" \
  "200"

# Cenário 3: Processar webhook com cliente não encontrado
PAYLOAD3='{
  "event_id": "evt_cliente_inexistente",
  "card_id": "card_789",
  "cliente_email": "nao.existe@example.com",
  "timestamp": "2026-05-27T10:00:00Z"
}'
executar_teste "CENÁRIO 3: Processar webhook com cliente não encontrado" \
  "Tentar processar webhook para cliente que não existe no banco" \
  "$PAYLOAD3" \
  "500"

# Cenário 4: Processar webhook duplicado (idempotência)
PAYLOAD4='{
  "event_id": "evt_prioridade_alta_001",
  "card_id": "card_123",
  "cliente_email": "maria.santos@example.com",
  "timestamp": "2026-05-27T10:00:00Z"
}'
executar_teste "CENÁRIO 4: Processar webhook duplicado (idempotência)" \
  "Tentar processar o mesmo evento novamente para testar idempotência" \
  "$PAYLOAD4" \
  "200"

# Cenário 5: Processar webhook sem event_id
PAYLOAD5='{
  "card_id": "card_999",
  "cliente_email": "teste@example.com",
  "timestamp": "2026-05-27T10:00:00Z"
}'
executar_teste "CENÁRIO 5: Processar webhook sem event_id" \
  "Tentar processar webhook sem fornecer event_id" \
  "$PAYLOAD5" \
  "400"

# Cenário 6: Processar webhook sem card_id
PAYLOAD6='{
  "event_id": "evt_sem_card",
  "cliente_email": "teste@example.com",
  "timestamp": "2026-05-27T10:00:00Z"
}'
executar_teste "CENÁRIO 6: Processar webhook sem card_id" \
  "Tentar processar webhook sem fornecer card_id" \
  "$PAYLOAD6" \
  "400"

# Cenário 7: Processar webhook sem cliente_email
PAYLOAD7='{
  "event_id": "evt_sem_email",
  "card_id": "card_999",
  "timestamp": "2026-05-27T10:00:00Z"
}'
executar_teste "CENÁRIO 7: Processar webhook sem cliente_email" \
  "Tentar processar webhook sem fornecer cliente_email" \
  "$PAYLOAD7" \
  "400"

# Cenário 8: Processar webhook sem timestamp
PAYLOAD8='{
  "event_id": "evt_sem_timestamp",
  "card_id": "card_999",
  "cliente_email": "teste@example.com"
}'
executar_teste "CENÁRIO 8: Processar webhook sem timestamp" \
  "Tentar processar webhook sem fornecer timestamp" \
  "$PAYLOAD8" \
  "400"

echo "========================================" >> "$OUTPUT_FILE"
echo "FIM DOS TESTES DE WEBHOOK" >> "$OUTPUT_FILE"
echo "Resultado salvo em: $OUTPUT_FILE" >> "$OUTPUT_FILE"
echo "========================================" >> "$OUTPUT_FILE"

# Exibir resultado no terminal
cat "$OUTPUT_FILE"
