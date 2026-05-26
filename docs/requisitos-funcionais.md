# Requisitos Funcionais - Mundo Invest

## Back-end

| ID | Requisito | Descrição | Critério de Aceite | Prioridade |
|----|-----------|-----------|-------------------|------------|
| RF-001 | Endpoint de criação de cliente | Sistema deve fornecer endpoint POST /clientes | Endpoint responde a requisições POST | Alta |
| RF-002 | Receber dados do cliente | API deve receber cliente_nome, cliente_email, tipo_solicitacao, valor_patrimonio | Request body contém todos os campos | Alta |
| RF-003 | Validar campos obrigatórios | Sistema deve validar que todos os campos estão presentes | Retorna HTTP 400 se campos faltando | Alta |
| RF-004 | Validar e-mail | Sistema deve validar formato do e-mail | Retorna HTTP 400 se e-mail inválido | Alta |
| RF-005 | Salvar cliente no banco | Sistema deve persistir cliente em banco local (SQLite ou PostgreSQL) | Cliente encontrado no banco após criação | Alta |
| RF-006 | Definir status inicial | Sistema deve definir status "Aguardando Análise" para novo cliente | Status no banco é "Aguardando Análise" | Alta |
| RF-007 | Estruturar mutation createCard | Sistema deve estruturar mutation GraphQL conforme especificação Pipefy | String/payload da mutation presente no código | Alta |
| RF-008 | Simular envio ao Pipefy | Sistema deve simular envio da mutation (não conectar ao Pipefy real) | Código contém mutation mas não faz requisição real | Alta |
| RF-009 | Retornar sucesso na criação | Sistema deve retornar HTTP 201 com dados do cliente criado | Response contém dados do cliente e card_id | Alta |
| RF-010 | Endpoint de webhook | Sistema deve fornecer endpoint POST /webhooks/pipefy/card-updated | Endpoint responde a requisições POST | Alta |
| RF-011 | Receber dados do webhook | API deve receber event_id, card_id, cliente_email, timestamp | Request body contém todos os campos | Alta |
| RF-012 | Verificar idempotência | Sistema deve verificar se event_id já foi processado | Retorna HTTP 200 sem processar se duplicado | Alta |
| RF-013 | Buscar cliente por e-mail | Sistema deve buscar cliente no banco usando cliente_email | Cliente encontrado se existir | Alta |
| RF-014 | Calcular prioridade alta | Sistema deve definir prioridade_alta se patrimonio >= 200.000 | Prioridade no banco é "prioridade_alta" | Alta |
| RF-015 | Calcular prioridade normal | Sistema deve definir prioridade_normal se patrimonio < 200.000 | Prioridade no banco é "prioridade_normal" | Alta |
| RF-016 | Estruturar mutation updateCard | Sistema deve estruturar mutation GraphQL conforme especificação Pipefy | String/payload da mutation presente no código | Alta |
| RF-017 | Atualizar status do cliente | Sistema deve mudar status para "Processado" após processamento | Status no banco é "Processado" | Alta |
| RF-018 | Salvar prioridade no banco | Sistema deve persistir prioridade calculada no banco | Prioridade no banco atualizada | Alta |
| RF-019 | Retornar sucesso no webhook | Sistema deve retornar HTTP 200 com dados do processamento | Response contém prioridade e status | Alta |
| RF-020 | Retornar erro cliente não encontrado | Sistema deve retornar HTTP 404 se cliente não existe | Response com mensagem de erro | Alta |
| RF-021 | Salvar evento processado | Sistema deve persistir event_id para controle de idempotência | Evento encontrado no banco após processamento | Alta |

## GraphQL (Pipefy Integration)

| ID | Requisito | Descrição | Critério de Aceite | Prioridade |
|----|-----------|-----------|-------------------|------------|
| RF-022 | Pesquisar documentação Pipefy | Desenvolvedor deve pesquisar documentação oficial do Pipefy | Código contém referência à documentação | Alta |
| RF-023 | Mutation createCard correta | Mutation deve seguir sintaxe oficial do Pipefy | Mutation estruturada corretamente no código | Alta |
| RF-024 | Mutation updateCard correta | Mutation deve seguir sintaxe oficial do Pipefy | Mutation estruturada corretamente no código | Alta |
| RF-025 | Variáveis da mutation | Mutation deve usar variáveis corretas (nome, e-mail, patrimônio) | Variáveis presentes no código | Alta |
| RF-026 | Comentários da fonte | Código deve ter comentários indicando fonte da especificação | Comentários presentes no código | Alta |

## Testes

| ID | Requisito | Descrição | Critério de Aceite | Prioridade |
|----|-----------|-----------|-------------------|------------|
| RF-027 | Teste de criação de cliente | Sistema deve ter teste automatizado para criação com payload válido | Teste passa e cliente salvo no banco | Alta |
| RF-028 | Teste de webhook com prioridade | Sistema deve ter teste para processamento de webhook aplicando regra de prioridade | Teste passa e prioridade correta definida | Alta |
| RF-029 | Teste de idempotência | Sistema deve ter teste para bloqueio de processamento de event_id duplicado | Teste passa e webhook não processado novamente | Alta |
| RF-030 | Teste de validação de campos | Sistema deve ter teste para validação de campos obrigatórios | Teste passa e retorna HTTP 400 | Alta |
| RF-031 | Teste de cliente não encontrado | Sistema deve ter teste para webhook com cliente inexistente | Teste passa e retorna HTTP 404 | Alta |
| RF-036 | Teste de validação de e-mail | Sistema deve ter teste para validação de formato de e-mail | Teste passa e retorna HTTP 400 | Alta |
| RF-037 | Teste de validação de patrimônio negativo | Sistema deve ter teste para rejeição de patrimônio negativo | Teste passa e retorna HTTP 400 | Alta |
| RF-038 | Teste de validação de patrimônio zero | Sistema deve ter teste para rejeição de patrimônio zero | Teste passa e retorna HTTP 400 | Alta |
| RF-039 | Teste de validação de nome vazio | Sistema deve ter teste para rejeição de nome vazio | Teste passa e retorna HTTP 400 | Alta |
| RF-040 | Teste de prioridade no limite | Sistema deve ter teste para patrimônio exatamente 200.000 | Teste passa e define prioridade_alta | Alta |
| RF-041 | Teste de validação de webhook | Sistema deve ter teste para campos obrigatórios do webhook | Teste passa e retorna HTTP 400 | Alta |
| RF-042 | Teste de event_id vazio | Sistema deve ter teste para rejeição de event_id vazio | Teste passa e retorna HTTP 400 | Alta |

## Documentação

| ID | Requisito | Descrição | Critério de Aceite | Prioridade |
|----|-----------|-----------|-------------------|------------|
| RF-032 | README com instruções de execução | Projeto deve ter README com como rodar localmente | README contém comandos para execução | Alta |
| RF-033 | README com instruções de testes | README deve ter como rodar os testes | README contém comandos para testes | Alta |
| RF-034 | README com exemplos de requisição | README deve ter exemplos curl para os dois endpoints | README contém exemplos funcionais | Alta |
| RF-035 | Visão de produção AWS | README deve explicar escalabilidade na AWS (opcional) | Texto explicando Lambda, API Gateway, DynamoDB/RDS | Baixa |
