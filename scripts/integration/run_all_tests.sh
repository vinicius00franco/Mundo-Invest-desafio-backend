#!/bin/bash

# Script principal para executar todos os testes de integração
# Organiza e executa os testes de sucesso, falha e validação

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT_DIR="${SCRIPT_DIR}/../exports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
SUMMARY_FILE="${OUTPUT_DIR}/test_summary_${TIMESTAMP}.txt"

# Criar diretório de exportação se não existir
mkdir -p "$OUTPUT_DIR"

# Inicializar arquivo de resumo
echo "========================================" > "$SUMMARY_FILE"
echo "RESUMO DE TESTES DE INTEGRAÇÃO" >> "$SUMMARY_FILE"
echo "Data: $(date)" >> "$SUMMARY_FILE"
echo "========================================" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Iniciando execução de todos os testes de integração...${NC}"
echo ""

# Executar testes de sucesso
echo -e "${YELLOW}=== EXECUTANDO TESTES DE SUCESSO ===${NC}"
if bash "${SCRIPT_DIR}/sucesso/testar_cenarios_sucesso.sh"; then
    echo -e "${GREEN}✓ Testes de sucesso concluídos${NC}"
    echo "✓ Testes de sucesso: PASSOU" >> "$SUMMARY_FILE"
else
    echo -e "${RED}✗ Testes de sucesso falharam${NC}"
    echo "✗ Testes de sucesso: FALHOU" >> "$SUMMARY_FILE"
fi
echo "" >> "$SUMMARY_FILE"
echo ""

# Executar testes de falha
echo -e "${YELLOW}=== EXECUTANDO TESTES DE FALHA ===${NC}"
if bash "${SCRIPT_DIR}/falha/testar_cenarios_falha.sh"; then
    echo -e "${GREEN}✓ Testes de falha concluídos${NC}"
    echo "✓ Testes de falha: PASSOU" >> "$SUMMARY_FILE"
else
    echo -e "${RED}✗ Testes de falha falharam${NC}"
    echo "✗ Testes de falha: FALHOU" >> "$SUMMARY_FILE"
fi
echo "" >> "$SUMMARY_FILE"
echo ""

# Executar testes de validação
echo -e "${YELLOW}=== EXECUTANDO TESTES DE VALIDAÇÃO ===${NC}"
if bash "${SCRIPT_DIR}/validacao/testar_cenarios_validacao.sh"; then
    echo -e "${GREEN}✓ Testes de validação concluídos${NC}"
    echo "✓ Testes de validação: PASSOU" >> "$SUMMARY_FILE"
else
    echo -e "${RED}✗ Testes de validação falharam${NC}"
    echo "✗ Testes de validação: FALHOU" >> "$SUMMARY_FILE"
fi
echo "" >> "$SUMMARY_FILE"
echo ""

echo "========================================" >> "$SUMMARY_FILE"
echo "FIM DOS TESTES" >> "$SUMMARY_FILE"
echo "Resumo salvo em: $SUMMARY_FILE" >> "$SUMMARY_FILE"
echo "========================================" >> "$SUMMARY_FILE"

echo -e "${YELLOW}=== RESUMO ===${NC}"
echo "Todos os resultados foram salvos em: $OUTPUT_DIR"
echo "Resumo geral: $SUMMARY_FILE"
cat "$SUMMARY_FILE"
