# Tarefas - Sprint 4: Testes Automatizados

## 1. Testes de Criação de Cliente (Feature Folder: gestao_clientes)

- [ ] Criar internal/gestao_clientes/controller_test.go com testes do ClienteController
- [ ] Teste: criação com payload válido
- [ ] Teste: validação de campos obrigatórios (nome)
- [ ] Teste: validação de campos obrigatórios (email)
- [ ] Teste: validação de campos obrigatórios (tipo_solicitacao)
- [ ] Teste: validação de campos obrigatórios (valor_patrimonio)
- [ ] Teste: validação de e-mail inválido
- [ ] Teste: validação de patrimônio negativo
- [ ] Teste: validação de patrimônio zero
- [ ] Teste: validação de nome vazio
- [ ] Teste: erro de banco de dados

## 2. Testes de Webhook (Feature Folder: processamento_eventos)

- [ ] Criar internal/processamento_eventos/controller_test.go com testes do WebhookController
- [ ] Teste: processamento com prioridade alta (patrimonio >= 200.000)
- [ ] Teste: processamento com prioridade normal (patrimonio < 200.000)
- [ ] Teste: processamento no limite (patrimonio == 200.000)
- [ ] Teste: idempotência (event_id duplicado)
- [ ] Teste: cliente não encontrado
- [ ] Teste: validação de campos obrigatórios (event_id)
- [ ] Teste: validação de campos obrigatórios (card_id)
- [ ] Teste: validação de campos obrigatórios (cliente_email)
- [ ] Teste: validação de campos obrigatórios (timestamp)
- [ ] Teste: validação de event_id vazio
- [ ] Teste: validação de card_id vazio
- [ ] Teste: validação de cliente_email vazio
- [ ] Teste: validação de email inválido
- [ ] Teste: validação de timestamp inválido
- [ ] Teste: erro de banco de dados

## 3. Testes de Integração (Feature Folder: shared)

- [ ] Criar internal/shared/database/connection_test.go com testes de integração
- [ ] Teste: integração com banco de dados
- [ ] Teste: transações atômicas
- [ ] Teste: restrições de unicidade

## 4. Configuração de CI/CD

- [ ] Configurar pipeline para execução automática de testes
- [ ] Configurar relatórios de cobertura de código
- [ ] Configurar notificações de falha de testes
