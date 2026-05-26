# Regras de Negócio - Mundo Invest

## Regras de Cliente

| ID | Regra | Descrição | Prioridade |
|----|-------|-----------|------------|
| RN-001 | Campos obrigatórios | Cliente deve ter nome, e-mail, tipo de solicitação e valor de patrimônio | Alta |
| RN-002 | Validação de e-mail | E-mail deve estar em formato válido | Alta |
| RN-003 | Status inicial | Cliente criado deve ter status "Aguardando Análise" | Alta |
| RN-004 | Card ID | Cliente deve ter um card_id associado (simulação Pipefy) | Alta |

## Regras de Priorização

| ID | Regra | Descrição | Prioridade |
|----|-------|-----------|------------|
| RN-005 | Prioridade alta | Se valor_patrimonio >= 200.000, definir prioridade como "prioridade_alta" | Alta |
| RN-006 | Prioridade normal | Se valor_patrimonio < 200.000, definir prioridade como "prioridade_normal" | Alta |
| RN-007 | Limite de patrimônio | Limite para prioridade alta é 200.000 (fixo) | Alta |

## Regras de Webhook

| ID | Regra | Descrição | Prioridade |
|----|-------|-----------|------------|
| RN-008 | Idempotência | Webhook deve verificar event_id para evitar processamento duplicado | Alta |
| RN-009 | Campos obrigatórios | Webhook deve ter event_id, card_id, cliente_email e timestamp | Alta |
| RN-010 | Busca por e-mail | Cliente deve ser buscado no banco usando cliente_email | Alta |
| RN-011 | Atualização de status | Após processamento, status do cliente deve ser "Processado" | Alta |
| RN-012 | Atualização de prioridade | Prioridade calculada deve ser salva no banco | Alta |

## Regras de GraphQL (Pipefy)

| ID | Regra | Descrição | Prioridade |
|----|-------|-----------|------------|
| RN-013 | Mutation createCard | Deve seguir especificação oficial do Pipefy para criação de card | Alta |
| RN-014 | Mutation updateCard | Deve seguir especificação oficial do Pipefy para atualização de card | Alta |
| RN-015 | Simulação de envio | Mutations devem ser estruturadas mas não enviadas ao Pipefy real | Alta |
| RN-016 | Campos do card | Card deve conter nome, e-mail e patrimônio do cliente | Alta |

## Regras de Processamento

| ID | Regra | Descrição | Prioridade |
|----|-------|-----------|------------|
| RN-017 | Ordem de processamento | Prioridade deve ser calculada antes de atualizar status | Alta |
| RN-018 | Persistência de eventos | Eventos processados devem ser salvos para controle de idempotência | Alta |
| RN-019 | Timestamp | Webhook deve ter timestamp para rastreabilidade | Média |

## Regras de Validação

| ID | Regra | Descrição | Prioridade |
|----|-------|-----------|------------|
| RN-020 | Validação de nome | Nome não pode ser vazio | Alta |
| RN-021 | Validação de tipo de solicitação | Tipo de solicitação não pode ser vazio | Alta |
| RN-022 | Validação de patrimônio | Patrimônio deve ser um número positivo | Alta |
| RN-023 | Validação de event_id | Event_id não pode ser vazio | Alta |
| RN-024 | Validação de card_id | Card_id não pode ser vazio | Alta |
| RN-034 | Validação de e-mail no webhook | E-mail no webhook deve ter formato válido | Alta |
| RN-035 | Validação de timestamp | Timestamp deve estar em formato ISO 8601 | Média |

## Regras de Banco de Dados

| ID | Regra | Descrição | Prioridade |
|----|-------|-----------|------------|
| RN-025 | Persistência local | Cliente deve ser salvo em banco local (SQLite ou PostgreSQL) | Alta |
| RN-026 | Transações | Operações de atualização devem ser atômicas | Alta |
| RN-027 | Índices | E-mail deve ter índice para busca eficiente | Média |
| RN-028 | Event_id único | Event_id deve ter restrição de unicidade | Alta |

## Regras de API

| ID | Regra | Descrição | Prioridade |
|----|-------|-----------|------------|
| RN-029 | Endpoint POST /clientes | Deve aceitar payload JSON com dados do cliente | Alta |
| RN-030 | Endpoint POST /webhooks/pipefy/card-updated | Deve aceitar payload JSON com dados do webhook | Alta |
| RN-031 | Resposta de sucesso | Deve retornar HTTP 201 para criação, HTTP 200 para webhook | Alta |
| RN-032 | Resposta de erro | Deve retornar HTTP 400 para validação, HTTP 404 para não encontrado | Alta |
| RN-033 | Mensagens em Português | Todas as mensagens de erro devem estar em Português | Alta |
