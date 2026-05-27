# Testes de Integração - Mundo Invest

Este diretório contém os scripts de teste de integração para a API do Mundo Invest, organizados por categoria de teste.

## Estrutura de Diretórios

```
scripts/integration/
├── run_all_tests.sh          # Script principal para executar todos os testes
├── sucesso/                  # Testes de cenários de sucesso
│   └── testar_cenarios_sucesso.sh
├── falha/                    # Testes de cenários de falha
│   └── testar_cenarios_falha.sh
├── validacao/                # Testes de cenários de validação
│   └── testar_cenarios_validacao.sh
└── exports/                  # Diretório para armazenar os resultados dos testes
```

## Pré-requisitos

1. Servidor da API rodando na porta 8080
2. Banco de dados PostgreSQL configurado e acessível
3. Utilitários: `curl`, `jq`

## Como Executar

### Executar todos os testes
```bash
./scripts/integration/run_all_tests.sh
```

### Executar testes específicos
```bash
# Testes de sucesso
./scripts/integration/sucesso/testar_cenarios_sucesso.sh

# Testes de falha
./scripts/integration/falha/testar_cenarios_falha.sh

# Testes de validação
./scripts/integration/validacao/testar_cenarios_validacao.sh
```

## Categorias de Testes

### Testes de Sucesso
Verificam se a API funciona corretamente com dados válidos:
- Criação de cliente com todos os campos obrigatórios
- Validação dos campos retornados na resposta
- Verificação do status inicial do cliente

### Testes de Falha
Verificam se a API trata corretamente erros esperados:
- Campos obrigatórios ausentes
- Erros de conexão com banco de dados
- Tratamento de erros internos

### Testes de Validação
Verificam se a API valida corretamente os dados de entrada:
- Formato de e-mail inválido
- Valores negativos ou zero para patrimônio
- Campos vazios
- Valores fora dos limites permitidos

## Padrão de Nomenclatura JSON

A API utiliza o padrão **camelCase** para todos os campos JSON, seguindo as melhores práticas de APIs REST modernas:

**Request (Campos em camelCase):**
```json
{
  "nome": "João Silva",
  "email": "joao.silva@example.com",
  "tipoSolicitacao": "Atualização cadastral",
  "valorPatrimonio": 250000
}
```

**Response (Campos em camelCase):**
```json
{
  "mensagem": "Cliente criado com sucesso",
  "identificadorInterno": 1,
  "identificadorExterno": "card_123",
  "nome": "João Silva",
  "email": "joao.silva@example.com",
  "valorPatrimonio": 250000,
  "tipoSolicitacao": "Atualização cadastral",
  "status": "Aguardando Análise",
  "dataCriacao": "2026-05-27T10:00:00Z"
}
```

## Resultados

Todos os resultados são salvos no diretório `scripts/exports/` com timestamp:
- `cenarios_sucesso_<timestamp>.txt`
- `cenarios_falha_<timestamp>.txt`
- `cenarios_validacao_<timestamp>.txt`
- `test_summary_<timestamp>.txt` (resumo geral)

## Análise de Resultados

Como Senior Backend, verifique:
1. **Códigos HTTP**: Estejam conforme o esperado (201 para sucesso, 400 para validação, 500 para erros)
2. **Respostas JSON**: Estrutura correta e campos esperados presentes
3. **Validações**: Regras de negócio aplicadas corretamente
4. **Tratamento de Erros**: Mensagens de erro claras e informativas
5. **Performance**: Tempo de resposta aceitável
